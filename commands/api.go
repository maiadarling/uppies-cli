package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"uppies/cli/config"
)

type APIClient struct {
	BaseURL string
	Token   string
}

func NewAPIClient() *APIClient {
	return &APIClient{
		BaseURL: "http://api.uppiesplz.com:3000",
		Token:   config.Token,
	}
}

type UploadRequest struct {
	Data string `json:"data"`
}

type UploadResponse struct {
	Data struct {
		Name   string `json:"name"`
		URL    string `json:"url"`
		Status string `json:"status"`
	} `json:"data"`
}

type SiteResponse struct {
	Data struct {
		Name   string `json:"name"`
		URL    string `json:"url"`
		Status string `json:"status"`
	} `json:"data"`
}

func (c *APIClient) UploadSite(data string) (UploadResponse, error) {
	var respData UploadResponse
	reqBody := UploadRequest{Data: data}
	jsonBytes, err := json.Marshal(reqBody)
	if err != nil {
		return respData, fmt.Errorf("marshaling request: %w", err)
	}

	req, err := http.NewRequest("POST", c.BaseURL+"/sites", bytes.NewReader(jsonBytes))
	if err != nil {
		return respData, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return respData, fmt.Errorf("sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return respData, fmt.Errorf("upload failed with status: %d", resp.StatusCode)
	}

	err = json.NewDecoder(resp.Body).Decode(&respData)
	if err != nil {
		return respData, fmt.Errorf("decoding response: %w", err)
	}

	return respData, nil
}

func (c *APIClient) GetSite(name string) (SiteResponse, error) {
	var respData SiteResponse
	req, err := http.NewRequest("GET", c.BaseURL+"/sites/"+name, nil)
	if err != nil {
		return respData, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return respData, fmt.Errorf("sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return respData, fmt.Errorf("get site failed with status: %d", resp.StatusCode)
	}

	err = json.NewDecoder(resp.Body).Decode(&respData)
	if err != nil {
		return respData, fmt.Errorf("decoding response: %w", err)
	}

	return respData, nil
}