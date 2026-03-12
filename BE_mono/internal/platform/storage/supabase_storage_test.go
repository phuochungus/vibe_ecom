package storage

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"

	storage_go "github.com/supabase-community/storage-go"
)

type fakeSupabaseClient struct {
	getBucketFn    func(id string) (storage_go.Bucket, error)
	createBucketFn func(id string, options storage_go.BucketOptions) (storage_go.Bucket, error)
	uploadFileFn   func(bucketID string, relativePath string, data io.Reader, fileOptions ...storage_go.FileOptions) (storage_go.FileUploadResponse, error)
	getPublicURLFn func(bucketID string, filePath string, urlOptions ...storage_go.UrlOptions) storage_go.SignedUrlResponse
}

func (f *fakeSupabaseClient) GetBucket(id string) (storage_go.Bucket, error) {
	if f.getBucketFn == nil {
		return storage_go.Bucket{}, nil
	}
	return f.getBucketFn(id)
}

func (f *fakeSupabaseClient) CreateBucket(id string, options storage_go.BucketOptions) (storage_go.Bucket, error) {
	if f.createBucketFn == nil {
		return storage_go.Bucket{}, nil
	}
	return f.createBucketFn(id, options)
}

func (f *fakeSupabaseClient) UploadFile(bucketID string, relativePath string, data io.Reader, fileOptions ...storage_go.FileOptions) (storage_go.FileUploadResponse, error) {
	if f.uploadFileFn == nil {
		return storage_go.FileUploadResponse{}, nil
	}
	return f.uploadFileFn(bucketID, relativePath, data, fileOptions...)
}

func (f *fakeSupabaseClient) GetPublicUrl(bucketID string, filePath string, urlOptions ...storage_go.UrlOptions) storage_go.SignedUrlResponse {
	if f.getPublicURLFn == nil {
		return storage_go.SignedUrlResponse{}
	}
	return f.getPublicURLFn(bucketID, filePath, urlOptions...)
}

func TestResolveSupabaseStorageAPIBaseURL(t *testing.T) {
	got := resolveSupabaseStorageAPIBaseURL("https://project.supabase.co")
	if got != "https://project.supabase.co/storage/v1" {
		t.Fatalf("expected derived storage api base url, got %q", got)
	}

	got = resolveSupabaseStorageAPIBaseURL("https://project.supabase.co/storage/v1/")
	if got != "https://project.supabase.co/storage/v1" {
		t.Fatalf("expected existing storage api base url to be preserved, got %q", got)
	}
}

func TestSupabasePublicURLEncodesObjectKey(t *testing.T) {
	store := &SupabaseStorage{
		bucket:        "golf-store",
		publicBaseURL: "https://project.supabase.co/storage/v1",
	}

	got := store.PublicURL("products/2026/03/my image #1.png")
	want := "https://project.supabase.co/storage/v1/object/public/golf-store/products/2026/03/my%20image%20%231.png"
	if got != want {
		t.Fatalf("expected public url %q, got %q", want, got)
	}
}

func TestSupabaseUploadReturnsErrorForNilReceiver(t *testing.T) {
	var store *SupabaseStorage

	url, err := store.Upload(context.Background(), "products/test.png", bytes.NewReader([]byte("x")), 1, "image/png")
	if err == nil {
		t.Fatal("expected upload error for nil receiver")
	}
	if url != "" {
		t.Fatalf("expected empty url, got %q", url)
	}
}

func TestSupabaseEnsureBucketCreatesMissingPublicBucket(t *testing.T) {
	client := &fakeSupabaseClient{
		getBucketFn: func(id string) (storage_go.Bucket, error) {
			if id != "golf-store" {
				t.Fatalf("unexpected bucket lookup %q", id)
			}
			return storage_go.Bucket{}, &storage_go.StorageError{Status: 404, Message: "not found"}
		},
		createBucketFn: func(id string, options storage_go.BucketOptions) (storage_go.Bucket, error) {
			if id != "golf-store" {
				t.Fatalf("unexpected bucket create %q", id)
			}
			if !options.Public {
				t.Fatal("expected bucket to be created as public")
			}
			return storage_go.Bucket{Id: id, Public: true}, nil
		},
	}

	store := &SupabaseStorage{client: client, bucket: "golf-store", publicBaseURL: "https://project.supabase.co/storage/v1"}
	if err := store.ensureBucket(context.Background()); err != nil {
		t.Fatalf("expected missing bucket to be created, got %v", err)
	}
}

func TestSupabaseEnsureBucketRejectsPrivateBucket(t *testing.T) {
	client := &fakeSupabaseClient{
		getBucketFn: func(id string) (storage_go.Bucket, error) {
			return storage_go.Bucket{Id: id, Public: false}, nil
		},
	}

	store := &SupabaseStorage{client: client, bucket: "golf-store", publicBaseURL: "https://project.supabase.co/storage/v1"}
	err := store.ensureBucket(context.Background())
	if err == nil {
		t.Fatal("expected error for private bucket")
	}
}

func TestSupabaseUploadUsesBucketAndReturnsPublicURL(t *testing.T) {
	content := []byte("image-bytes")
	client := &fakeSupabaseClient{
		uploadFileFn: func(bucketID string, relativePath string, data io.Reader, fileOptions ...storage_go.FileOptions) (storage_go.FileUploadResponse, error) {
			if bucketID != "golf-store" {
				t.Fatalf("expected golf-store bucket, got %q", bucketID)
			}
			if relativePath != "products/2026/test image.png" {
				t.Fatalf("unexpected relative path %q", relativePath)
			}
			body, err := io.ReadAll(data)
			if err != nil {
				t.Fatalf("read upload body: %v", err)
			}
			if !bytes.Equal(body, content) {
				t.Fatal("uploaded content mismatch")
			}
			if len(fileOptions) != 1 {
				t.Fatalf("expected one file option set, got %d", len(fileOptions))
			}
			if fileOptions[0].ContentType == nil || *fileOptions[0].ContentType != "image/png" {
				t.Fatal("expected image/png content type")
			}
			if fileOptions[0].Upsert == nil || *fileOptions[0].Upsert {
				t.Fatal("expected upsert to be disabled")
			}
			return storage_go.FileUploadResponse{Key: relativePath}, nil
		},
	}

	store := &SupabaseStorage{
		client:        client,
		bucket:        "golf-store",
		publicBaseURL: "https://project.supabase.co/storage/v1",
	}

	url, err := store.Upload(context.Background(), "products/2026/test image.png", bytes.NewReader(content), int64(len(content)), "image/png")
	if err != nil {
		t.Fatalf("expected upload to succeed, got %v", err)
	}
	want := "https://project.supabase.co/storage/v1/object/public/golf-store/products/2026/test%20image.png"
	if url != want {
		t.Fatalf("expected public url %q, got %q", want, url)
	}
}

func TestSupabaseUploadReturnsContextError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	store := &SupabaseStorage{
		client:        &fakeSupabaseClient{},
		bucket:        "golf-store",
		publicBaseURL: "https://project.supabase.co/storage/v1",
	}
	_, err := store.Upload(ctx, "products/test.png", bytes.NewReader([]byte("x")), 1, "image/png")
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context canceled, got %v", err)
	}
}
