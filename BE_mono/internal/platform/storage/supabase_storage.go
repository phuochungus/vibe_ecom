package storage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	storage_go "github.com/supabase-community/storage-go"
)

type SupabaseConfig struct {
	ProjectURL     string
	ServiceRoleKey string
	Bucket         string
	PublicBaseURL  string
}

type supabaseBucketClient interface {
	GetBucket(id string) (storage_go.Bucket, error)
	CreateBucket(id string, options storage_go.BucketOptions) (storage_go.Bucket, error)
}

type supabaseObjectClient interface {
	UploadFile(bucketID string, relativePath string, data io.Reader, fileOptions ...storage_go.FileOptions) (storage_go.FileUploadResponse, error)
	GetPublicUrl(bucketID string, filePath string, urlOptions ...storage_go.UrlOptions) storage_go.SignedUrlResponse
}

type supabaseClient interface {
	supabaseBucketClient
	supabaseObjectClient
}

type SupabaseStorage struct {
	client        supabaseClient
	bucket        string
	publicBaseURL string
}

func NewSupabaseStorage(ctx context.Context, cfg SupabaseConfig) (*SupabaseStorage, error) {
	projectURL := strings.TrimSpace(cfg.ProjectURL)
	serviceRoleKey := strings.TrimSpace(cfg.ServiceRoleKey)
	bucket := strings.TrimSpace(cfg.Bucket)

	if projectURL == "" {
		return nil, fmt.Errorf("supabase project url is required")
	}
	if serviceRoleKey == "" {
		return nil, fmt.Errorf("supabase service role key is required")
	}
	if bucket == "" {
		return nil, fmt.Errorf("supabase bucket is required")
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	client := storage_go.NewClient(resolveSupabaseStorageAPIBaseURL(projectURL), serviceRoleKey, nil)
	store := &SupabaseStorage{
		client:        client,
		bucket:        bucket,
		publicBaseURL: resolveSupabasePublicBaseURL(cfg),
	}
	if err := store.ensureBucket(ctx); err != nil {
		return nil, err
	}

	return store, nil
}

func (s *SupabaseStorage) Upload(ctx context.Context, objectKey string, reader io.Reader, size int64, contentType string) (string, error) {
	if s == nil || s.client == nil {
		return "", fmt.Errorf("supabase storage client is not configured")
	}
	if err := ctx.Err(); err != nil {
		return "", err
	}

	key := strings.TrimLeft(strings.TrimSpace(objectKey), "/")
	if key == "" {
		return "", fmt.Errorf("object key is required")
	}
	if size <= 0 {
		return "", fmt.Errorf("object size must be greater than 0")
	}

	trimmedContentType := strings.TrimSpace(contentType)
	upsert := false
	_, err := s.client.UploadFile(s.bucket, key, io.LimitReader(reader, size), storage_go.FileOptions{
		ContentType: &trimmedContentType,
		Upsert:      &upsert,
	})
	if err != nil {
		return "", fmt.Errorf("upload file to supabase: %w", err)
	}

	return s.PublicURL(key), nil
}

func (s *SupabaseStorage) PublicURL(objectKey string) string {
	if s == nil {
		return ""
	}

	key := strings.TrimLeft(strings.TrimSpace(objectKey), "/")
	if key == "" {
		return strings.TrimRight(s.publicBaseURL, "/") + "/object/public/" + s.bucket
	}
	return strings.TrimRight(s.publicBaseURL, "/") + "/object/public/" + s.bucket + "/" + encodeObjectKey(key)
}

func (s *SupabaseStorage) ensureBucket(ctx context.Context) error {
	if s == nil || s.client == nil {
		return fmt.Errorf("supabase storage client is not configured")
	}
	if err := ctx.Err(); err != nil {
		return err
	}

	bucket, err := s.client.GetBucket(s.bucket)
	if err == nil {
		if !bucket.Public {
			return fmt.Errorf("supabase bucket %q must be public", s.bucket)
		}
		return nil
	}
	if !isSupabaseBucketNotFound(err) {
		return fmt.Errorf("check supabase bucket: %w", err)
	}

	_, err = s.client.CreateBucket(s.bucket, storage_go.BucketOptions{Public: true})
	if err != nil {
		return fmt.Errorf("create supabase bucket: %w", err)
	}
	return nil
}

func resolveSupabaseStorageAPIBaseURL(projectURL string) string {
	trimmed := strings.TrimRight(strings.TrimSpace(projectURL), "/")
	if trimmed == "" {
		return ""
	}
	if strings.HasSuffix(trimmed, "/storage/v1") {
		return trimmed
	}
	return trimmed + "/storage/v1"
}

func resolveSupabasePublicBaseURL(cfg SupabaseConfig) string {
	explicit := strings.TrimRight(strings.TrimSpace(cfg.PublicBaseURL), "/")
	if explicit != "" {
		return explicit
	}
	return resolveSupabaseStorageAPIBaseURL(cfg.ProjectURL)
}

func isSupabaseBucketNotFound(err error) bool {
	if err == nil {
		return false
	}

	var storageErr *storage_go.StorageError
	if errors.As(err, &storageErr) {
		return storageErr.Status == 404
	}
	return strings.Contains(strings.ToLower(err.Error()), "not found")
}
