package untappd

import (
	"math"
	"net/http"
	"net/url"
	"strconv"
)

// Distance is a distance unit accepted by the Untappd APIv4.
// A set of Distance constants are provided for ease of use.
type Distance string

const (
	// DistanceMiles requests a radius in miles for local checkins.
	DistanceMiles Distance = "m"

	// DistanceKilometers requests a radius in kilometers for local checkins.
	DistanceKilometers Distance = "km"
)

// LocalService is a "service" which allows access to API methods involving checkins
// in a localized area.
type LocalService struct {
	client *Client
}

// Checkins queries for information about checkins in a local area, specified
// by latitude and longitude.
//
// This method returns up to 25 of a local area's most recent checkins.
// For more granular control, and to page through the checkins list using ID
// parameters, use CheckinsMinMaxIDLimit instead.
func (l *LocalService) Checkins(latitude float64, longitude float64) ([]*Checkin, *http.Response, error) {
	return l.CheckinsMinMaxIDLimitRadius(latitude, longitude, 0, math.MaxInt32, 25, 25, DistanceMiles)
}

// CheckinsMinMaxIDLimit queries for information about a Local's checkins,
// but also accepts minimum checkin ID, maximum checkin ID, and a limit
// parameter to enable paging through checkins. The ID parameter
// specifies the Local ID, which will return a list of recent checkins
// for a given Local.
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
	return getCheckins(l.client, "thepub/local", url.Values{
		"lat":       []string{strconv.FormatFloat(latitude, 'f', -1, 64)},
		"lng":       []string{strconv.FormatFloat(longitude, 'f', -1, 64)},
		"min_id":    []string{strconv.Itoa(minID)},
		"max_id":    []string{strconv.Itoa(maxID)},
		"limit":     []string{strconv.Itoa(limit)},
		"radius":    []string{strconv.Itoa(radius)},
		"dist_pref": []string{string(units)},
	})
}
