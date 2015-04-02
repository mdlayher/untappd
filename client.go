// Package untappd provides an Untappd APIv4 client, written in Go.  MIT Licensed.
//
// To use this client with the Untappd APIv4, you must register for an API key
// here: https://untappd.com/api/register.
//
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

	// untappdUserAgent is the default user agent this package will report to
	// the Untappd APIv4.
	untappdUserAgent = "github.com/mdlayher/untappd"
)

// Sort is a sorting method accepted by the Untappd APIv4.
// A set of Sort constants are provided for ease of use.
type Sort string

// Constants that define various methods that the Untappd APIv4 can use to
// sort Beer results.
const (
	// SortDate sorts a list of beers by most recent date checked in.
	SortDate Sort = "date"

	// SortCheckin sorts a list of beers by highest number of checkins.
	SortCheckin Sort = "checkin"

	// SortHighestRated sorts a list of beers by highest rated overall on Untappd.
	SortHighestRated Sort = "highest_rated"

	// SortLowestRated sorts a list of beers by lowest rated overall on Untappd.
	SortLowestRated Sort = "lowest_rated"

	// SortUserHighestRated sorts a list of beers by highest rated by this user
	// on Untappd.
	SortUserHighestRated Sort = "highest_rated_you"

	// SortUserLowestRated sorts a list of beers by lowest rated by this user
	// on Untappd.
	SortUserLowestRated Sort = "lowest_rated_you"

	// SortHighestABV sorts a list of beers by highest alcohol by volume on Untappd.
	SortHighestABV = "highest_abv"

	// SortLowestABV sorts a list of beers by lowest alcohol by volume on Untappd.
	SortLowestABV = "lowest_abv"
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
	UserAgent string

	client *http.Client
	url    *url.URL

	clientID     string
	clientSecret string

	// Methods involving a User
	User interface {
		// https://untappd.com/api/docs#userinfo
		Info(username string, compact bool) (*User, *http.Response, error)

		// https://untappd.com/api/docs#userfriends
		Friends(username string) ([]*User, *http.Response, error)
		FriendsOffsetLimit(username string, offset int, limit int) ([]*User, *http.Response, error)

		// https://untappd.com/api/docs#userbadges
		Badges(username string) ([]*Badge, *http.Response, error)
		BadgesOffset(username string, offset int) ([]*Badge, *http.Response, error)

		// https://untappd.com/api/docs#userwishlist
		WishList(username string) ([]*Beer, *http.Response, error)
		WishListOffsetLimitSort(username string, offset int, limit int, sort Sort) ([]*Beer, *http.Response, error)

		// https://untappd.com/api/docs#userbeers
		Beers(username string) ([]*Beer, *http.Response, error)
		BeersOffsetLimitSort(username string, offset int, limit int, sort Sort) ([]*Beer, *http.Response, error)
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
		UserAgent: untappdUserAgent,

		client: client,
		url: &url.URL{
			Scheme: "https",
			Host:   "api.untappd.com",
			Path:   "v4",
		},

		clientID:     clientID,
		clientSecret: clientSecret,
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
	req.Header.Add("User-Agent", c.UserAgent)

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
			Code              int              `json:"code"`
			ErrorDetail       string           `json:"error_detail"`
			ErrorType         string           `json:"error_type"`
			DeveloperFriendly string           `json:"developer_friendly"`
			ResponseTime      responseDuration `json:"response_time"`
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
		Duration:          time.Duration(m.ResponseTime),
	}
}
