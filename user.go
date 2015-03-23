package untappd

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

var (
	// errOverLimit is returned when a client attempts to use a limit parameter
	// of greater than 25 in a user friends request.
	errOverLimit = errors.New("limit must not be greater than 25")
)

// User represents an Untappd user.
//
// BUG(mdlayher): write out fields to access more user information.
type User struct {
	UID       int64  `json:"uid"`
	ID        int64  `json:"id"`
	UserName  string `json:"user_name"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// UserService is a "service" which allows access to API methods involving users.
type UserService struct {
	client *Client
}

// Info queries for information about a User with the specified username.
// If the compact parameter is set to 'true', only basic user information will
// be populated.
func (u *UserService) Info(username string, compact bool) (*User, *http.Response, error) {
	// Determine if a compact response is requested
	q := url.Values{}
	if compact {
		q.Set("compact", "true")
	}

	// Temporary struct to unmarshal user JSON
	var v struct {
		// BUG(mdlayher): determine if Meta struct should be returned as well
		Response struct {
			User *User `json:"user"`
		} `json:"response"`
	}

	// Perform request for user information by username
	res, err := u.client.request("GET", "user/info/"+username, q, &v)
	if err != nil {
		return nil, res, err
	}

	// Return results
	return v.Response.User, res, nil
}

// Friends queries for information about a User's friends.  The username
// parameter specifies the User whose friends will be returned.
//
// This method returns up to a maximum of 25 friends.  For more granular
// control, and to page through the friends list, use FriendsOffsetLimit
// instead.
func (u *UserService) Friends(username string) ([]*User, *http.Response, error) {
	// Use default parameters as specified by API
	return u.FriendsOffsetLimit(username, 0, 25)
}

// FriendsOffsetLimit queries for information about a User's friends, but also
// accepts offset and limit parameters to enable paging through more than 25
// friends.  The username parameter specifies the User whose friends will be
// returned.
//
// 25 friends is the maximum number of friends which may be returned by one call.
// Attempts to specify a limit of greater than 25 will result in an error.
func (u *UserService) FriendsOffsetLimit(username string, offset uint, limit uint) ([]*User, *http.Response, error) {
	// API only allows a max of 25 for limit
	// Documentation: https://untappd.com/api/docs#userfriends
	if limit > 25 {
		return nil, nil, errOverLimit
	}

	q := url.Values{
		"offset": []string{strconv.Itoa(int(offset))},
		"limit":  []string{strconv.Itoa(int(limit))},
	}

	// Temporary struct to unmarshal friends JSON
	var v struct {
		// BUG(mdlayher): determine if Meta struct should be returned as well
		Response struct {
			Count int `json:"count"`
			// BUG(mdlayher): Untappd API may return an object here for one
			// result, instead of an array.  This may require further JSON
			// parsing before unmarshaling into a struct.
			Items []struct {
				// BUG(mdlayher): parse more fields later
				User *User `json:"user"`
			} `json:"items"`
		} `json:"response"`
	}

	// Perform request for user friends by username
	res, err := u.client.request("GET", "user/friends/"+username, q, &v)
	if err != nil {
		return nil, res, err
	}

	// Build result slice from struct
	users := make([]*User, v.Response.Count)
	for i := range v.Response.Items {
		users[i] = v.Response.Items[i].User
	}

	// Return results
	return users, res, nil
}
