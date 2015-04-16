package untappd

import (
	"net/http"
	"net/url"
	"strconv"
)

// Info queries for information about a Beer with the specified ID.
// If the compact parameter is set to 'true', only basic beer information will
// be populated.
func (b *BeerService) Info(id int, compact bool) (*Beer, *http.Response, error) {
	// Determine if a compact response is requested
	q := url.Values{}
	if compact {
		q.Set("compact", "true")
	}

	// Temporary struct to unmarshal raw beer JSON
	var v struct {
		Response struct {
			Beer rawBeer `json:"beer"`
		} `json:"response"`
	}

	// Perform request for beer information by ID
	res, err := b.client.request("GET", "beer/info/"+strconv.Itoa(id), q, &v)
	if err != nil {
		return nil, res, err
	}

	return v.Response.Beer.export(), res, nil
}
