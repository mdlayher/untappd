package untappd

import (
	"errors"
	"net/url"
	"testing"
	"time"
)

var (
	errBadJSON = errors.New("invalid character '}' looking for beginning of value")
)

// Test_responseTimeUnmarshalJSON verifies that responseTime.UnmarshalJSON
// provides proper time.Duration for a variety of responseTime JSON values
// from the Untappd APIv4.
func Test_responseTimeUnmarshalJSON(t *testing.T) {
	var tests = []struct {
		description string
		body        []byte
		result      time.Duration
		err         error
	}{
		{
			description: "0.05 milliseconds",
			body:        []byte(`{"time":0.05,"measure":"milliseconds"}`),
			result:      time.Duration(5*time.Millisecond) / 100,
		},
		{
			description: "5 milliseconds",
			body:        []byte(`{"time":5,"measure":"milliseconds"}`),
			result:      time.Duration(5 * time.Millisecond),
		},
		{
			description: "500 milliseconds",
			body:        []byte(`{"time":500,"measure":"milliseconds"}`),
			result:      time.Duration(500 * time.Millisecond),
		},
		{
			description: "0.5 seconds",
			body:        []byte(`{"time":0.5,"measure":"seconds"}`),
			result:      time.Duration(500 * time.Millisecond),
		},
		{
			description: "1 seconds",
			body:        []byte(`{"time":1,"measure":"seconds"}`),
			result:      time.Duration(1 * time.Second),
		},
		{
			description: "10 seconds",
			body:        []byte(`{"time":10,"measure":"seconds"}`),
			result:      time.Duration(10 * time.Second),
		},
		{
			description: "0.5 minutes",
			body:        []byte(`{"time":0.5,"measure":"minutes"}`),
			result:      time.Duration(30 * time.Second),
		},
		{
			description: "1 minutes",
			body:        []byte(`{"time":1,"measure":"minutes"}`),
			result:      time.Duration(1 * time.Minute),
		},
		{
			description: "2 minutes",
			body:        []byte(`{"time":2,"measure":"minutes"}`),
			result:      time.Duration(2 * time.Minute),
		},
		{
			description: "invalid: 100 hours",
			body:        []byte(`{"time":100,"measure":"hours"}`),
			err:         errInvalidTimeUnit,
		},
		{
			description: "invalid: 10 days",
			body:        []byte(`{"time":10,"measure":"days"}`),
			err:         errInvalidTimeUnit,
		},
		{
			description: "invalid: 1 lightyears",
			body:        []byte(`{"time":1,"measure":"lightyears"}`),
			err:         errInvalidTimeUnit,
		},
		{
			description: "bad JSON",
			body:        []byte(`}`),
			err:         errBadJSON,
		},
	}

	for _, tt := range tests {
		r := new(responseTime)
		err := r.UnmarshalJSON(tt.body)
		if tt.err == nil && err != nil {
			t.Fatal(err)
		}
		if tt.err != nil && err.Error() != tt.err.Error() {
			t.Fatalf("unexpected error for test %q: %v != %v", tt.description, err, tt.err)
		}

		if *r != responseTime(tt.result) {
			t.Fatalf("unexpected duration for test %q: %v != %v", tt.description, r, tt.result)
		}
	}
}

// Test_responseURLUnmarshalJSON verifies that responseURL.UnmarshalJSON
// provides proper url.URL value for a variety of responseURL JSON values
// from the Untappd APIv4.
func Test_responseURLUnmarshalJSON(t *testing.T) {
	var tests = []struct {
		description string
		body        []byte
		result      url.URL
		err         error
	}{
		{
			description: "empty string",
			body:        []byte(`""`),
			result:      url.URL{},
		},
		{
			description: "scheme only",
			body:        []byte(`"https://"`),
			result:      url.URL{Scheme: "https"},
		},
		{
			description: "scheme and host",
			body:        []byte(`"http://foo.com:80"`),
			result:      url.URL{Scheme: "http", Host: "foo.com:80"},
		},
		{
			description: "scheme, host, and path",
			body:        []byte(`"http://foo.com:80/bar"`),
			result:      url.URL{Scheme: "http", Host: "foo.com:80", Path: "/bar"},
		},
		{
			description: "bad JSON",
			body:        []byte(`}`),
			err:         errBadJSON,
		},
	}

	for _, tt := range tests {
		r := new(responseURL)
		err := r.UnmarshalJSON(tt.body)
		if tt.err == nil && err != nil {
			t.Fatal(err)
		}
		if tt.err != nil && err.Error() != tt.err.Error() {
			t.Fatalf("unexpected error for test %q: %v != %v", tt.description, err, tt.err)
		}

		if *r != responseURL(tt.result) {
			t.Fatalf("unexpected url.URL for test %q: %#v != %#v", tt.description, r, tt.result)
		}
	}
}

// Test_responseBoolUnmarshalJSON verifies that responseBool.UnmarshalJSON
// provides proper bool value for a variety of responseBool JSON values
// from the Untappd APIv4.
func Test_responseBoolUnmarshalJSON(t *testing.T) {
	var tests = []struct {
		description string
		body        []byte
		result      bool
		err         error
	}{
		{
			description: "0 (false)",
			body:        []byte(`0`),
			result:      false,
		},
		{
			description: "1 (true)",
			body:        []byte(`1`),
			result:      true,
		},
		{
			description: "2 (invalid)",
			body:        []byte(`2`),
			err:         errInvalidBool,
		},
		{
			description: "bad JSON",
			body:        []byte(`}`),
			err:         errBadJSON,
		},
	}

	for _, tt := range tests {
		r := new(responseBool)
		err := r.UnmarshalJSON(tt.body)
		if tt.err == nil && err != nil {
			t.Fatal(err)
		}
		if tt.err != nil && err.Error() != tt.err.Error() {
			t.Fatalf("unexpected error for test %q: %v != %v", tt.description, err, tt.err)
		}

		if *r != responseBool(tt.result) {
			t.Fatalf("unexpected bool for test %q: %v != %v", tt.description, r, tt.result)
		}
	}
}
