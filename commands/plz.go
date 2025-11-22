package commands

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"uppies/cli/api"
	"uppies/cli/internal/terminal"
	"uppies/cli/internal/utils"
)

func resolvePath(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		fmt.Printf("Error resolving absolute path for %s: %v\n", path, err)
		os.Exit(1)
	}
	return abs
}

func packageFolder(folder string) []byte {
	var buf bytes.Buffer
	zipWriter := zip.NewWriter(&buf)

	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		relPath, err := filepath.Rel(folder, path)
		if err != nil {
			return err
		}
		if relPath == "." {
			return nil
		}
		if info.IsDir() {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		w, err := zipWriter.Create(relPath)
		if err != nil {
			return err
		}

		_, err = io.Copy(w, f)
		return err
	})

	if err != nil {
		fmt.Printf("Error packaging folder %s: %v\n", folder, err)
		os.Exit(1)
	}

	if err := zipWriter.Close(); err != nil {
		fmt.Printf("Error closing zip writer: %v\n", err)
		os.Exit(1)
	}

	return buf.Bytes()
}

func uploadArchive(archive []byte) api.UploadResponse {
	encoded := base64.StdEncoding.EncodeToString(archive)
	client := api.NewAPIClient()
	resp, err := client.UploadSite(encoded)
	if err != nil {
		fmt.Println("Error uploading:", err)
		os.Exit(1)
	}
	return resp
}

func verifySite(name string) {
	client := api.NewAPIClient()
	for {
		resp, err := client.GetSite(name)
		if err != nil {
			fmt.Println("Error verifying site:", err)
			os.Exit(1)
		}

		if resp.Data.Status == "live" {
			break
		}

		if resp.Data.Status == "error" {
			fmt.Println("\nSite deployment failed.")
			os.Exit(1)
		}

		time.Sleep(terminal.StatusRetryInterval)
		verifySite(name)
	}
}

func plzRun(cmd *cobra.Command, args []string) {
	var absFolder string
	var zipBytes []byte
	var resp api.UploadResponse

	folder := args[0]

	fmt.Println("\\o/ Uppies!")
	fmt.Println("")

	terminal.RunStage("Processing", func() { absFolder = resolvePath(folder) })
	terminal.RunStage("Archiving", func() { zipBytes = packageFolder(absFolder) })
	terminal.RunStage("Uploading", func() { resp = uploadArchive(zipBytes) })
	terminal.RunStage("Verifying", func() { verifySite(resp.Data.Name) })

	fmt.Printf("You can now access your site at: %s\n", resp.Data.URL)
}

func PlzCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "plz [folder]",
		Short:   "Execute the plz command",
		Args:    cobra.ExactArgs(1),
		PreRun:  utils.RequireLogin,
		Run:     plzRun,
	}
}