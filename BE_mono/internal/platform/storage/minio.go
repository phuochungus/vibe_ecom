package storage

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinIOConfig struct {
	Endpoint      string
	PublicBaseURL string
	AccessKey     string
	SecretKey     string
	Bucket        string
	UseSSL        bool
}

type MinIO struct {
	client        *minio.Client
	bucket        string
	publicBaseURL string
}

func NewMinIO(ctx context.Context, cfg MinIOConfig) (*MinIO, error) {
	endpoint := strings.TrimSpace(cfg.Endpoint)
	accessKey := strings.TrimSpace(cfg.AccessKey)
	secretKey := strings.TrimSpace(cfg.SecretKey)
	bucket := strings.TrimSpace(cfg.Bucket)

	if endpoint == "" {
		return nil, fmt.Errorf("minio endpoint is required")
	}
	if accessKey == "" || secretKey == "" {
		return nil, fmt.Errorf("minio credentials are required")
	}
	if bucket == "" {
		return nil, fmt.Errorf("minio bucket is required")
	}

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("create minio client: %w", err)
	}

	store := &MinIO{
		client:        client,
		bucket:        bucket,
		publicBaseURL: resolvePublicBaseURL(cfg),
	}
	if err := store.ensureBucket(ctx); err != nil {
		return nil, err
	}
	if err := store.ensurePublicReadPolicy(ctx); err != nil {
		return nil, err
	}

	return store, nil
}

func (m *MinIO) Upload(ctx context.Context, objectKey string, reader io.Reader, size int64, contentType string) (string, error) {
	key := strings.TrimLeft(strings.TrimSpace(objectKey), "/")
	if key == "" {
		return "", fmt.Errorf("object key is required")
	}
	if size <= 0 {
		return "", fmt.Errorf("object size must be greater than 0")
	}

	_, err := m.client.PutObject(ctx, m.bucket, key, reader, size, minio.PutObjectOptions{
		ContentType: strings.TrimSpace(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("put object: %w", err)
	}

	return m.PublicURL(key), nil
}

func (m *MinIO) PublicURL(objectKey string) string {
	key := strings.TrimLeft(strings.TrimSpace(objectKey), "/")
	if key == "" {
		return strings.TrimRight(m.publicBaseURL, "/") + "/" + m.bucket
	}
	return strings.TrimRight(m.publicBaseURL, "/") + "/" + m.bucket + "/" + encodeObjectKey(key)
}

func (m *MinIO) ensureBucket(ctx context.Context) error {
	exists, err := m.client.BucketExists(ctx, m.bucket)
	if err != nil {
		return fmt.Errorf("check minio bucket: %w", err)
	}
	if exists {
		return nil
	}
	if err := m.client.MakeBucket(ctx, m.bucket, minio.MakeBucketOptions{}); err != nil {
		return fmt.Errorf("create minio bucket: %w", err)
	}
	return nil
}

func (m *MinIO) ensurePublicReadPolicy(ctx context.Context) error {
	policy := fmt.Sprintf(`{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:GetObject"],"Resource":["arn:aws:s3:::%s/*"]}]}`, m.bucket)
	if err := m.client.SetBucketPolicy(ctx, m.bucket, policy); err != nil {
		return fmt.Errorf("set minio bucket policy: %w", err)
	}
	return nil
}

func resolvePublicBaseURL(cfg MinIOConfig) string {
	baseURL := strings.TrimRight(strings.TrimSpace(cfg.PublicBaseURL), "/")
	if baseURL != "" {
		return baseURL
	}

	scheme := "http"
	if cfg.UseSSL {
		scheme = "https"
	}
	return scheme + "://" + strings.TrimSpace(cfg.Endpoint)
}

func encodeObjectKey(value string) string {
	parts := strings.Split(strings.TrimSpace(value), "/")
	encoded := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" {
			continue
		}
		encoded = append(encoded, url.PathEscape(trimmed))
	}
	return strings.Join(encoded, "/")
}
