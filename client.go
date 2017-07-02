package gfs

import (
	"net/http"
	"path"
	"bytes"
	"encoding/json"
	"errors"
)

type Client struct {
	url string
	token string
	client *http.Client
}

func (c *Client) getUrl(p string) string {
	return path.Join(c.url, p)
}

func (c *Client) Login(username, password string) (error) {
	loginRequest := LoginRequest{
		Username:username,
		Password:password,
	}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(loginRequest)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", c.getUrl("/login"), &buf)
	req.Header.Set("accept", FormatJson)
	if err != nil {
		return err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()


	var response AuthorizationResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return err
	}

	if response.Error != "" {
		return errors.New(response.Error)
	}

	if response.Token == "" {
		return errors.New("No error on login, but no token was returned. This shouldn't be able to happen...")
	}

	c.token = response.Token

	return nil
}

func (c *Client) getContent(p string, out interface{}) error {
	req, err := http.NewRequest("GET", c.getUrl(p), nil)
	if err != nil {
		return err
	}

	req.Header.Set("token", c.token)
	req.Header.Set("accept", FormatJson)

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(out)
	if err != nil {
		return err
	}

	return nil
}

// Gets the content of the given directory
func (c *Client) GetDirectoryContent(p string) (*DirectoryStats, error) {
	var stats DirectoryStats
	err := c.getContent(p, &stats)
	return &stats, err
}

// Gets the metadata about the given file
func (c *Client) GetFileData(p string) (*FileStats, error) {
	var stats FileStats
	err := c.getContent(p, &stats)
	return &stats, err
}

func NewClient(url, username, password string) (*Client, error) {
	client := Client{
		url:url,
	}

	err := client.Login(username, password)
	if err != nil {
		return nil, err
	}

	return &client, nil
}
