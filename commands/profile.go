package commands

import (
	"fmt"
	"os"

	"uppies/cli/config"

	"github.com/spf13/cobra"
)

func ProfileCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "profile",
		Short: "Manage Uppies profiles",
	}

	cmd.AddCommand(
		profileListCmd(),
		profileCreateCmd(),
		profileRmCmd(),
		profileSetCmd(),
		profileGetCmd(),
	)

	return cmd
}

func profileListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all profiles",
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := config.GetConfig()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
				os.Exit(1)
			}

			if len(cfg.Profiles) == 0 {
				fmt.Println("No profiles found.")
				return
			}

			fmt.Println("Profiles:")
			for _, profile := range cfg.Profiles {
				active := ""
				if profile.Name == cfg.ActiveProfile {
					active = " (active)"
				}
				fmt.Printf("  %s%s\n", profile.Name, active)
			}
		},
	}
}

func profileCreateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create <name>",
		Short: "Create a new profile",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]

			cfg, err := config.GetConfig()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
				os.Exit(1)
			}

			// Check if profile already exists
			for _, profile := range cfg.Profiles {
				if profile.Name == name {
					fmt.Fprintf(os.Stderr, "Profile '%s' already exists\n", name)
					os.Exit(1)
				}
			}

			// Add new profile
			cfg.Profiles = append(cfg.Profiles, config.Profile{Name: name, Host: "", Token: ""})
			if err := config.SaveConfigData(cfg); err != nil {
				fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("Profile '%s' created\n", name)
		},
	}
}

func profileRmCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "rm <name>",
		Short: "Remove a profile",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]

			cfg, err := config.GetConfig()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
				os.Exit(1)
			}

			// Find and remove profile
			found := false
			newProfiles := make([]config.Profile, 0, len(cfg.Profiles))
			for _, profile := range cfg.Profiles {
				if profile.Name == name {
					found = true
					if profile.Name == cfg.ActiveProfile {
						fmt.Fprintf(os.Stderr, "Cannot remove active profile '%s'\n", name)
						os.Exit(1)
					}
				} else {
					newProfiles = append(newProfiles, profile)
				}
			}

			if !found {
				fmt.Fprintf(os.Stderr, "Profile '%s' does not exist\n", name)
				os.Exit(1)
			}

			cfg.Profiles = newProfiles
			if err := config.SaveConfigData(cfg); err != nil {
				fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("Profile '%s' removed\n", name)
		},
	}
}

func profileSetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "set <name>",
		Short: "Set the active profile",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]

			cfg, err := config.GetConfig()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
				os.Exit(1)
			}

			// Check if profile exists
			found := false
			for _, profile := range cfg.Profiles {
				if profile.Name == name {
					found = true
					break
				}
			}

			if !found {
				fmt.Fprintf(os.Stderr, "Profile '%s' does not exist\n", name)
				os.Exit(1)
			}

			cfg.ActiveProfile = name
			if err := config.SaveConfigData(cfg); err != nil {
				fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("Active profile set to '%s'\n", name)
		},
	}
}

func profileGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get",
		Short: "Get the active profile",
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := config.GetConfig()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("Active profile: %s\n", cfg.ActiveProfile)
		},
	}
}