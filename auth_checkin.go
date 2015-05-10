package untappd

import (
	"net/http"
	"net/url"
	"strconv"
)

// CheckinRequest represents a request to check-in a beer to Untappd.
// To perform a successful checkin, the BeerID, GMTOffset, and TimeZone
// members must be filled in.  The easiest way to obtain the GMTOffset
// and TimeZone for your current system is to use the time package from
// the standard library:
//   beerID := 1
//   timezone, offset := time.Now().Zone()
//   offset = offset / 60 / 60
//
//   request := CheckinRequest{
//       BeerID:    beerID,
//       GMTOffset: offset,
//       TimeZone:  timezone,
//   }
type CheckinRequest struct {
	// Mandatory parameters
	BeerID    int
	GMTOffset int
	TimeZone  string

	// Optional parameters

	// Checkin location
	FoursquareID string
	Latitude     float64
	Longitude    float64

	// User comment and rating
	Comment string
	Rating  float64

	// Send to social media?
	Facebook bool
	Twitter  bool
	// FoursquareID is required if this is true
	Foursquare bool
}

// Checkin checks-in a beer specified by the input CheckinRequest struct.
// A variety of struct members can be filled in to specify the rating,
// comment, etc. for a checkin.
func (a *AuthService) Checkin(r CheckinRequest) (*Checkin, *http.Response, error) {
	// Add required parameters
	q := url.Values{
		"bid":        []string{strconv.Itoa(r.BeerID)},
		"gmt_offset": []string{strconv.Itoa(r.GMTOffset)},
		"timezone":   []string{r.TimeZone},
	}

	// Add optional parameters, if not empty
	if r.FoursquareID != "" {
		q.Set("foursquare_id", r.FoursquareID)
	}
	if r.Latitude != 0 {
		q.Set("geolat", formatFloat(r.Latitude))
	}
	if r.Longitude != 0 {
		q.Set("geolng", formatFloat(r.Longitude))
	}

	if r.Comment != "" {
		q.Set("shout", r.Comment)
	}
	if r.Rating != 0 {
		q.Set("rating", formatFloat(r.Rating))
	}

	if r.Facebook {
		q.Set("facebook", "on")
	}
	if r.Twitter {
		q.Set("twitter", "on")
	}
	if r.Foursquare {
		q.Set("foursquare", "on")
	}

	// Temporary struct to unmarshal checkin JSON
	var v struct {
		Response rawCheckin `json:"response"`
	}

	// Perform request to check in a beer
	res, err := a.client.request("POST", "checkin/add", q, nil, &v)
	if err != nil {
		return nil, res, err
	}

	return v.Response.export(), res, nil
}
