package untappd

import (
	"math"
	"net/http"
	"net/url"
	"strconv"
)

// Checkins queries for information about a Beer's checkins.
// The ID parameter specifies the Beer ID, which will return
// a list of recent checkins for a given Beer.
//
// This method returns up to 25 of the Beer's most recent checkins.
// For more granular control, and to page through the checkins list using ID
// parameters, use CheckinsMinMaxIDLimit instead.
func (b *BeerService) Checkins(id int) ([]*Checkin, *http.Response, error) {
	// Use default parameters as specified by API.  Max ID is somewhat
	// arbitrary, but should provide plenty of headroom, just in case.
	return b.CheckinsMinMaxIDLimit(id, 0, math.MaxInt32, 25)
}

// CheckinsMinMaxIDLimit queries for information about a Beer's checkins,
// but also accepts minimum checkin ID, maximum checkin ID, and a limit
// parameter to enable paging through checkins. The ID parameter
// specifies the Beer ID, which will return a list of recent checkins
// for a given Beer.
//
// 25 checkins is the maximum number of checkins which may be returned by
// one call.
func (b *BeerService) CheckinsMinMaxIDLimit(id int, minID int, maxID int, limit int) ([]*Checkin, *http.Response, error) {
	q := url.Values{
		"min_id": []string{strconv.Itoa(minID)},
		"max_id": []string{strconv.Itoa(maxID)},
		"limit":  []string{strconv.Itoa(limit)},
	}

	// Temporary struct to unmarshal checkin JSON
	var v struct {
		Response struct {
			Checkins struct {
				Count int           `json:"count"`
				Items []*rawCheckin `json:"items"`
			} `json:"checkins"`
		} `json:"response"`
	}

	// Perform request for user checkins by ID
	res, err := b.client.request("GET", "beer/checkins/"+strconv.Itoa(id), q, &v)
	if err != nil {
		return nil, res, err
	}

	// Build result slice from struct
	checkins := make([]*Checkin, v.Response.Checkins.Count)
	for i := range v.Response.Checkins.Items {
		checkins[i] = v.Response.Checkins.Items[i].export()
	}

	return checkins, res, nil
}
