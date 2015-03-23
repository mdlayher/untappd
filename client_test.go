package untappd

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
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

// TestClient_requestContainsAPIKeys verifies that both client_id and client_secret
// are always present in API requests.
func TestClient_requestContainsAPIKeys(t *testing.T) {
	method := "GET"
	c, done := testClient(t, func(t *testing.T, w http.ResponseWriter, r *http.Request) {
		if m := r.Method; m != method {
			t.Fatalf("unexpected method: %q != %q", m, method)
		}

		q := r.URL.Query()

		if q.Get("client_id") == "" {
			t.Fatal("empty client_id query parameter")
		}
		if q.Get("client_secret") == "" {
			t.Fatal("empty client_secret query parameter")
		}
	})
	defer done()

	if _, err := c.request(method, "foo", nil, nil); err != nil {
		t.Fatal(err)
	}
}

// TestClient_requestContainsQueryParameters verifies that all custom query
// parameters are present in API requests.
func TestClient_requestContainsQueryParameters(t *testing.T) {
	method := "POST"
	c, done := testClient(t, func(t *testing.T, w http.ResponseWriter, r *http.Request) {
		if m := r.Method; m != method {
			t.Fatalf("unexpected method: %q != %q", m, method)
		}

		q := r.URL.Query()

		if s := q.Get("foo"); s != "bar" {
			t.Fatalf("unexpected query parameter: %q != %q", s, "bar")
		}
		if s := q.Get("bar"); s != "baz" {
			t.Fatalf("unexpected query parameter: %q != %q", s, "baz")
		}

		s, ok := q["baz"]
		if !ok {
			t.Fatal("missing query parameter: baz")
		}
		for _, ss := range s {
			if ss != "qux" && ss != "corge" {
				t.Fatal("did not find \"qux\" or \"corge\" in key \"baz\"")
			}
		}
	})
	defer done()

	if _, err := c.request(method, "foo", url.Values{
		"foo": []string{"bar"},
		"bar": []string{"baz"},
		"baz": []string{"qux", "corge"},
	}, nil); err != nil {
		t.Fatal(err)
	}
}

// TestClient_requestContainsHeaders verifies that all typical headers are set
// by the client during an API request.
func TestClient_requestContainsHeaders(t *testing.T) {
	method := "PUT"
	c, done := testClient(t, func(t *testing.T, w http.ResponseWriter, r *http.Request) {
		if m := r.Method; m != method {
			t.Fatalf("unexpected method: %q != %q", m, method)
		}

		h := r.Header

		if s := h.Get("Accept"); s != jsonContentType {
			t.Fatalf("unexpected Accept header: %q != %q", s, jsonContentType)
		}
		if s := h.Get("Content-Type"); s != jsonContentType {
			t.Fatalf("unexpected Content-Type header: %q != %q", s, jsonContentType)
		}

		if s := h.Get("User-Agent"); s != untappdUserAgent {
			t.Fatalf("unexpected User-Agent header: %q != %q", s, untappdUserAgent)
		}
	})
	defer done()

	if _, err := c.request(method, "foo", nil, nil); err != nil {
		t.Fatal(err)
	}
}

// TestClient_requestContainsBody verifies that a response body can be
// unmarshaled from JSON following an API request.
func TestClient_requestContainsBody(t *testing.T) {
	method := "GET"
	c, done := testClient(t, func(t *testing.T, w http.ResponseWriter, r *http.Request) {
		if m := r.Method; m != method {
			t.Fatalf("unexpected method: %q != %q", m, method)
		}

		// Use canned JSON with HTTP 500, though the HTTP code here will
		// return 200, for processing
		w.Write(apiErrJSON)
	})
	defer done()

	var v struct {
		Meta struct {
			Code int `json:"code"`
		} `json:"meta"`
	}

	if _, err := c.request(method, "foo", nil, &v); err != nil {
		t.Fatal(err)
	}

	if c := v.Meta.Code; c != http.StatusInternalServerError {
		t.Fatalf("unexpected code in response body: %d != %d", c, http.StatusInternalServerError)
	}
}

// Test_checkResponseWrongContentType verifies that checkResponse returns an error
// when the Content-Type header does not indicate application/json.
func Test_checkResponseWrongContentType(t *testing.T) {
	withHTTPResponse(t, http.StatusOK, "foo/bar", nil, func(t *testing.T, res *http.Response) {
		if err := checkResponse(res); err.Error() != "expected application/json content type, but received foo/bar" {
			t.Fatal(err)
		}
	})
}

// Test_checkResponseEOF verifies that checkResponse returns an io.EOF when no
// JSON body is found in the HTTP response body.
func Test_checkResponseJSONEOF(t *testing.T) {
	withHTTPResponse(t, http.StatusInternalServerError, jsonContentType, nil, func(t *testing.T, res *http.Response) {
		if err := checkResponse(res); err != io.EOF {
			t.Fatal(err)
		}
	})
}

// Test_checkResponseEOF verifies that checkResponse returns an io.ErrUnexpectedEOF
// when a short JSON body is found in the HTTP response body.
func Test_checkResponseJSONUnexpectedEOF(t *testing.T) {
	withHTTPResponse(t, http.StatusInternalServerError, jsonContentType, []byte("{"), func(t *testing.T, res *http.Response) {
		if err := checkResponse(res); err != io.ErrUnexpectedEOF {
			t.Fatal(err)
		}
	})
}

// Test_checkResponseEOF verifies that checkResponse returns the appropriate error
// assuming all sanity checks pass, but the API did return a client-consumable error.
func Test_checkResponseErrorOK(t *testing.T) {
	withHTTPResponse(t, http.StatusInternalServerError, jsonContentType, apiErrJSON, func(t *testing.T, res *http.Response) {
		apiErr := &Error{
			Code:              500,
			Detail:            "The user has not authorized this application or the token is invalid.",
			Type:              "invalid_auth",
			DeveloperFriendly: "The user has not authorized this application or the token is invalid.",
			Duration:          time.Duration(0 * time.Second),
		}

		if err := checkResponse(res); err.Error() != apiErr.Error() {
			t.Fatalf("unexpected API error: %v != %v", err, apiErr)
		}
	})
}

// Test_checkResponseEOF verifies that checkResponse returns no error when HTTP
// status is OK, but response body is empty.
func Test_checkResponseOKNoBody(t *testing.T) {
	withHTTPResponse(t, http.StatusOK, jsonContentType, nil, func(t *testing.T, res *http.Response) {
		if err := checkResponse(res); err != nil {
			t.Fatal(err)
		}
	})
}

// Test_checkResponseEOF verifies that checkResponse returns no error when HTTP
// status is OK, but response body contains data.
func Test_checkResponseOKWithBody(t *testing.T) {
	withHTTPResponse(t, http.StatusOK, jsonContentType, []byte("{}"), func(t *testing.T, res *http.Response) {
		if err := checkResponse(res); err != nil {
			t.Fatal(err)
		}
	})
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

// withHTTPResponse is a test helper which generates a *http.Response and invokes
// an input closure, used for testing.
func withHTTPResponse(t *testing.T, code int, contentType string, body []byte, fn func(t *testing.T, res *http.Response)) {
	res := &http.Response{
		StatusCode: code,
		Header: http.Header{
			"Content-Type": []string{contentType},
		},
		Body: ioutil.NopCloser(bytes.NewReader(body)),
	}

	fn(t, res)
}

// testClient wires up a new Client with a HTTP test server, allowing for easy
// setup and teardown of repetitive code.  The input closure is invoked in the
// HTTP server, to change the functionality as needed for each test.
func testClient(t *testing.T, fn func(t *testing.T, w http.ResponseWriter, r *http.Request)) (*Client, func()) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", jsonContentType)

		if fn != nil {
			fn(t, w, r)
		}
	}))

	client, err := NewClient("foo", "bar", nil)
	if err != nil {
		t.Fatal(err)
	}

	u, err := url.Parse(srv.URL + "/v4")
	if err != nil {
		t.Fatal(err)
	}

	client.url = u

	return client, func() {
		srv.Close()
	}
}

// assertInvalidUserErr asserts that an input error was generated from the
// invalidUserErrJSON used in some tests.
func assertInvalidUserErr(t *testing.T, err error) {
	if err == nil {
		t.Fatal("error should have occurred, but error is nil")
	}

	uErr, ok := err.(*Error)
	if !ok {
		t.Fatal("error is not of type *Error")
	}

	if c := uErr.Code; c != http.StatusNotFound {
		t.Fatalf("unexpected error code: %d != %d", c, http.StatusNotFound)
	}
	detail := "Invalid user."
	if d := uErr.Detail; d != detail {
		t.Fatalf("unexpected error detail: %q != %q", d, detail)
	}
	eType := "invalid_user"
	if e := uErr.Type; e != eType {
		t.Fatalf("unexpected error type: %q != %q", e, eType)
	}
}

// JSON taken from Untappd APIv4 documentation: https://untappd.com/api/docs
var apiErrJSON = []byte(`{
  "meta": {
    "code": 500,
    "error_detail": "The user has not authorized this application or the token is invalid.",
    "error_type": "invalid_auth",
    "developer_friendly": "The user has not authorized this application or the token is invalid.",
    "response_time": {
      "time": 0,
      "measure": "seconds"
    }
  }
}`)

// invalidUserErrJSON is canned JSON used to test for invalid user handling
var invalidUserErrJSON = []byte(`{"meta":{"code":404,"error_detail":"Invalid user.","error_type":"invalid_user","response_time":{"time":0,"measure":"seconds"}}}`)
