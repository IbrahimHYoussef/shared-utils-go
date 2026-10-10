package testutils_test

import (
	"context"
	"testing"

	"github.com/IbrahimHYoussef/shared-utils-go/pkg/testutils"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/testcontainers/testcontainers-go"
)

func TestCreateS3Container(t *testing.T) {
	testcontainers.SkipIfProviderIsNotHealthy(t)
	ctx := context.Background()

	cases := map[string]*testutils.S3TestCred{
		"nil creds use defaults": nil,
		"given creds are kept":   {Endpoint: "http://localhost:9000", AccessKey: testutils.DefaultAccessKey, SecretKey: testutils.DefaultSecretKey, UseSSl: false},
		"empty fields fall back to test":{Endpoint: "http//:localhost:9000"},
	}

	for name, cred := range cases {
		t.Run(name, func(t *testing.T) {
			container, err := testutils.CreateMinioContainer(ctx, "", cred)
			if err != nil {
				t.Fatalf("CreateMinoContainer: %v", err)
			}
			t.Cleanup(func() { container.MinioContainer.Terminate(ctx) })
			t.Log(container.MinioContainer.Username)
			endpoint := container.Credintials.Endpoint
			accessKey := container.Credintials.AccessKey
			secredKey := container.Credintials.SecretKey
			useSsl := container.Credintials.UseSSl

			minio, err := minio.New(endpoint,
				&minio.Options{
					Creds:  credentials.NewStaticV4(accessKey, secredKey, ""),
					Secure: useSsl,
				},
			)
			if err != nil {
				t.Fatalf("connect: %v", err)
			}
			online := minio.IsOnline()
			if !online {
				t.Fatalf("Object Storage is not Online")
			}
		})
	}
}
