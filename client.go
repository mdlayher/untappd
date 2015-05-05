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
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	// formEncodedContentType is the content type for key/value POST
	// body requests.
	formEncodedContentType = "application/x-www-form-urlencoded"

	// jsonContentType is the content type for JSON data.
	jsonContentType = "application/json"

	// untappdUserAgent is the default user agent this package will report to
	// the Untappd APIv4.
	untappdUserAgent = "github.com/mdlayher/untappd"
)

var (
	// ErrNoAccessToken is returned when an empty AccessToken is passed to
	// NewAuthenticatedClient.
	ErrNoAccessToken = errors.New("no client ID")

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

	accessToken string

	// Methods which require authentication
	Auth interface {
		// https://untappd.com/api/docs#activityfeed
		Checkins() ([]*Checkin, *http.Response, error)
		CheckinsMinMaxIDLimit(minID int, maxID int, limit int) ([]*Checkin, *http.Response, error)
	}

	// Methods involving a Beer
	Beer interface {
		// https://untappd.com/api/docs#beeractivityfeed
		Checkins(id int) ([]*Checkin, *http.Response, error)
		CheckinsMinMaxIDLimit(id int, minID int, maxID int, limit int) ([]*Checkin, *http.Response, error)

		// https://untappd.com/api/docs#beerinfo
		Info(id int, compact bool) (*Beer, *http.Response, error)

		// https://untappd.com/api/docs#beersearch
		Search(query string) ([]*Beer, *http.Response, error)
		SearchOffsetLimitSort(query string, offset int, limit int, sort Sort) ([]*Beer, *http.Response, error)
	}

	// Methods involving a Brewery
	Brewery interface {
		// https://untappd.com/api/docs#breweryactivityfeed
		Checkins(id int) ([]*Checkin, *http.Response, error)
		CheckinsMinMaxIDLimit(id int, minID int, maxID int, limit int) ([]*Checkin, *http.Response, error)

		// https://untappd.com/api/docs#breweryinfo
		Info(id int, compact bool) (*Brewery, *http.Response, error)

		// https://untappd.com/api/docs#brewerysearch
		Search(query string) ([]*Brewery, *http.Response, error)
		SearchOffsetLimit(query string, offset int, limit int) ([]*Brewery, *http.Response, error)
	}

	// Methods involving a Local area
	Local interface {
		// https://untappd.com/api/docs#theppublocal
		Checkins(latitude float64, longitude float64) ([]*Checkin, *http.Response, error)
		CheckinsMinMaxIDLimitRadius(
			latitude float64,
			longitude float64,
			minID int,
			maxID int,
			limit int,
			radius int,
			units Distance,
		) ([]*Checkin, *http.Response, error)
	}

	// Methods involving a User
	User interface {
		// https://untappd.com/api/docs#userbadges
		Badges(username string) ([]*Badge, *http.Response, error)
		BadgesOffsetLimit(username string, offset int, limit int) ([]*Badge, *http.Response, error)

		// https://untappd.com/api/docs#userbeers
		Beers(username string) ([]*Beer, *http.Response, error)
		BeersOffsetLimitSort(username string, offset int, limit int, sort Sort) ([]*Beer, *http.Response, error)

		// https://untappd.com/api/docs#useractivityfeed
		Checkins(username string) ([]*Checkin, *http.Response, error)
		CheckinsMinMaxIDLimit(username string, minID int, maxID int, limit int) ([]*Checkin, *http.Response, error)

		// https://untappd.com/api/docs#userfriends
		Friends(username string) ([]*User, *http.Response, error)
		FriendsOffsetLimit(username string, offset int, limit int) ([]*User, *http.Response, error)

		// https://untappd.com/api/docs#userinfo
		Info(username string, compact bool) (*User, *http.Response, error)

		// https://untappd.com/api/docs#userwishlist
		WishList(username string) ([]*Beer, *http.Response, error)
		WishListOffsetLimitSort(username string, offset int, limit int, sort Sort) ([]*Beer, *http.Response, error)
	}

	// Methods involving a Venue
	Venue interface {
		// https://untappd.com/api/docs#venueactivityfeed
		Checkins(id int) ([]*Checkin, *http.Response, error)
		CheckinsMinMaxIDLimit(id int, minID int, maxID int, limit int) ([]*Checkin, *http.Response, error)

		// https://untappd.com/api/docs#venueinfo
		Info(id int, compact bool) (*Venue, *http.Response, error)
	}
}

// NewClient creates a properly initialized instance of Client, using the input
// client ID, client secret, and http.Client.
//
// To use a Client with the Untappd APIv4, you must register for an API key
// here: https://untappd.com/api/register.
func NewClient(clientID string, clientSecret string, client *http.Client) (*Client, error) {
	// Disallow empty ID and secret
	if clientID == "" {
		return nil, ErrNoClientID
	}
	if clientSecret == "" {
		return nil, ErrNoClientSecret
	}

	// Perform common client setup
	return newClient(clientID, clientSecret, "", client)
}

// NewAuthenticatedClient creates a properly initialized and authenticated instance
// of Client, using the input access token and http.Client.
//
// NewAuthenticatedClient must be called in order to create a Client which can
// access authenticated API actions, such as checking in beers, toasting other
// users' checkins, adding comments, etc.
//
// To use an authenticated Client with the Untappd APIv4, you must register
// for an API key here: https://untappd.com/api/register.  Next, you must follow
// the OAuth Authentication procedure documented here:
// https://untappd.com/api/docs#authentication.  Upon successful OAuth Authentication,
// you will receive an access token which can be used with NewAuthenticatedClient.
func NewAuthenticatedClient(accessToken string, client *http.Client) (*Client, error) {
	// Disallow empty access token
	if accessToken == "" {
		return nil, ErrNoAccessToken
	}

	// Perform common client setup
	return newClient("", "", accessToken, client)
}

// newClient handles common setup logic for a Client for NewClient and
// NewAuthenticatedClient.
func newClient(clientID string, clientSecret string, accessToken string, client *http.Client) (*Client, error) {
	// If input client is nil, use http.DefaultClient
	if client == nil {
		client = http.DefaultClient
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

		accessToken: accessToken,
	}

	// Add "services" which allow access to various API methods
	c.Auth = &AuthService{client: c}
	c.User = &UserService{client: c}
	c.Beer = &BeerService{client: c}
	c.Brewery = &BreweryService{client: c}
	c.Venue = &VenueService{client: c}
	c.Local = &LocalService{client: c}

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
// Additionally, it accepts POST body parameters, GET query parameters, and an
// optional struct which can be used to unmarshal result JSON.
func (c *Client) request(method string, endpoint string, body url.Values, query url.Values, v interface{}) (*http.Response, error) {
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

	// Always prefer authenticated client access, using an access token.
	// If no token is found, fall back to unauthenticated client ID and
	// client secret.
	if c.accessToken != "" {
		q.Set("access_token", c.accessToken)
	} else {
		q.Set("client_id", c.clientID)
		q.Set("client_secret", c.clientSecret)
	}
	u.RawQuery = q.Encode()

	// Determine if request will contain a POST body
	hasBody := method == "POST" && len(body) > 0
	var length int

	// If performing a POST request and body parameters exist, encode
	// them now
	buf := bytes.NewBuffer(nil)
	if hasBody {
		// Encode and retrieve length to send to server
		buf = bytes.NewBufferString(body.Encode())
		length = buf.Len()
	}

	// Generate new HTTP request for appropriate URL
	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	// Set headers to indicate proper content type
	req.Header.Add("Accept", jsonContentType)

	// For POST requests, add proper headers
	if hasBody {
		req.Header.Add("Content-Type", formEncodedContentType)
		req.Header.Add("Content-Length", strconv.Itoa(length))
	}

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

// getCheckins is the backing method for both any request which returns a
// list of checkins.  It handles performing the necessary HTTP request
// with the correct parameters, and returns a list of Checkins.
func (c *Client) getCheckins(endpoint string, q url.Values) ([]*Checkin, *http.Response, error) {
	// Temporary struct to unmarshal checkin JSON
	var v struct {
		Response struct {
			Checkins struct {
				Count int           `json:"count"`
				Items []*rawCheckin `json:"items"`
			} `json:"checkins"`
		} `json:"response"`
	}

	// Perform request for user checkins by ID
	res, err := c.request("GET", endpoint, nil, q, &v)
	if err != nil {
		return nil, res, err
	}

	// Build result slice from struct
	checkins := make([]*Checkin, v.Response.Checkins.Count)
	for i := range v.Response.Checkins.Items {
		checkins[i] = v.Response.Checkins.Items[i].export()
	}

	return checkins, res, nil
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
