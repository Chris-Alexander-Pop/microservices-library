// Package image provides an Image Generation interface.
package image

import "context"

// Service is the image generation interface.
type Service interface {
	Generate(ctx context.Context, prompt string, opts Options) ([]string, error)
}

// Options configures the generation.
type Options struct {
	N     int    `json:"n"`
	Size  string `json:"size"` // 1024x1024
	Model string `json:"model"`
}
