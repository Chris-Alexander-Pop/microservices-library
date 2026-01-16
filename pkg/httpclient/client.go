package httpclient

import (
	"net/http"
	"time"
)

type Client struct {
	*http.Client
}

func New() *Client {
	return &Client{
		Client: &http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:       100,
				IdleConnTimeout:    90 * time.Second,
				DisableCompression: true,
			},
		},
	}
}
