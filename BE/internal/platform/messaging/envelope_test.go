package messaging

import "testing"

func TestNewEnvelopeSetsDefaults(t *testing.T) {
	env := NewEnvelope("api-gateway", "", "order.create.requested", "", nil)

	if env.ID == "" {
		t.Fatalf("expected id to be generated")
	}
	if env.CorrelationID == "" {
		t.Fatalf("expected correlation id to be generated")
	}
	if env.Version != "v1" {
		t.Fatalf("expected default version v1, got %s", env.Version)
	}
	if string(env.Payload) != "{}" {
		t.Fatalf("expected default payload {}")
	}
	if err := env.Validate(); err != nil {
		t.Fatalf("expected valid envelope, got %v", err)
	}
}
