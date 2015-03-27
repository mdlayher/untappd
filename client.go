// Package untappd provides an Untappd APIv4 client, written in Go.  MIT Licensed.
//
// To use this client with the Untappd APIv4, you must register for an API key
// here: https://untappd.com/api/register.

// This package is inspired by Google's go-github library, as well as
// Antoine Grondin's canlii library.  Both can be found on GitHub:
//  - https://github.com/google/go-github
//  - https://github.com/aybabtme/canlii
package untappd

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const (
	// jsonContentType is the content type for JSON data
	jsonContentType = "application/json"

	// untappdUserAgent is the user agent this package will rpoert to
	// the Untappd APIv4.
	untappdUserAgent = "github.com/mdlayher/untappd"
)

var (
	// ErrNoClientID is returned when an empty Client ID is passed to NewClient.
	ErrNoClientID = errors.New("no client ID")

	// ErrNoClientSecret is returned when an empty Client Secret is passed
	// to NewClient.
	ErrNoClientSecret = errors.New("no client secret")
)

// Client is a HTTP client for the Untappd APIv4.  It enables access to various
// methods of the Untappd APIv4.
type Client struct {
	client *http.Client
	url    *url.URL

	clientID     string
	clientSecret string

	userAgent string

	// Methods involving a User
	User interface {
		Info(username string, compact bool) (*User, *http.Response, error)
		Friends(username string) ([]*User, *http.Response, error)
		FriendsOffsetLimit(username string, offset uint, limit uint) ([]*User, *http.Response, error)
		Badges(username string) ([]*Badge, *http.Response, error)
		BadgesOffset(username string, offset uint) ([]*Badge, *http.Response, error)
	}
}

// NewClient creates a properly initialized instance of Client, using the input
// client ID, client secret, and http.Client.
//
// To use a Client with the Untappd APIv4, you must register for an API key
// here: https://untappd.com/api/register.
func NewClient(clientID string, clientSecret string, client *http.Client) (*Client, error) {
	// If input client is nil, use http.DefaultClient
	if client == nil {
		client = http.DefaultClient
	}

	// Disallow empty ID and secret
	if clientID == "" {
		return nil, ErrNoClientID
	}
	if clientSecret == "" {
		return nil, ErrNoClientSecret
	}

	// Set up basic client
	c := &Client{
		client: client,
		url: &url.URL{
			Scheme: "https",
			Host:   "api.untappd.com",
			Path:   "v4",
		},

		clientID:     clientID,
		clientSecret: clientSecret,

		userAgent: untappdUserAgent,
	}

	// Add "services" which allow access to various API methods
	c.User = &UserService{client: c}

	return c, nil
}

// Error represents an error returned from the Untappd APIv4.
type Error struct {
	Code              int
	Detail            string
	Type              string
	DeveloperFriendly string
	Duration          time.Duration
}

// Error returns the string representation of an Error.
func (e Error) Error() string {
	// Per APIv4 documentation, the "developer friendly" string should be used
	// in place of the regular "details" string wherever available
	details := e.Detail
	if e.DeveloperFriendly != "" {
		details = e.DeveloperFriendly
	}

	return fmt.Sprintf("%d [%s]: %s", e.Code, e.Type, details)
}

// request creates a new HTTP request, using the specified HTTP method and API endpoint.
func (c *Client) request(method string, endpoint string, query url.Values, v interface{}) (*http.Response, error) {
	// Generate relative URL using API root and endpoint
	rel, err := url.Parse(fmt.Sprintf("%s/%s/", c.url.Path, endpoint))
	if err != nil {
		return nil, err
	}

	// Resolve relative URL to base, using input host
	u := c.url.ResolveReference(rel)

	// Add any URL requested URL query parameters
	q := u.Query()
	for k, v := range query {
		for _, vv := range v {
			q.Add(k, vv)
		}
	}

	// Add required client ID and client secret
	q.Set("client_id", c.clientID)
	q.Set("client_secret", c.clientSecret)
	u.RawQuery = q.Encode()

	// Generate new HTTP request for appropriate URL
	req, err := http.NewRequest(method, u.String(), nil)
	if err != nil {
		return nil, err
	}

	// Set headers to indicate proper content type
	req.Header.Add("Accept", jsonContentType)
	req.Header.Add("Content-Type", jsonContentType)

	// Identify the client
	req.Header.Add("User-Agent", c.userAgent)

	// Invoke request using underlying HTTP client
	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Check response for errors
	if err := checkResponse(res); err != nil {
		return res, err
	}

	// If no second parameter was passed, do not attempt to handle response
	if v == nil {
		return res, nil
	}

	// Decode response body into v, returning response
	return res, json.NewDecoder(res.Body).Decode(v)
}

// checkResponse checks for a non-200 HTTP status code, and returns any errors
// encountered.
func checkResponse(res *http.Response) error {
	// Ensure correct content type
	if cType := res.Header.Get("Content-Type"); cType != jsonContentType {
		return fmt.Errorf("expected %s content type, but received %s", jsonContentType, cType)
	}

	// Check for 200-range status code
	if c := res.StatusCode; 200 <= c && c <= 299 {
		return nil
	}

	// Used as an intermediary form, but the contents are packed into
	// a more consumable form on error output
	var apiErr struct {
		Meta struct {
			Code              int    `json:"code"`
			ErrorDetail       string `json:"error_detail"`
			ErrorType         string `json:"error_type"`
			DeveloperFriendly string `json:"developer_friendly"`
			ResponseTime      struct {
				Time    float64 `json:"time"`
				Measure string  `json:"measure"`
			} `json:"response_time"`
		} `json:"meta"`
	}

	// Unmarshal error response
	if err := json.NewDecoder(res.Body).Decode(&apiErr); err != nil {
		return err
	}

	// Assemble Error struct from API response
	m := apiErr.Meta
	return &Error{
		Code:              m.Code,
		Detail:            m.ErrorDetail,
		Type:              m.ErrorType,
		DeveloperFriendly: m.DeveloperFriendly,
		Duration: timeUnitToDuration(
			m.ResponseTime.Time,
			m.ResponseTime.Measure,
		),
	}
}

// timeUnitToDuration parses a time float64 and measure string from the Untappd
// APIv4, and converts them into a native Go time.Duration.
func timeUnitToDuration(timeFloat float64, measure string) time.Duration {
	// Known measure strings mapped to Go parse-able equivalents
	timeUnits := map[string]string{
		"milliseconds": "ms",
		"seconds":      "s",
		"minutes":      "m",
	}

	// Parse a Go time.Duration from string
	duration, err := time.ParseDuration(fmt.Sprintf("%f%s", timeFloat, timeUnits[measure]))
	if err != nil {
		// If error, return no duration
		return 0
	}

	return duration
}
