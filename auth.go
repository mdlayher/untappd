package untappd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const (
	// untappdOAuthAuthenticate is a URL pattern which is interpolated with
	// an Untappd APIv4 client ID and redirect URL.  The resulting URL can
	// be returned to clients to begin the authentication process.
	untappdOAuthAuthenticate = "https://untappd.com/oauth/authenticate/?client_id=%s&response_type=code&redirect_url=%s"

	// untappdOAuthAuthorize is a URL pattern which is interpolated with
	// an Untappd APIv4 client ID, client secret, and redirect URL.
	// The resulting URL is requested by AuthHandler.ServeHTTP, generating
	// an Access Token for client consumption.
	untappdOAuthAuthorize = "https://untappd.com/oauth/authorize/?client_id=%s&client_secret=%s&response_type=code&redirect_url=%s"
)

// AuthService is a "service" which allows access to API methods which require
// authentication.
type AuthService struct {
	client *Client
}

// AuthHandler implements http.Handler, and provides a simple process for
// authenticating users using OAuth with Untappd APIv4.
type AuthHandler struct {
	clientID     string
	clientSecret string
	redirectURL  *url.URL
	oAuthURL     *url.URL
	handler      TokenHandlerFunc
	client       HTTPClient
}

// TokenHandlerFunc is a function which is invoked at the end of a successful
// AuthHandler authentication process.  The token generated during the process is
// provided via the token parameter, and the HTTP request and response writers are
// available for further HTTP processing.
type TokenHandlerFunc func(token string, w http.ResponseWriter, r *http.Request)

// defaultTokenFn is the default implementation of TokenHandlerFunc, and is used
// automatically by NewAuthHandler, unless a custom TokenHandlerFunc is provided.
// This function simply prints the token to the HTTP response writer.
var defaultTokenFn = func(token string, w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte(token)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// NewAuthHandler creates a http.Handler which can be used to easily authenticate
// a user using the Server Side Authentication process, documented here:
// https://untappd.com/api/docs#authentication.
//
// The first return parameter is the http.Handler described above.  The second is
// a URL which should be provided to a user, so that they can begin the
// authentication flow.  The third contains any errors which may have occurred
// during setup.
//
// The client ID, client secret, and redirectURL parameters are mandatory.
//
// The TokenHandlerFunc parameter can be used to provide a custom handler which
// contains an access token, and HTTP request and response writers, for further
// processing. The TokenHandlerFunc is only called upon successful authentication.
// Otherwise, an HTTP error is returned to the client.  If no TokenHandlerFunc is
// provided, a default one which writes the token to the HTTP response body will
// be used.
//
// The http.Client parameter can be used to provide a custom http.Client which
// obeys timeouts, etc.  This client is used to communicate with an upstream
// OAuth authentication server.  If no http.Client is provided, http.DefaultClient
// will be used.
func NewAuthHandler(clientID, clientSecret, redirectURL string, fn TokenHandlerFunc, client HTTPClient) (*AuthHandler, *url.URL, error) {
	if clientID == "" {
		return nil, nil, ErrNoClientID
	}
	if clientSecret == "" {
		return nil, nil, ErrNoClientSecret
	}

	ru, err := url.Parse(redirectURL)
	if err != nil {
		return nil, nil, err
	}

	cu, err := url.Parse(fmt.Sprintf(
		untappdOAuthAuthenticate,
		clientID,
		ru.String(),
	))
	if err != nil {
		return nil, nil, err
	}

	ou, err := url.Parse(fmt.Sprintf(
		untappdOAuthAuthorize,
		clientID,
		clientSecret,
		ru.String(),
	))
	if err != nil {
		return nil, nil, err
	}

	if fn == nil {
		fn = defaultTokenFn
	}

	if client == nil {
		client = http.DefaultClient
	}

	return &AuthHandler{
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURL:  ru,
		oAuthURL:     ou,
		handler:      fn,
		client:       client,
	}, cu, nil
}

// ServeHTTP implements http.Handler, and provides a simple http.Handler which
// can properly authenticate using the Server Side Authentication method outlined
// in Untappd documentation: https://untappd.com/api/docs#authentication.
func (a *AuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "only GET requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "no 'code' GET parameter", http.StatusBadRequest)
		return
	}

	res, err := a.client.Get(a.oAuthURL.String() + "&code=" + code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	if c := res.StatusCode; c > 299 || c < 200 {
		http.Error(w, fmt.Sprintf("authentication server error: HTTP %03d", c), http.StatusBadGateway)
		return
	}

	if !strings.Contains(res.Header.Get("Content-Type"), JSONContentType) {
		http.Error(w, "authentication server sent non-JSON content", http.StatusBadGateway)
		return
	}

	var v struct {
		Response struct {
			AccessToken string `json:"access_token"`
		} `json:"response"`
	}

	if err := json.NewDecoder(res.Body).Decode(&v); err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	a.handler(v.Response.AccessToken, w, r)
}
