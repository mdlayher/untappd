package untappd

import (
	"net/http"
	"net/url"
	"strconv"
)

// LocalCheckinsRequest represents a request to view checkins in a local area,
// specified by latitude and longitude.  All other parameters are optional,
// but may be used to filter checkins which meet a set of criteria.
type LocalCheckinsRequest struct {
	// Mandatory parameters
	Latitude  float64
	Longitude float64

	// Optional parameters

	// Minimum and maximum checkin IDs to query
	MinID int
	MaxID int

	// Maximum number of results to return
	Limit int

	// Distance radius from latitude/longitude pair, and units
	// for the radius
	Radius int
	Units  Distance
}

// Checkins queries for information about checkins in a local area, specified
// by latitude and longitude.
//
// This method returns up to 25 of a local area's most recent checkins within
// a distance of 25 miles.
// For more granular control, and to page through the checkins list using ID
// parameters, use CheckinsMinMaxIDLimitRadius instead.
func (l *LocalService) Checkins(latitude float64, longitude float64) ([]*Checkin, *http.Response, error) {
	return l.CheckinsMinMaxIDLimitRadius(LocalCheckinsRequest{
		Latitude:  latitude,
		Longitude: longitude,

		Limit: 25,

		Radius: 25,
		Units:  DistanceMiles,
	})
}

// CheckinsMinMaxIDLimitRadius queries for information about a local area's
// checkins, but also accepts a variety of parameters to query and page
// through checkins.  The latitude and longitude parameters specify the
// local area where recent checkins will be queried.
//
// 25 checkins is the maximum number of checkins which may be returned by
// one call.
func (l *LocalService) CheckinsMinMaxIDLimitRadius(r LocalCheckinsRequest) ([]*Checkin, *http.Response, error) {
	// Add required parameters
	q := url.Values{
		"lat": []string{FormatFloat(r.Latitude)},
		"lng": []string{FormatFloat(r.Longitude)},
	}

	// Add optional parameters, if not empty
	if r.MinID != 0 {
		q.Set("min_id", strconv.Itoa(r.MinID))
	}
	if r.MaxID != 0 {
		q.Set("max_id", strconv.Itoa(r.MaxID))
	}

	if r.Limit != 0 {
		q.Set("limit", strconv.Itoa(r.Limit))
	}

	if r.Radius != 0 {
		q.Set("radius", strconv.Itoa(r.Radius))
	}
	if r.Units != "" {
		q.Set("dist_pref", string(r.Units))
	}

	return l.client.getCheckins("thepub/local", q)
}
