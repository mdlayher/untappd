package untappd

import (
	"math"
	"net/http"
	"net/url"
	"strconv"
)

// Checkins queries for information about checkins in a local area, specified
// by latitude and longitude.
//
// This method returns up to 25 of a local area's most recent checkins.
// For more granular control, and to page through the checkins list using ID
// parameters, use CheckinsMinMaxIDLimitRadius instead.
func (l *LocalService) Checkins(latitude float64, longitude float64) ([]*Checkin, *http.Response, error) {
	return l.CheckinsMinMaxIDLimitRadius(latitude, longitude, 0, math.MaxInt32, 25, 25, DistanceMiles)
}

// CheckinsMinMaxIDLimitRadius queries for information about a local area's
// checkins, but also accepts a variety of parameters to query and page
// through checkins.  The latitude and longitude parameters specify the
// local area where recent checkins will be queried.
//
// 25 checkins is the maximum number of checkins which may be returned by
// one call.
func (l *LocalService) CheckinsMinMaxIDLimitRadius(
	latitude float64,
	longitude float64,
	minID int,
	maxID int,
	limit int,
	radius int,
	units Distance,
) ([]*Checkin, *http.Response, error) {
	return l.client.getCheckins("thepub/local", url.Values{
		"lat":       []string{strconv.FormatFloat(latitude, 'f', -1, 64)},
		"lng":       []string{strconv.FormatFloat(longitude, 'f', -1, 64)},
		"min_id":    []string{strconv.Itoa(minID)},
		"max_id":    []string{strconv.Itoa(maxID)},
		"limit":     []string{strconv.Itoa(limit)},
		"radius":    []string{strconv.Itoa(radius)},
		"dist_pref": []string{string(units)},
	})
}
