package service

import (
	"context"
	"testing"

	"golf-store/be-mono/internal/platform/storage"
)

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
