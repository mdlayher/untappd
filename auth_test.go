package untappd_test

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"github.com/mdlayher/untappd"
	"github.com/nelsam/hel/pers"
)

// TestNewAuthHandler verifies that NewAuthHandler returns appropriate errors
// for various types of input parameters.
func TestNewAuthHandler(t *testing.T) {
	const badURL = "http://%20.com"

	var tests = []struct {
		description  string
		clientID     string
		clientSecret string
		redirectURL  string
		err          error
	}{
		{
			description: "no client ID or client secret",
			err:         untappd.ErrNoClientID,
		},
		{
			description:  "no client ID",
			clientSecret: "bar",
			err:          untappd.ErrNoClientID,
		},
		{
			description: "no client secret",
			clientID:    "foo",
			err:         untappd.ErrNoClientSecret,
		},
		{
			description:  "bad redirect URL",
			clientID:     "foo",
			clientSecret: "bar",
			redirectURL:  badURL,
			err: &url.Error{
				Op:  "parse",
				URL: badURL,
			},
		},
		{
			description:  "ok",
			clientID:     "foo",
			clientSecret: "bar",
			redirectURL:  "http://foo.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			t.Parallel()
			if _, _, err := untappd.NewAuthHandler(tt.clientID, tt.clientSecret, tt.redirectURL, nil, nil); err != tt.err {
				// Special case: check for matching type *url.Error
				if reflect.TypeOf(err) == reflect.TypeOf(tt.err) {
					return
				}

				t.Fatalf("unexpected error for test %q: %v != %v", tt.description, err, tt.err)
			}
		})
	}
}

// TestAuthHandlerServeHTTPBadMethod verifies that AuthHandler returns a
// HTTP 405 on non-GET method.
func TestAuthHandlerServeHTTPBadMethod(t *testing.T) {
	url, done := listeningTestAuthHandler(t, "http://foo.com/", nil, newMockHTTPClient())
	defer done()

	res, err := http.Post(url, "", nil)
	if err != nil {
		log.Fatal(err)
	}

	if got, want := res.StatusCode, http.StatusMethodNotAllowed; got != want {
		log.Fatalf("unexpected HTTP status code: %d != %d", got, want)
	}
}

// TestAuthHandlerServeHTTPBadMethod verifies that AuthHandler returns a
// HTTP 401 if no code parameter is passed via query string.
func TestAuthHandlerServeHTTPNoCodeParameter(t *testing.T) {
	url, done := listeningTestAuthHandler(t, "http://foo.com", nil, newMockHTTPClient())
	defer done()

	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	if got, want := res.StatusCode, http.StatusBadRequest; got != want {
		log.Fatalf("unexpected HTTP status code: %d != %d", got, want)
	}
}

// TestAuthHandlerServeHTTPOAuthInternalServerError verifies that AuthHandler
// returns a HTTP 502 if the upstream server returns a non-200 status code.
func TestAuthHandlerServeHTTPOAuthInternalServerError(t *testing.T) {
	testOAuthBadGateway(t, func(t *testing.T, w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	})
}

// TestAuthHandlerServeHTTPOAuthNotJSON verifies that AuthHandler returns a HTTP 502
// if the upstream server returns a non-JSON Content-Type header.
func TestAuthHandlerServeHTTPOAuthNotJSON(t *testing.T) {
	testOAuthBadGateway(t, func(t *testing.T, w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("hello world"))
	})
}

// TestAuthHandlerServeHTTPOAuthBadJSON verifies that AuthHandler returns a HTTP
// 502 if the upstream server returns broken JSON.
func TestAuthHandlerServeHTTPOAuthBadJSON(t *testing.T) {
	testOAuthBadGateway(t, func(t *testing.T, w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", untappd.JSONContentType)
		w.Write([]byte("{"))
	})
}

// TestAuthHandlerServeHTTPOK verifies that AuthHandler can complete an
// entire mock authentication cycle, and return the correct final token upon
// successful authentication.
func TestAuthHandlerServeHTTPOK(t *testing.T) {
	expectedToken := "ABCDEF0123456789"

	client := newMockHTTPClient()
	body := bytes.NewBufferString(`{"response":{"access_token":"` + expectedToken + `"}}`)
	pers.Return(client.GetOutput, &http.Response{
		StatusCode: http.StatusOK,
		Header: http.Header{
			"Content-Type": []string{untappd.JSONContentType},
		},
		Body: ioutil.NopCloser(body),
	}, nil)

	var accessToken string
	tokenFn := func(token string, w http.ResponseWriter, r *http.Request) {
		accessToken = token
	}

	url, done := listeningTestAuthHandler(t, "http://foo.com", tokenFn, client)
	defer done()

	res, err := http.Get(url + "?code=foo")
	if err != nil {
		t.Fatal(err)
	}

	if got, want := res.StatusCode, http.StatusOK; got != want {
		t.Fatalf("unexpected HTTP status code: %d != %d", got, want)
	}

	if got, want := accessToken, expectedToken; got != want {
		t.Fatalf("unexpected access token: %q != %q", got, want)
	}
}

// errBrokenWriter is always returned by brokenResponseWriter.Write.
var errBrokenWriter = errors.New("broken writer")

// brokenResponseWriter is a http.ResponseWriter which always returns
// an error when its Write method is called.
type brokenResponseWriter struct {
	*httptest.ResponseRecorder
}

// Write implements a broken Write method for a brokenResponseWriter.
func (w *brokenResponseWriter) Write(b []byte) (int, error) {
	return 0, errBrokenWriter
}

// Test_defaultTokenFnBadWriter verifies that defaultTokenFn returns an
// internal server error if it is unable to write a response body.
func Test_defaultTokenFnBadWriter(t *testing.T) {
	w := &brokenResponseWriter{
		ResponseRecorder: httptest.NewRecorder(),
	}

	client := newMockHTTPClient()
	body := bytes.NewBufferString(`{"response":{"access_token":"foo"}}`)
	resp := &http.Response{
		StatusCode: 200,
		Header: http.Header{
			"Content-Type": []string{"application/json"},
		},
		Body: ioutil.NopCloser(body),
	}
	pers.Return(client.GetOutput, resp, nil)

	h := newTestAuthHandler(t, "http://foo.com", nil, client)

	req := httptest.NewRequest("GET", "http://bar.com?code=foo", ioutil.NopCloser(bytes.NewBufferString("")))
	h.ServeHTTP(w, req)

	if got, want := w.Code, http.StatusInternalServerError; got != want {
		t.Fatalf("unexpected HTTP status code: %d != %d", got, want)
	}
}

// Test_defaultTokenFnOK verifies that defaultTokenFn returns a valid token
// if it is unable to write a response body.
func Test_defaultTokenFnOK(t *testing.T) {
	for _, tok := range []string{"foo", "bar", "baz"} {
		t.Run(fmt.Sprintf("auth token %s", tok), func(t *testing.T) {
			client := newMockHTTPClient()
			body := bytes.NewBufferString(fmt.Sprintf(`{"response":{"access_token":"%s"}}`, tok))
			resp := &http.Response{
				StatusCode: 200,
				Header: http.Header{
					"Content-Type": []string{"application/json"},
				},
				Body: ioutil.NopCloser(body),
			}
			pers.Return(client.GetOutput, resp, nil)

			h := newTestAuthHandler(t, "http://foo.com", nil, client)

			req := httptest.NewRequest("GET", "http://bar.com?code=foo", ioutil.NopCloser(bytes.NewBufferString("")))

			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)

			if got, want := rec.Code, http.StatusOK; got != want {
				t.Fatalf("unexpected HTTP status code: %d != %d", got, want)
			}
			if b := rec.Body.String(); b != tok {
				t.Fatalf("unexpected response body: %q != %q", b, tok)
			}
		})
	}
}

// listeningTestAuthHandler creates a mocked AuthHandler which points at a httptest server,
// and returns that server's URL and a function to shut it down.
func listeningTestAuthHandler(t *testing.T, redirectURL string, fn untappd.TokenHandlerFunc, client untappd.HTTPClient) (string, func()) {
	t.Helper()
	h := newTestAuthHandler(t, redirectURL, fn, client)

	srv := httptest.NewServer(h)
	return srv.URL, func() {
		srv.Close()
	}
}

// newTestAuthHandler creates an AuthHandler for testing.
func newTestAuthHandler(t *testing.T, redirectURL string, fn untappd.TokenHandlerFunc, client untappd.HTTPClient) *untappd.AuthHandler {
	t.Helper()

	if client == nil {
		t.Fatal("A mocked out HTTP client must be provided to prevent hitting untappd's API directly")
	}
	h, _, err := untappd.NewAuthHandler(
		"foo",
		"bar",
		redirectURL,
		fn,
		client,
	)
	if err != nil {
		t.Fatal(err)
	}

	return h
}

// testOAuthBadGateway handles common setup procedures for tests which check
// for a HTTP 503 Bad Gateway error.
func testOAuthBadGateway(t *testing.T, fn func(t *testing.T, w http.ResponseWriter, r *http.Request)) {
	t.Helper()

	client, done := newMockHTTPClientHandler(t, fn)
	defer done()
	url, done2 := listeningTestAuthHandler(t, "http://foo.com", nil, client)
	defer done2()

	res, err := http.Get(url + "?code=foo")
	if err != nil {
		t.Fatal(err)
	}

	if got, want := res.StatusCode, http.StatusBadGateway; got != want {
		t.Fatalf("unexpected HTTP status code: %d != %d", got, want)
	}
}
