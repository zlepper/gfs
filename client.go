package gfs

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
)

type Client struct {
	url    *url.URL
	token  string
	client http.Client
}

func urlJoin(p1, p2 string) string {
	if !strings.HasSuffix(p1, "/") {
		p1 = p1 + "/"
	}

	if strings.HasPrefix(p2, "/") {
		p2 = p2[1:]
	}

	return p1 + p2
}

func (c *Client) getUrl(p string) (string, error) {
	target, err := url.Parse(p)
	if err != nil {
		return "", err
	}

	u := c.url.ResolveReference(target)

	return u.String(), nil
}

func (c *Client) setHeaders(req *http.Request) {
	req.Header.Set("gfs-token", c.token)
	req.Header.Set("accept", FormatJson)
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

	sUrl, err := c.getUrl("/login")
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", sUrl, &buf)
	if err != nil {
		return err
	}
	req.Header.Set("accept", FormatJson)
	req.Header.Set("Content-Type", FormatJson)

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
	sUrl, err := c.getUrl(p)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("GET", sUrl, nil)
	if err != nil {
		return err
	}

	c.setHeaders(req)

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
	// The path on which the file is uploaded on the GFS server
	UploadPath string
}

// Creates a new instance of upload file.
func NewUploadFile(filename, uploadPath string, reader io.ReadCloser) UploadFile {
	return UploadFile{
		Filename:   filename,
		Reader:     reader,
		UploadPath: uploadPath,
	}
}

// Helper method to quickly create new UploadFile from a file on the disk
// Make sure to call `f.Reader.Close()` to make sure no leaks happens in
// case of errors
func NewUploadFileFromDisk(filepath, uploadPath string) (f UploadFile, err error) {
	file, err := os.Open(filepath)
	if err != nil {
		return f, err
	}

	return NewUploadFile(path.Base(filepath), uploadPath, file), nil
}

// Call this to upload files
// Pass a function to progressUpdater to receive updates about the progress of the upload
func (c *Client) UploadFiles(files []UploadFile) error {
	for _, file := range files {
		err := c.UploadFile(file)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) UploadFile(file UploadFile) error {
	defer file.Reader.Close()

	sUrl, err := c.getUrl("/upload")
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", sUrl, file.Reader)
	if err != nil {
		return err
	}

	c.setHeaders(req)
	req.Header.Set("Content-Type", FormatOctetStream)

	q := req.URL.Query()

	uploadPath := urlJoin(file.UploadPath, file.Filename)

	q.Add("filename", uploadPath)
	req.URL.RawQuery = q.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		var response invalidRequest
		err = json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			return err
		}

		if response.Error != "" {
			return errors.New(response.Error)
		}
	}

	return nil
}

func NewClient(host, username, password string) (*Client, error) {

	u, err := url.Parse(host)
	if err != nil {
		return nil, err
	}

	client := Client{
		url: u,
	}

	err = client.Login(username, password)
	if err != nil {
		return nil, err
	}

	return &client, nil
}
