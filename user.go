package untappd

import (
	"net/http"
	"net/url"
)

// User represents an Untappd user.
// TODO(mdlayher): write out fields to access more user information.
type User struct {
	UID       int64  `json:"uid"`
	ID        int64  `json:"id"`
	UserName  string `json:"user_name"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// userService is a "service" which allows access to API methods involving users.
type userService struct {
	client *Client
}

// Info queries for information about a User with the specified username.
// If the compact parameter is set to 'true', only basic user information will
// be populated.
func (u *userService) Info(username string, compact bool) (*User, *http.Response, error) {
	// Determine if a compact response is requested
	q := url.Values{}
	if compact {
		q.Set("compact", "true")
	}

	// Temporary struct to unmarshal user JSON
	var v struct {
		User *User `json:"user"`
	}

	// Perform request for user information by username
	res, err := u.client.request("GET", "user/info/"+username, q, &v)
	if err != nil {
		return nil, res, err
	}

	// Return results
	return v.User, res, nil
}