package untappd

import (
	"testing"
	"time"
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

// TestErrorError tests for consistent output from the Error.Error method.
func TestErrorError(t *testing.T) {
	var tests = []struct {
		description string
		code        int
		eType       string
		details     string
		developer   string
		result      string
	}{
		{
			description: "only details",
			code:        500,
			eType:       "auth_failed",
			details:     "authentication failed",
			developer:   "",
			result:      "500 [auth_failed]: authentication failed",
		},
		{
			description: "only developer friendly",
			code:        501,
			eType:       "auth_failed",
			details:     "",
			developer:   "authentication failed due to server error",
			result:      "501 [auth_failed]: authentication failed due to server error",
		},
		{
			description: "both details and developer friendly",
			code:        502,
			eType:       "auth_failed",
			details:     "authentication failed",
			developer:   "authentication failed due to server error",
			result:      "502 [auth_failed]: authentication failed due to server error",
		},
	}

	for _, tt := range tests {
		err := &Error{
			Code:              tt.code,
			Detail:            tt.details,
			Type:              tt.eType,
			DeveloperFriendly: tt.developer,
		}

		if res := err.Error(); res != tt.result {
			t.Fatalf("unexpected result string for test %q: %q != %q", tt.description, res, tt.result)
		}
	}
}

// Test_timeUnitToDuration verifies that timeUnitToDuration provides proper
// output for a variety of time number and measure values.
func Test_timeUnitToDuration(t *testing.T) {
	var tests = []struct {
		description string
		number      float64
		measure     string
		result      time.Duration
	}{
		{"0.05 milliseconds", 0.05, "milliseconds", time.Duration(5*time.Millisecond) / 100},
		{"5 milliseconds", 5, "milliseconds", time.Duration(5 * time.Millisecond)},
		{"500 milliseconds", 500, "milliseconds", time.Duration(500 * time.Millisecond)},
		{"0.5 seconds", 0.5, "seconds", time.Duration(500 * time.Millisecond)},
		{"1 seconds", 1, "seconds", time.Duration(1 * time.Second)},
		{"10 seconds", 10, "seconds", time.Duration(10 * time.Second)},
		{"0.5 minutes", 0.5, "minutes", time.Duration(30 * time.Second)},
		{"1 minutes", 1, "minutes", time.Duration(1 * time.Minute)},
		{"2 minutes", 2, "minutes", time.Duration(2 * time.Minute)},
		{"invalid: 100 hours", 100, "hours", time.Duration(0)},
		{"invalid: 10 days", 10, "days", time.Duration(0)},
		{"invalid: 1 lightyears", 1, "lightyears", time.Duration(0)},
	}

	for _, tt := range tests {
		if dur := timeUnitToDuration(tt.number, tt.measure); dur != tt.result {
			t.Fatalf("unexpected duration for test %q: %v != %v", tt.description, dur, tt.result)
		}
	}
}
