package storage

import (
	"context"
	"strings"
	"testing"
)

func TestResolvePublicBaseURL(t *testing.T) {
	t.Run("uses explicit public base url", func(t *testing.T) {
		got := resolvePublicBaseURL(MinIOConfig{
			Endpoint:      "minio:9000",
			PublicBaseURL: "https://cdn.example.com/storage/",
			UseSSL:        false,
		})
		if got != "https://cdn.example.com/storage" {
			t.Fatalf("expected explicit public base url to win, got %q", got)
		}
	})

	t.Run("builds from endpoint and ssl flag", func(t *testing.T) {
		got := resolvePublicBaseURL(MinIOConfig{
			Endpoint: "minio:9000",
			UseSSL:   true,
		})
		if got != "https://minio:9000" {
			t.Fatalf("expected derived public base url, got %q", got)
		}
	})
}

func TestEncodeObjectKey(t *testing.T) {
	got := encodeObjectKey("products/2026/03/my image #1.png")
	want := "products/2026/03/my%20image%20%231.png"
	if got != want {
		t.Fatalf("expected encoded object key %q, got %q", want, got)
	}
}

func TestUploadReturnsErrorForNilReceiver(t *testing.T) {
	var store *MinIO

	url, err := store.Upload(context.Background(), "products/test.png", strings.NewReader("x"), 1, "image/png")
	if err == nil {
		t.Fatal("expected upload error for nil receiver")
	}
	if url != "" {
		t.Fatalf("expected empty url, got %q", url)
	}
}
