package untappd

import (
	"math"
	"net/http"
	"net/url"
	"strconv"
)

// Checkins queries for information about a Venue's checkins.
// The ID parameter specifies the Venue ID, which will return
// a list of recent checkins for a given Venue.
//
// This method returns up to 25 of the Venue's most recent checkins.
// For more granular control, and to page through the checkins list using ID
// parameters, use CheckinsMinMaxIDLimit instead.
func (v *VenueService) Checkins(id int) ([]*Checkin, *http.Response, error) {
	// Use default parameters as specified by API.  Max ID is somewhat
	// arbitrary, but should provide plenty of headroom, just in case.
	return v.CheckinsMinMaxIDLimit(id, 0, math.MaxInt32, 25)
}

// CheckinsMinMaxIDLimit queries for information about a Venue's checkins,
// but also accepts minimum checkin ID, maximum checkin ID, and a limit
// parameter to enable paging through checkins. The ID parameter
// specifies the Venue ID, which will return a list of recent checkins
// for a given Venue.
//
// 25 checkins is the maximum number of checkins which may be returned by
// one call.
func (v *VenueService) CheckinsMinMaxIDLimit(id int, minID int, maxID int, limit int) ([]*Checkin, *http.Response, error) {
	return v.client.getCheckins("venue/checkins/"+strconv.Itoa(id), url.Values{
		"min_id": []string{strconv.Itoa(minID)},
		"max_id": []string{strconv.Itoa(maxID)},
		"limit":  []string{strconv.Itoa(limit)},
	})
}
