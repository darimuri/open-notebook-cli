package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/darimuri/open-notebook-cli/internal/auth"
)

type Client struct {
	baseURL    string
	auth       *auth.Middleware
	httpClient *http.Client
	debug      bool
}

func NewClient(baseURL string, authMiddleware *auth.Middleware) *Client {
	return &Client{
		baseURL:    baseURL,
		auth:       authMiddleware,
		httpClient: http.DefaultClient,
		debug:      false,
	}
}

func NewClientWithDebug(baseURL string, authMiddleware *auth.Middleware, debug bool) *Client {
	return &Client{
		baseURL:    baseURL,
		auth:       authMiddleware,
		httpClient: http.DefaultClient,
		debug:      debug,
	}
}

func (c *Client) BaseURL() string {
	return c.baseURL
}

func (c *Client) NewRequest(method, path string, body interface{}) (*http.Request, error) {
	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(data)
		if c.debug {
			fmt.Fprintf(os.Stderr, "DEBUG: request body: %s\n", string(data))
		}
	}

	req, err := http.NewRequest(method, c.baseURL+path, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	c.auth.AddAuth(req)

	if c.debug {
		fmt.Fprintf(os.Stderr, "DEBUG: request headers: %v\n", req.Header)
	}

	return req, nil
}

func (c *Client) Do(req *http.Request, v interface{}) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if c.debug {
		fmt.Fprintf(os.Stderr, "DEBUG: %s %s -> %d\n", req.Method, req.URL.Path, resp.StatusCode)
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("API error: %d", resp.StatusCode)
	}

	if v != nil {
		return json.NewDecoder(resp.Body).Decode(v)
	}
	return nil
}

func (c *Client) Get(path string, v interface{}) error {
	req, err := c.NewRequest("GET", path, nil)
	if err != nil {
		return err
	}
	return c.Do(req, v)
}

func (c *Client) Post(path string, body interface{}, v interface{}) error {
	req, err := c.NewRequest("POST", path, body)
	if err != nil {
		return err
	}
	return c.Do(req, v)
}

func (c *Client) Put(path string, body interface{}, v interface{}) error {
	req, err := c.NewRequest("PUT", path, body)
	if err != nil {
		return err
	}
	return c.Do(req, v)
}

func (c *Client) Delete(path string, v interface{}) error {
	req, err := c.NewRequest("DELETE", path, nil)
	if err != nil {
		return err
	}
	return c.Do(req, v)
}