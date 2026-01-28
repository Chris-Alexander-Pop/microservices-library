package hdfs

import (
	"io"
	"os"

	"github.com/colinmarc/hdfs/v2"
)

type Client struct {
	client *hdfs.Client
}

func New(namenode string, user string) (*Client, error) {
	c, err := hdfs.NewClient(hdfs.ClientOptions{
		Addresses: []string{namenode},
		User:      user,
	})
	if err != nil {
		return nil, err
	}
	return &Client{client: c}, nil
}

func (c *Client) Create(path string) (io.WriteCloser, error) {
	return c.client.Create(path)
}

func (c *Client) Open(path string) (io.ReadCloser, error) {
	return c.client.Open(path)
}

func (c *Client) ReadDir(path string) ([]os.FileInfo, error) {
	return c.client.ReadDir(path)
}

func (c *Client) Close() error {
	return c.client.Close()
}
