package untappd

import (
	"errors"
	"net/url"
	"reflect"
	"testing"
	"time"
)

var (
	errBadJSON = errors.New("invalid character '}' looking for beginning of value")
)

// Test_responseDurationUnmarshalJSON verifies that responseDuration.UnmarshalJSON
// provides proper time.Duration for a variety of responseDuration JSON values
// from the Untappd APIv4.
func Test_responseDurationUnmarshalJSON(t *testing.T) {
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
		r := new(responseDuration)
		err := r.UnmarshalJSON(tt.body)
		if tt.err == nil && err != nil {
			t.Fatal(err)
		}
		if tt.err != nil && err.Error() != tt.err.Error() {
			t.Fatalf("unexpected error for test %q: %v != %v", tt.description, err, tt.err)
		}

		if *r != responseDuration(tt.result) {
			t.Fatalf("unexpected duration for test %q: %v != %v", tt.description, r, tt.result)
		}
	}
}

// Test_responseTimeUnmarshalJSON verifies that responseTime.UnmarshalJSON
// provides proper time.Time for a variety of responseTime JSON values
// from the Untappd APIv4.
func Test_responseTimeUnmarshalJSON(t *testing.T) {
	mst, err := time.LoadLocation("America/Denver")
	if err != nil {
		t.Fatal(err)
	}

	var tests = []struct {
		description string
		body        []byte
		result      time.Time
		err         error
	}{
		{
			description: "default format",
			body:        []byte(`"` + time.RFC1123Z + `"`),
			result:      time.Date(2006, time.January, 2, 15, 4, 5, 0, mst),
		},
		{
			description: "bad time",
			body:        []byte(`"01-01-2001"`),
			err:         errors.New(`parsing time "01-01-2001" as "Mon, 02 Jan 2006 15:04:05 -0700": cannot parse "01-01-2001" as "Mon"`),
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

		ry, rm, rd := time.Time(*r).Date()
		ty, tm, td := tt.result.Date()
		if ry != ty {
			t.Fatalf("unexpected year for test %q: %d != %d", tt.description, ry, ty)
		}
		if rm != tm {
			t.Fatalf("unexpected month for test %q: %d != %d", tt.description, rm, tm)
		}
		if rd != td {
			t.Fatalf("unexpected day for test %q: %d != %d", tt.description, rd, td)
		}

		rh, rmi, rs := time.Time(*r).Clock()
		th, tmi, ts := tt.result.Clock()
		if rh != th {
			t.Fatalf("unexpected hour time for test %q: %d != %d", tt.description, rh, th)
		}
		if rmi != tmi {
			t.Fatalf("unexpected minute time for test %q: %d != %d", tt.description, rmi, tmi)
		}
		if rs != ts {
			t.Fatalf("unexpected second time for test %q: %d != %d", tt.description, rs, ts)
		}
	}
}

// Test_responseURLUnmarshalJSON verifies that responseURL.UnmarshalJSON
// provides proper url.URL value for a variety of responseURL JSON values
// from the Untappd APIv4.
func Test_responseURLUnmarshalJSON(t *testing.T) {
	// Bad URL used to validate URL parsing
	badURL := "http://www.%20.com/foo"

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
		// Test data courtesy of Bryan Liles, thanks
		// https://github.com/bryanl
		{
			description: "bad URL",
			body:        []byte(`"` + badURL + `"`),
			err: &url.Error{
				Op:  "parse",
				URL: badURL,
			},
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

		// Due to a change in Go tip's net/url library, we have to check
		// for individual fields on URL parse error
		uErr, ok := err.(*url.Error)
		if tt.err != nil && err != nil && ok {
			if uErr.Op != "parse" || uErr.URL != badURL {
				t.Fatalf("unexpected URL parse error: %v", uErr)
			}

			continue
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

// Test_responseBadgeLevelsUnmarshalJSON verifies that responseBadgeLevels.UnmarshalJSON
// provides proper badge count and items values for a variety of responseBadgeLevels
// JSON values from the Untappd APIv4.
func Test_responseBadgeLevelsUnmarshalJSON(t *testing.T) {
	var tests = []struct {
		description string
		body        []byte
		result      responseBadgeLevels
		err         error
	}{
		{
			description: "no badge levels (special API case)",
			body:        []byte(`[]`),
			result:      responseBadgeLevels{},
		},
		{
			description: "no badge levels (possibly non-existant case)",
			body:        []byte(`{"count":0,"items":[]}`),
			result: responseBadgeLevels{
				Count: 0,
				Items: []*rawBadge{},
			},
		},
		{
			description: "1 badge level",
			body:        []byte(`{"count":1,"items":[{"badge_name":"Foo (Level 1)"}]}`),
			result: responseBadgeLevels{
				Count: 1,
				Items: []*rawBadge{
					&rawBadge{
						Name: "Foo (Level 1)",
					},
				},
			},
		},
		{
			description: "2 badge levels",
			body:        []byte(`{"count":2,"items":[{"badge_name":"Foo (Level 2)"},{"badge_name":"Foo (Level 1)"}]}`),
			result: responseBadgeLevels{
				Count: 2,
				Items: []*rawBadge{
					&rawBadge{
						Name: "Foo (Level 2)",
					},
					&rawBadge{
						Name: "Foo (Level 1)",
					},
				},
			},
		},
		{
			description: "bad JSON",
			body:        []byte(`}`),
			err:         errBadJSON,
		},
	}

	for _, tt := range tests {
		r := new(responseBadgeLevels)
		err := r.UnmarshalJSON(tt.body)
		if tt.err == nil && err != nil {
			t.Fatal(err)
		}
		if tt.err != nil && err.Error() != tt.err.Error() {
			t.Fatalf("unexpected error for test %q: %v != %v", tt.description, err, tt.err)
		}

		if !reflect.DeepEqual(*r, responseBadgeLevels(tt.result)) {
			t.Fatalf("unexpected responseBadgeLevels for test %q: %v != %v", tt.description, r, tt.result)
		}
	}
}

// Test_responseVenueUnmarshalJSON verifies that responseVenue.UnmarshalJSON
// provides proper rawVenue output for a variety of responseVenue
// JSON values from the Untappd APIv4.
func Test_responseVenueUnmarshalJSON(t *testing.T) {
	var tests = []struct {
		description string
		body        []byte
		result      responseVenue
		err         error
	}{
		{
			description: "no venue (empty array, special API case)",
			body:        []byte(`[]`),
			result:      responseVenue{},
		},
		{
			description: "no venue (empty object, possibly non-existant case)",
			body:        []byte(`{}`),
			result:      responseVenue{},
		},
		{
			description: "venue exists",
			body:        []byte(`{"venue_id":1,"venue_name":"foo"}`),
			result: responseVenue{
				ID:   1,
				Name: "foo",
			},
		},
		{
			description: "bad JSON",
			body:        []byte(`}`),
			err:         errBadJSON,
		},
	}

	for _, tt := range tests {
		r := new(responseVenue)
		err := r.UnmarshalJSON(tt.body)
		if tt.err == nil && err != nil {
			t.Fatal(err)
		}
		if tt.err != nil && err.Error() != tt.err.Error() {
			t.Fatalf("unexpected error for test %q: %v != %v", tt.description, err, tt.err)
		}

		if !reflect.DeepEqual(*r, responseVenue(tt.result)) {
			t.Fatalf("unexpected responseVenue for test %q: %v != %v", tt.description, r, tt.result)
		}
	}
}
