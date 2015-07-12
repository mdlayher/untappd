package untappd

import (
	"math"
	"net/http"
	"net/url"
	"strconv"
)

// Checkins queries for information about a User's checkins.
// The username parameter specifies the User whose checkins will be
// returned.
//
// This method returns up to 25 of the User's most recent checkins.
// For more granular control, and to page through the checkins list using ID
// parameters, use CheckinsMinMaxIDLimit instead.
func (u *UserService) Checkins(username string) ([]*Checkin, *http.Response, error) {
	// Use default parameters as specified by API.  Max ID is somewhat
	// arbitrary, but should provide plenty of headroom, just in case.
	return u.CheckinsMinMaxIDLimit(username, 0, math.MaxInt32, 25)
}

// CheckinsMinMaxIDLimit queries for information about a User's checkins,
// but also accepts minimum checkin ID, maximum checkin ID, and a limit
// parameter to enable paging through checkins. The username parameter
// specifies the User whose checkins will be returned.
//
// 50 checkins is the maximum number of checkins which may be returned by
// one call.
func (u *UserService) CheckinsMinMaxIDLimit(username string, minID int, maxID int, limit int) ([]*Checkin, *http.Response, error) {
	v := url.Values{}
	if minID != 0 {
		v.Set("min_id", strconv.Itoa(minID))
	}
	if maxID != math.MaxInt32 {
		v.Set("max_id", strconv.Itoa(maxID))
	}
	v.Set("limit", strconv.Itoa(limit))
	return u.client.getCheckins("user/checkins/"+username, v)
}
