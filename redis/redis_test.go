package redis

import (
	"context"
	"log"
	"testing"
)

// TestGet tests the Get function from the redis package.
func TestGet(t *testing.T) {
	ctx := context.Background()

	k := "key"
	v, err := Get(ctx, k)
	if err == nil {
		t.Errorf("Expected error when Redis is not connected, but got nil")
	}

	// Connect to Redis before testing
	Connect()

	v, err = Get(ctx, "nonexistent_key")
	if err != nil {
		t.Errorf("Failed to get value for key %s: %v", k, err)
	}
	if v != "" {
		t.Errorf("Expected empty value for key %s, but got: %s", k, v)
	}

	v, err = Get(nil, k)
	if err == nil {
		t.Errorf("Expected error for nil context, but got value: %s", v)
	}

	v, err = Get(ctx, k)
	if err != nil {
		t.Errorf("Failed to get value for key %s: %v", k, err)
	}
	if v != "value" {
		t.Errorf("Expected value 'value' for key %s, but got: %s", k, v)
	}

	_ = v
	log.Println("v = ", v)
}
