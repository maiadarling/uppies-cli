package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"uppies/cli/config"
)

type APIClient struct {
	BaseURL *url.URL
	Token   string
	Client  *http.Client
}

func NewAPIClient() *APIClient {
	u, _ := url.Parse(config.Host)

	return &APIClient{
		BaseURL: u,
		Token:   config.Token,
		Client:  &http.Client{},
	}
}

// MARK: Types
type Item struct {
	Name    string   `json:"name"`
	URL     string   `json:"url"`
	Domains []string `json:"domains"`
	Status  string   `json:"status"`
}

type Pagination struct {
}

type SuccessfulResponse[T any] struct {
	Data       T           `json:"data"`
	Pagination *Pagination `json:"pagination,omitempty"`
	Message    *string     `json:"message,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type APIResponse[T any] struct {
	Data       T           `json:"data,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
	Message    *string     `json:"message,omitempty"`
	Error      *string     `json:"error,omitempty"`
}

type SingleResponse = APIResponse[Item]
type ListResponse = APIResponse[[]Item]

type UploadRequest struct {
	Data string `json:"data"`
}

// MARK: Helpers
func (c *APIClient) request(method, path string, reqBody any, respBody any) error {
	rel, err := url.Parse(path)
	if err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}

	fullURL := c.BaseURL.ResolveReference(rel)

	// Encode request body if provided
	var buf *bytes.Buffer
	if reqBody != nil {
		b, err := json.Marshal(reqBody)
		if err != nil {
			return fmt.Errorf("marshal request: %w", err)
		}
		buf = bytes.NewBuffer(b)
	} else {
		buf = &bytes.Buffer{}
	}

	req, err := http.NewRequest(method, fullURL.String(), buf)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}

	req.Header.Set("X-Uppies-Key", c.Token)
	if reqBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return fmt.Errorf("perform request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	// Decode JSON response into `respBody`
	if respBody != nil {
		if err := json.NewDecoder(resp.Body).Decode(respBody); err != nil {
			return fmt.Errorf("decode response: %w", err)
		}
	}

	return nil
}

// MARK: API Methods
func (c *APIClient) UploadSite(data string) (SingleResponse, error) {
	var out SingleResponse
	err := c.request("POST", "/sites", UploadRequest{Data: data}, &out)
	return out, err
}

func (c *APIClient) GetSite(name string) (SingleResponse, error) {
	var out SingleResponse
	err := c.request("GET", "/sites/"+name, nil, &out)
	return out, err
}

func (c *APIClient) ListSites() (ListResponse, error) {
	var out ListResponse
	err := c.request("GET", "/sites", nil, &out)
	return out, err
}
