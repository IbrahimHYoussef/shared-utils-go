package testutils

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/testcontainers/testcontainers-go"
	tclog "github.com/testcontainers/testcontainers-go/log"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/testcontainers/testcontainers-go/modules/minio"
)

type S3TestCred struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	UseSSl    bool
}

type S3Container struct {
	MinioContainer *minio.MinioContainer
	Credintials    S3TestCred
}

// Defaults CreateMinioContainer uses for an empty image or empty credentials.
//
// DefaultMinioImage is pinned by digest: cgr.dev/chainguard/minio only
// publishes a floating "latest" tag, and pinning keeps local and CI runs on
// identical bits. Update the digest deliberately when bumping Minio.
const (
	DefaultMinioImage   = "cgr.dev/chainguard/minio@sha256:7df55d002b698e4494a131b848dd553361bbe58b7e19f62d9623c29af1ba9ccf"
	DefaultEndpoint     = "localhost"
	DefaultAccessKey    = "minioAdmin"
	DefaultSecretKey    = "minioAdmin"
	MinioStartupTimeout = 60 * time.Second
)

// testcontainers logs every container step when tests run with -v; turn its
// logs off.
func init() {
	tclog.SetDefault(log.New(io.Discard, "", 0))
}

func CreateMinioContainer(ctx context.Context, image string, cred *S3TestCred) (*S3Container, error) {
	if len(image) == 0 {
		image = DefaultMinioImage
	}
	resolved := S3TestCred{Endpoint: DefaultEndpoint, AccessKey: DefaultAccessKey, SecretKey: DefaultSecretKey, UseSSl: false}
	if cred != nil {
		if cred.Endpoint != "" {
			resolved.Endpoint = cred.Endpoint
		}
		if cred.AccessKey != "" {
			resolved.AccessKey = cred.AccessKey
		}
		if cred.SecretKey != "" {
			resolved.SecretKey = cred.SecretKey
		}
		resolved.UseSSl = cred.UseSSl
	}
	cred = &resolved

	minioContainer, err := minio.Run(ctx,
		image,
		minio.WithUsername(cred.AccessKey),
		minio.WithPassword(cred.SecretKey),
		// Wait on the readiness endpoint rather than a log line. Minio prints
		// its "API:" banner once; matching it twice only worked locally
		// because an unrelated storage error log also began with "API:".
		testcontainers.WithWaitStrategy(
			wait.ForHTTP("/minio/health/live").
				WithPort("9000/tcp").
				WithStartupTimeout(MinioStartupTimeout),
		),
	)
	if err != nil {
		return nil, err
	}
	return &S3Container{
		MinioContainer: minioContainer,
		Credintials:    *cred,
	}, nil
}
