package untappd

import (
	"net/http"
	"net/url"
	"strconv"
)

// Badges queries for information about a User's badges.  The username
// parameter specifies the User whose badges will be returned.
//
// This method returns up to 50 of the User's most recently earned badges.
// For more granular control, and to page through the badges list, use
// BadgesOffsetLimit instead.
func (u *UserService) Badges(username string) ([]*Badge, *http.Response, error) {
	// Use default parameters as specified by API
	return u.BadgesOffsetLimit(username, 0, 50)
}

// BadgesOffsetLimit queries for information about a User's badges, but also
// accepts offset and limit parameters to enable paging through more than 50
// badges.  The username parameter specifies the User whose badges will be
// returned.
//
// 50 badges is the maximum number of badges which may be returned by one call.
func (u *UserService) BadgesOffsetLimit(username string, offset int, limit int) ([]*Badge, *http.Response, error) {
	q := url.Values{
		"offset": []string{strconv.Itoa(offset)},
		"limit":  []string{strconv.Itoa(limit)},
	}

	// Temporary struct to unmarshal badges JSON
	var v struct {
		Response struct {
			Count int         `json:"count"`
			Items []*rawBadge `json:"items"`
		} `json:"response"`
	}

	// Perform request for user badges by username
	res, err := u.client.request("GET", "user/badges/"+username, nil, q, &v)
	if err != nil {
		return nil, res, err
	}

	// Build result slice from struct
	badges := make([]*Badge, v.Response.Count)
	for i := range v.Response.Items {
		badges[i] = v.Response.Items[i].export()
	}

	return badges, res, nil
}
