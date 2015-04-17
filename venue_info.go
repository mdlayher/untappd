package untappd

import (
	"net/http"
	"net/url"
	"strconv"
)

// Info queries for information about a Venue with the specified ID.
// If the compact parameter is set to 'true', only basic venue information will
// be populated.
func (b *VenueService) Info(id int, compact bool) (*Venue, *http.Response, error) {
	// Determine if a compact response is requested
	q := url.Values{}
	if compact {
		q.Set("compact", "true")
	}

	// Temporary struct to unmarshal raw user JSON
	var v struct {
		Response struct {
			Venue rawVenue `json:"venue"`
		} `json:"response"`
	}

	// Perform request for venue information by ID
	res, err := b.client.request("GET", "venue/info/"+strconv.Itoa(id), q, &v)
	if err != nil {
		return nil, res, err
	}

	return v.Response.Venue.export(), res, nil
}
