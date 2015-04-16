package untappd

import (
	"math"
	"net/http"
	"net/url"
	"strconv"
	"time"
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
	q := url.Values{
		"min_id": []string{strconv.Itoa(minID)},
		"max_id": []string{strconv.Itoa(maxID)},
		"limit":  []string{strconv.Itoa(limit)},
	}

	// Temporary struct to unmarshal checkin JSON
	var v struct {
		Response struct {
			Checkins struct {
				Count int `json:"count"`
				Items []struct {
					ID         int           `json:"checkin_id"`
					Beer       rawBeer       `json:"beer"`
					Brewery    rawBrewery    `json:"brewery"`
					User       rawUser       `json:"user"`
					Venue      responseVenue `json:"venue"`
					UserRating float64       `json:"rating_score"`
					Comment    string        `json:"checkin_comment"`
					Created    responseTime  `json:"created_at"`
				} `json:"items"`
			} `json:"checkins"`
		} `json:"response"`
	}

	// Perform request for user checkins by username
	res, err := u.client.request("GET", "user/checkins/"+username, q, &v)
	if err != nil {
		return nil, res, err
	}

	// Build result slice from struct
	checkins := make([]*Checkin, v.Response.Checkins.Count)
	for i := range v.Response.Checkins.Items {
		// Information about the beer itself
		checkin := &Checkin{
			ID:         v.Response.Checkins.Items[i].ID,
			Comment:    v.Response.Checkins.Items[i].Comment,
			UserRating: v.Response.Checkins.Items[i].UserRating,
			Created:    time.Time(v.Response.Checkins.Items[i].Created),
		}
		checkins[i] = checkin

		checkins[i].Beer = v.Response.Checkins.Items[i].Beer.export()
		checkins[i].Brewery = v.Response.Checkins.Items[i].Brewery.export()
		checkins[i].User = v.Response.Checkins.Items[i].User.export()

		// If no venue was set in the response JSON, venue will be nil
		venue := v.Response.Checkins.Items[i].Venue
		if venue.ID != 0 && venue.Name != "" {
			// Since venue was not empty, add it to the struct
			rv := rawVenue(v.Response.Checkins.Items[i].Venue)
			checkins[i].Venue = rv.export()
		}
	}

	return checkins, res, nil
}
