package utils

import (
	"context"
	"testing"
	"time"
)

func TestIsServerAlive(t *testing.T) {
	tests := []struct {
		name string
		addr string
		want bool
	}{
		{"with scheme", "http://google.com", true},
		{"without scheme", "google.com", false}, // должно фейлиться
		{"https", "https://google.com", true},
		{"invalid", "http://definitely-not-real-domain-12345.xyz", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			got := IsServerAlive(ctx, tt.addr)
			if got != tt.want {
				t.Errorf("IsServerAlive(%s) = %v, want %v", tt.addr, got, tt.want)
			}
		})
	}
}
