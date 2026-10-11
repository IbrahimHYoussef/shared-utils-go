package objects

import (
	"bytes"
	"context"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
)

type ObjectStorage struct {
	minioClient *minio.Client
}

func NewObjectStorage(clinet *minio.Client) *ObjectStorage {
	return &ObjectStorage{
		minioClient: clinet,
	}
}

func (o *ObjectStorage) Upload(ctx context.Context, bucket, contentType string, key string, data *bytes.Buffer) error {
	_, err := o.minioClient.PutObject(ctx, bucket, key, data, int64(data.Len()),
		minio.PutObjectOptions{ContentType: contentType})
	return err
}

func (o *ObjectStorage) Download(ctx context.Context, bucket, key string) (io.ReadSeekCloser, error) {
	return o.minioClient.GetObject(ctx, bucket, key, minio.GetObjectOptions{})
}

func (o *ObjectStorage) Delete(ctx context.Context, bucket, key string) error {
	return o.minioClient.RemoveObject(ctx, bucket, key, minio.RemoveObjectOptions{})
}

func (o *ObjectStorage) UploadLink(ctx context.Context, bucket, key string, expire time.Duration) (string, error) {

	url, err := o.minioClient.PresignedPutObject(ctx, bucket, key, expire)
	if err != nil {
		return "", err
	}

	return url.String(), nil
}

func (o *ObjectStorage) DownloadLink(ctx context.Context, bucket, key string, expire time.Duration) (string, error) {
	url, err := o.minioClient.PresignedGetObject(ctx, bucket, key, expire, nil)
	if err != nil {
		return "", err
	}
	return url.String(), nil
}
