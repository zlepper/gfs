package gfs

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"path"
)

type Client struct {
	url    string
	token  string
	client *http.Client
}

func (c *Client) getUrl(p string) string {
	return path.Join(c.url, p)
}

func (c *Client) Login(username, password string) error {
	loginRequest := LoginRequest{
		Username: username,
		Password: password,
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

	req.Header.Set("gfs-token", c.token)
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

// Requirements for uploading a file to GFS
type UploadFile struct {
	// The name of the file to upload
	Filename string
	// The file reader that allows reading the file
	Reader io.ReadCloser
	// The length of the file
	FileLength int64
}

// Creates a new instance of upload file.
func NewUploadFile(filename string, reader io.ReadCloser, fileLength int64) UploadFile {
	return UploadFile{
		Filename:   filename,
		Reader:     reader,
		FileLength: fileLength,
	}
}

// Helper method to quickly create new UploadFile from a file on the disk
// Make sure to call `f.Reader.Close()` to make sure no leaks happens in
// case of errors
func NewUploadFileFromDisk(filepath string) (f UploadFile, err error) {
	stats, err := os.Stat(filepath)
	if err != nil {
		return f, err
	}

	file, err := os.Open(filepath)
	if err != nil {
		return f, err
	}

	return NewUploadFile(path.Base(filepath), file, stats.Size()), nil
}

func (c *Client) UploadFiles(files []UploadFile, targetPath string) error {
	//
	//rp, wp, err := os.Pipe()
	//if err != nil {
	//	return err
	//}
	//
	//bodyWriter := multipart.NewWriter(wp)
	//
	//
	//
	//bodyWriter.CreateFormFile()

	return errors.New("gfs.UploadFiles is not yet supported.")
}

func NewClient(url, username, password string) (*Client, error) {
	client := Client{
		url: url,
	}

	err := client.Login(username, password)
	if err != nil {
		return nil, err
	}

	return &client, nil
}
