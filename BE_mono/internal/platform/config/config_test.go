package config

import "testing"

func TestLoadSupabaseStorageConfig(t *testing.T) {
	t.Setenv("SUPABASE_STORAGE_ENABLED", "true")
	t.Setenv("SUPABASE_URL", "https://project.supabase.co")
	t.Setenv("SUPABASE_SERVICE_ROLE_KEY", "service-role")
	t.Setenv("SUPABASE_STORAGE_BUCKET", "product-images")
	t.Setenv("SUPABASE_STORAGE_PUBLIC_BASE_URL", "https://cdn.example.com/storage")

	cfg := Load("be-mono", "8080")

	if !cfg.SupabaseStorageEnabled {
		t.Fatal("expected supabase storage to be enabled")
	}
	if cfg.SupabaseURL != "https://project.supabase.co" {
		t.Fatalf("expected supabase url to be normalized, got %q", cfg.SupabaseURL)
	}
	if cfg.SupabaseServiceRoleKey != "service-role" {
		t.Fatalf("unexpected service role key %q", cfg.SupabaseServiceRoleKey)
	}
	if cfg.SupabaseStorageBucket != "product-images" {
		t.Fatalf("unexpected supabase bucket %q", cfg.SupabaseStorageBucket)
	}
	if cfg.SupabaseStoragePublicBaseURL != "https://cdn.example.com/storage" {
		t.Fatalf("unexpected supabase public base url %q", cfg.SupabaseStoragePublicBaseURL)
	}
}
