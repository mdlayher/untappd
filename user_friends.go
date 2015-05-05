package untappd

import (
	"net/http"
	"net/url"
	"strconv"
)

// Friends queries for information about a User's friends.  The username
// parameter specifies the User whose friends will be returned.
//
// This method returns up to a maximum of 25 friends.  For more granular
// control, and to page through the friends list, use FriendsOffsetLimit
// instead.
//
// The resulting slice of User structs contains a more limited set of user
// information than a call to Info would.  However, basic information such as
// user ID, username, first name, last name, bio, etc. is available.
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
func (u *UserService) FriendsOffsetLimit(username string, offset int, limit int) ([]*User, *http.Response, error) {
	q := url.Values{
		"offset": []string{strconv.Itoa(offset)},
		"limit":  []string{strconv.Itoa(limit)},
	}

	// Temporary struct to unmarshal friends JSON
	var v struct {
		Response struct {
			Count int `json:"count"`
			Items []struct {
				User rawUser `json:"user"`
			} `json:"items"`
		} `json:"response"`
	}

	// Perform request for user friends by username
	res, err := u.client.request("GET", "user/friends/"+username, nil, q, &v)
	if err != nil {
		return nil, res, err
	}

	// Build result slice from struct
	users := make([]*User, v.Response.Count)
	for i := range v.Response.Items {
		users[i] = v.Response.Items[i].User.export()
	}

	return users, res, nil
}
