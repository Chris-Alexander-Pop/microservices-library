package blob

import (
	"context"
	"io"
)

type BlobStorage interface {
	Upload(ctx context.Context, key string, data io.Reader) error
	Download(ctx context.Context, key string) (io.ReadCloser, error)
	Delete(ctx context.Context, key string) error
	URL(key string) string
}

type LocalStore struct {
	BaseDir string
}
