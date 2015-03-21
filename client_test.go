package untappd

import (
	"testing"
)

// TestNewClient tests for all possible errors which can occur during a call
// to NewClient.
func TestNewClient(t *testing.T) {
	var tests = []struct {
		description  string
		clientID     string
		clientSecret string
		expErr       error
	}{
		{"no client ID or client secret", "", "", ErrNoClientID},
		{"no client ID", "", "bar", ErrNoClientID},
		{"no client secret", "foo", "", ErrNoClientSecret},
		{"ok", "foo", "bar", nil},
	}

	for _, tt := range tests {
		if _, err := NewClient(tt.clientID, tt.clientSecret, nil); err != tt.expErr {
			t.Fatalf("unexpected error for test %q: %v != %v", tt.description, err, tt.expErr)
		}
	}
}
