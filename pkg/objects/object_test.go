package objects_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/IbrahimHYoussef/shared-utils-go/pkg/objects"
	"github.com/IbrahimHYoussef/shared-utils-go/pkg/testutils"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/testcontainers/testcontainers-go"
)

// ObjectStorage must satisfy the interface the rest of the code depends on.
var _ objects.StandardObjectStorage = (*objects.ObjectStorage)(nil)

const testBucket = "test-bucket"

func CreateTestClient(endpoint, accessKey, secretKey string, useSsl bool) (*minio.Client, error) {
	client, err := minio.New(endpoint,
		&minio.Options{Creds: credentials.NewStaticV4(accessKey, secretKey, ""), Secure: useSsl})
	if err != nil {
		return nil, err
	}

	return client, nil
}

// storageEnv is a running Minio container with testBucket created, the raw
// client used to verify results, and the ObjectStorage under test.
type storageEnv struct {
	ctx      context.Context
	client   *minio.Client
	storage  *objects.ObjectStorage
	endpoint string
}

// newStorageEnv starts a Minio container, creates testBucket, and tears
// everything down when the test ends.
func newStorageEnv(t *testing.T) *storageEnv {
	t.Helper()
	testcontainers.SkipIfProviderIsNotHealthy(t)
	ctx := context.Background()

	container, err := testutils.CreateMinioContainer(ctx, "", nil)
	if err != nil {
		t.Fatalf("CreateMinioContainer: %v", err)
	}
	t.Cleanup(func() { _ = container.MinioContainer.Terminate(ctx) })

	endpoint, err := container.MinioContainer.ConnectionString(ctx)
	if err != nil {
		t.Fatalf("ConnectionString: %v", err)
	}

	client, err := CreateTestClient(endpoint, container.Credintials.AccessKey, container.Credintials.SecretKey, container.Credintials.UseSSl)
	if err != nil {
		t.Fatalf("CreateTestClient: %v", err)
	}
	if err := client.MakeBucket(ctx, testBucket, minio.MakeBucketOptions{}); err != nil {
		t.Fatalf("MakeBucket: %v", err)
	}

	return &storageEnv{
		ctx:      ctx,
		client:   client,
		storage:  objects.NewObjectStorage(client),
		endpoint: endpoint,
	}
}

// mustUpload puts data under key with the given content type or fails the test.
func (e *storageEnv) mustUpload(t *testing.T, key, contentType string, data []byte) {
	t.Helper()
	if err := e.storage.Upload(e.ctx, testBucket, contentType, key, bytes.NewBuffer(data)); err != nil {
		t.Fatalf("Upload(%q): %v", key, err)
	}
}

// readDownload downloads key and returns its bytes, or the first error from
// either the Download call or the read. minio-go fetches lazily, so a missing
// object may only surface on read.
func (e *storageEnv) readDownload(key string) ([]byte, error) {
	rc, err := e.storage.Download(e.ctx, testBucket, key)
	if err != nil {
		return nil, err
	}
	defer rc.Close()
	return io.ReadAll(rc)
}

func TestObjectStorage_Upload(t *testing.T) {
	env := newStorageEnv(t)

	t.Run("stores the bytes and content type", func(t *testing.T) {
		want := []byte(`{"hello":"world"}`)
		env.mustUpload(t, "docs/hello.json", "application/json", want)

		info, err := env.client.StatObject(env.ctx, testBucket, "docs/hello.json", minio.StatObjectOptions{})
		if err != nil {
			t.Fatalf("StatObject: %v", err)
		}
		if info.Size != int64(len(want)) {
			t.Fatalf("size %d, want %d", info.Size, len(want))
		}
		if info.ContentType != "application/json" {
			t.Fatalf("content type %q, want application/json", info.ContentType)
		}
	})

	t.Run("overwrites an existing key", func(t *testing.T) {
		env.mustUpload(t, "overwrite.txt", "text/plain", []byte("first"))
		env.mustUpload(t, "overwrite.txt", "text/plain", []byte("second version"))

		got, err := env.readDownload("overwrite.txt")
		if err != nil {
			t.Fatalf("Download: %v", err)
		}
		if string(got) != "second version" {
			t.Fatalf("got %q, want %q", got, "second version")
		}
	})

	t.Run("accepts an empty buffer", func(t *testing.T) {
		env.mustUpload(t, "empty.bin", "application/octet-stream", nil)

		info, err := env.client.StatObject(env.ctx, testBucket, "empty.bin", minio.StatObjectOptions{})
		if err != nil {
			t.Fatalf("StatObject: %v", err)
		}
		if info.Size != 0 {
			t.Fatalf("size %d, want 0", info.Size)
		}
	})

	t.Run("fails when the bucket does not exist", func(t *testing.T) {
		err := env.storage.Upload(env.ctx, "no-such-bucket", "text/plain", "k", bytes.NewBufferString("x"))
		if err == nil {
			t.Fatal("expected an error for a missing bucket")
		}
		if code := minio.ToErrorResponse(err).Code; code != "NoSuchBucket" {
			t.Fatalf("error code %q, want NoSuchBucket (%v)", code, err)
		}
	})
}

func TestObjectStorage_Download(t *testing.T) {
	env := newStorageEnv(t)

	t.Run("returns the uploaded bytes", func(t *testing.T) {
		want := bytes.Repeat([]byte("abc123"), 1000)
		env.mustUpload(t, "big.txt", "text/plain", want)

		got, err := env.readDownload("big.txt")
		if err != nil {
			t.Fatalf("Download: %v", err)
		}
		if !bytes.Equal(got, want) {
			t.Fatalf("downloaded %d bytes that differ from the %d uploaded", len(got), len(want))
		}
	})

	t.Run("reader supports seeking", func(t *testing.T) {
		env.mustUpload(t, "seek.txt", "text/plain", []byte("0123456789"))

		rc, err := env.storage.Download(env.ctx, testBucket, "seek.txt")
		if err != nil {
			t.Fatalf("Download: %v", err)
		}
		defer rc.Close()
		if _, err := rc.Seek(7, io.SeekStart); err != nil {
			t.Fatalf("Seek: %v", err)
		}
		got, err := io.ReadAll(rc)
		if err != nil {
			t.Fatalf("ReadAll: %v", err)
		}
		if string(got) != "789" {
			t.Fatalf("got %q after seeking to 7, want %q", got, "789")
		}
	})

	t.Run("fails for a missing key", func(t *testing.T) {
		_, err := env.readDownload("does-not-exist")
		if err == nil {
			t.Fatal("expected an error for a missing key")
		}
		if code := minio.ToErrorResponse(err).Code; code != "NoSuchKey" {
			t.Fatalf("error code %q, want NoSuchKey (%v)", code, err)
		}
	})
}

func TestObjectStorage_Delete(t *testing.T) {
	env := newStorageEnv(t)

	t.Run("removes the object", func(t *testing.T) {
		env.mustUpload(t, "to-delete.txt", "text/plain", []byte("bye"))

		if err := env.storage.Delete(env.ctx, testBucket, "to-delete.txt"); err != nil {
			t.Fatalf("Delete: %v", err)
		}

		_, err := env.client.StatObject(env.ctx, testBucket, "to-delete.txt", minio.StatObjectOptions{})
		if code := minio.ToErrorResponse(err).Code; code != "NoSuchKey" {
			t.Fatalf("object still present after delete: err=%v", err)
		}
	})

	t.Run("leaves other objects alone", func(t *testing.T) {
		env.mustUpload(t, "keep/a.txt", "text/plain", []byte("a"))
		env.mustUpload(t, "keep/b.txt", "text/plain", []byte("b"))

		if err := env.storage.Delete(env.ctx, testBucket, "keep/a.txt"); err != nil {
			t.Fatalf("Delete: %v", err)
		}

		got, err := env.readDownload("keep/b.txt")
		if err != nil {
			t.Fatalf("Download of untouched key: %v", err)
		}
		if string(got) != "b" {
			t.Fatalf("got %q, want %q", got, "b")
		}
	})

	t.Run("is a no-op for a missing key", func(t *testing.T) {
		if err := env.storage.Delete(env.ctx, testBucket, "never-existed"); err != nil {
			t.Fatalf("Delete of a missing key should succeed, got %v", err)
		}
	})
}

func TestObjectStorage_UploadLink(t *testing.T) {
	env := newStorageEnv(t)

	t.Run("link accepts a PUT that stores the object", func(t *testing.T) {
		link, err := env.storage.UploadLink(env.ctx, testBucket, "via-link.txt", time.Minute)
		if err != nil {
			t.Fatalf("UploadLink: %v", err)
		}
		assertPresignedURL(t, link, env.endpoint, testBucket+"/via-link.txt")

		req, err := http.NewRequestWithContext(env.ctx, http.MethodPut, link, bytes.NewBufferString("put through link"))
		if err != nil {
			t.Fatalf("NewRequest: %v", err)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("PUT to upload link: %v", err)
		}
		resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("PUT status %d, want 200", resp.StatusCode)
		}

		got, err := env.readDownload("via-link.txt")
		if err != nil {
			t.Fatalf("Download: %v", err)
		}
		if string(got) != "put through link" {
			t.Fatalf("got %q, want %q", got, "put through link")
		}
	})

	t.Run("rejects an expiry over seven days", func(t *testing.T) {
		if _, err := env.storage.UploadLink(env.ctx, testBucket, "k", 8*24*time.Hour); err == nil {
			t.Fatal("expected an error for an expiry over seven days")
		}
	})

	t.Run("rejects a zero expiry", func(t *testing.T) {
		if _, err := env.storage.UploadLink(env.ctx, testBucket, "k", 0); err == nil {
			t.Fatal("expected an error for a zero expiry")
		}
	})
}

func TestObjectStorage_DownloadLink(t *testing.T) {
	env := newStorageEnv(t)

	t.Run("link serves the object over GET", func(t *testing.T) {
		env.mustUpload(t, "shared.txt", "text/plain", []byte("shared content"))

		link, err := env.storage.DownloadLink(env.ctx, testBucket, "shared.txt", time.Minute)
		if err != nil {
			t.Fatalf("DownloadLink: %v", err)
		}
		assertPresignedURL(t, link, env.endpoint, testBucket+"/shared.txt")

		resp, err := http.Get(link)
		if err != nil {
			t.Fatalf("GET download link: %v", err)
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("GET status %d, want 200: %s", resp.StatusCode, body)
		}
		if string(body) != "shared content" {
			t.Fatalf("got %q, want %q", body, "shared content")
		}
	})

	t.Run("link for a missing key returns 404 when used", func(t *testing.T) {
		link, err := env.storage.DownloadLink(env.ctx, testBucket, "ghost.txt", time.Minute)
		if err != nil {
			t.Fatalf("DownloadLink: %v", err)
		}
		resp, err := http.Get(link)
		if err != nil {
			t.Fatalf("GET download link: %v", err)
		}
		resp.Body.Close()
		if resp.StatusCode != http.StatusNotFound {
			t.Fatalf("GET status %d, want 404", resp.StatusCode)
		}
	})

	t.Run("rejects an expiry over seven days", func(t *testing.T) {
		if _, err := env.storage.DownloadLink(env.ctx, testBucket, "k", 8*24*time.Hour); err == nil {
			t.Fatal("expected an error for an expiry over seven days")
		}
	})
}

// assertPresignedURL checks that link is an absolute URL pointing at the test
// endpoint, addressing wantPath, and carrying a signature.
func assertPresignedURL(t *testing.T, link, endpoint, wantPath string) {
	t.Helper()
	u, err := url.Parse(link)
	if err != nil {
		t.Fatalf("link %q is not a URL: %v", link, err)
	}
	if u.Host != endpoint {
		t.Fatalf("link host %q, want %q", u.Host, endpoint)
	}
	if u.Path != "/"+wantPath {
		t.Fatalf("link path %q, want %q", u.Path, "/"+wantPath)
	}
	if u.Query().Get("X-Amz-Signature") == "" {
		t.Fatalf("link %q has no X-Amz-Signature", link)
	}
}
