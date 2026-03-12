package service

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"

	"golf-store/be-mono/internal/platform/storage"
)

var samplePNG = []byte{
	0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a,
	0x00, 0x00, 0x00, 0x0d, 0x49, 0x48, 0x44, 0x52,
	0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
	0x08, 0x02, 0x00, 0x00, 0x00, 0x90, 0x77, 0x53,
	0xde, 0x00, 0x00, 0x00, 0x0c, 0x49, 0x44, 0x41,
	0x54, 0x08, 0x99, 0x63, 0x60, 0x60, 0x00, 0x00,
	0x00, 0x04, 0x00, 0x01, 0x0b, 0xe7, 0x02, 0x9d,
	0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4e, 0x44,
	0xae, 0x42, 0x60, 0x82,
}

type fakeImageStorage struct {
	uploadFn func(ctx context.Context, objectKey string, content []byte, size int64, contentType string) (string, error)
}

func (f *fakeImageStorage) Upload(ctx context.Context, objectKey string, reader io.Reader, size int64, contentType string) (string, error) {
	content, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}
	if f.uploadFn == nil {
		return "", nil
	}
	return f.uploadFn(ctx, objectKey, content, size, contentType)
}

func TestAdminUploadImageReturnsServiceUnavailableForTypedNilStorage(t *testing.T) {
	var minioStore *storage.MinIO

	svc := New(nil, minioStore)

	out, apiErr := svc.AdminUploadImage(context.Background(), "image.webp", []byte("fake-image"))
	if out != nil {
		t.Fatalf("expected no upload output, got %#v", out)
	}
	if apiErr == nil {
		t.Fatal("expected api error, got nil")
	}
	if apiErr.Status != 503 {
		t.Fatalf("expected status 503, got %d", apiErr.Status)
	}
	if apiErr.Code != "IMAGE_STORAGE_NOT_CONFIGURED" {
		t.Fatalf("expected IMAGE_STORAGE_NOT_CONFIGURED, got %q", apiErr.Code)
	}
}

func TestAdminUploadImageRejectsNonImagePayload(t *testing.T) {
	svc := New(nil, &fakeImageStorage{})

	out, apiErr := svc.AdminUploadImage(context.Background(), "notes.txt", []byte("not an image"))
	if out != nil {
		t.Fatalf("expected no upload output, got %#v", out)
	}
	if apiErr == nil {
		t.Fatal("expected api error, got nil")
	}
	if apiErr.Status != 400 {
		t.Fatalf("expected status 400, got %d", apiErr.Status)
	}
	if apiErr.Code != "VALIDATION_ERROR" {
		t.Fatalf("expected VALIDATION_ERROR, got %q", apiErr.Code)
	}
}

func TestAdminUploadImageUploadsAndReturnsMetadata(t *testing.T) {
	store := &fakeImageStorage{
		uploadFn: func(_ context.Context, objectKey string, content []byte, size int64, contentType string) (string, error) {
			if size != int64(len(samplePNG)) {
				t.Fatalf("expected size %d, got %d", len(samplePNG), size)
			}
			if contentType != "image/png" {
				t.Fatalf("expected image/png, got %q", contentType)
			}
			if !bytes.Equal(content, samplePNG) {
				t.Fatalf("uploaded content mismatch")
			}
			if objectKey == "" {
				t.Fatal("expected object key to be generated")
			}
			return "http://127.0.0.1:9000/golf-store/" + objectKey, nil
		},
	}
	svc := New(nil, store)

	out, apiErr := svc.AdminUploadImage(context.Background(), "image.png", samplePNG)
	if apiErr != nil {
		t.Fatalf("expected no api error, got %v", apiErr)
	}
	if out == nil {
		t.Fatal("expected upload output, got nil")
	}
	if out.URL == "" {
		t.Fatal("expected URL in upload output")
	}
	if out.ObjectKey == "" {
		t.Fatal("expected object key in upload output")
	}
	if out.ContentType != "image/png" {
		t.Fatalf("expected image/png content type, got %q", out.ContentType)
	}
	if out.Size != int64(len(samplePNG)) {
		t.Fatalf("expected size %d, got %d", len(samplePNG), out.Size)
	}
}

func TestAdminUploadImageReturnsInternalErrorWhenStorageFails(t *testing.T) {
	store := &fakeImageStorage{
		uploadFn: func(_ context.Context, _ string, _ []byte, _ int64, _ string) (string, error) {
			return "", errors.New("upload failed")
		},
	}
	svc := New(nil, store)

	out, apiErr := svc.AdminUploadImage(context.Background(), "image.png", samplePNG)
	if out != nil {
		t.Fatalf("expected no upload output, got %#v", out)
	}
	if apiErr == nil {
		t.Fatal("expected api error, got nil")
	}
	if apiErr.Status != 500 {
		t.Fatalf("expected status 500, got %d", apiErr.Status)
	}
	if apiErr.Code != "IMAGE_UPLOAD_FAILED" {
		t.Fatalf("expected IMAGE_UPLOAD_FAILED, got %q", apiErr.Code)
	}
}
