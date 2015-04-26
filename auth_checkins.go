package untappd

import (
	"math"
	"net/http"
	"net/url"
	"strconv"
)

// Checkins queries for information about checkins from friends of an
// authenticated user.  This is akin to the "Recent Friend Activity" feed
// displayed on the homepage of Untappd for an authenticated user.
//
// This method returns up to 25 of an authenticated user's friends' recent
// checkins.  For more granular control, and to page through the checkins
// list using ID parameters, use CheckinsMinMaxIDLimit instead.
func (a *AuthService) Checkins() ([]*Checkin, *http.Response, error) {
	// Use default parameters as specified by API.  Max ID is somewhat
	// arbitrary, but should provide plenty of headroom, just in case.
	return a.CheckinsMinMaxIDLimit(0, math.MaxInt32, 25)
}

// CheckinsMinMaxIDLimit queries for information about checkins from friends
// of an authenticated user, but also accepts minimum checkin ID, maximum
// checkin ID, and a limit parameter to enable paging through checkins.
// This is akin to the "Recent Friend Activity" feed displayed on the homepage
// of Untappd for an authenticated user.
//
// 50 checkins is the maximum number of checkins which may be returned by
// one call.
func (a *AuthService) CheckinsMinMaxIDLimit(minID int, maxID int, limit int) ([]*Checkin, *http.Response, error) {
	return a.client.getCheckins("checkin/recent", url.Values{
		"min_id": []string{strconv.Itoa(minID)},
		"max_id": []string{strconv.Itoa(maxID)},
		"limit":  []string{strconv.Itoa(limit)},
	})
}
