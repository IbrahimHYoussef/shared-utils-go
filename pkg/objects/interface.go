package objects

import (
	"bytes"
	"context"
	"io"
	"time"
)

type StandardObjectStorage interface {
	Upload(ctx context.Context, bucket, contentType, key string, data *bytes.Buffer) error
	Download(ctx context.Context, bucket, key string) (io.ReadSeekCloser, error)
	Delete(ctx context.Context, bucket, key string) error
	UploadLink(ctx context.Context, bucket, key string, expire time.Duration) (string, error)
	DownloadLink(ctx context.Context, bucket, key string, expire time.Duration) (string, error)
}
