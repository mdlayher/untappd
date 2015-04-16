package untappd

import (
	"net/http"
	"net/url"
	"strconv"
)

// Info queries for information about a Brewery with the specified ID.
// If the compact parameter is set to 'true', only basic brewery information will
// be populated.
func (b *BreweryService) Info(id int, compact bool) (*Brewery, *http.Response, error) {
	// Determine if a compact response is requested
	q := url.Values{}
	if compact {
		q.Set("compact", "true")
	}

	// Temporary struct to unmarshal raw user JSON
	var v struct {
		Response struct {
			Brewery rawBrewery `json:"brewery"`
		} `json:"response"`
	}

	// Perform request for brewery information by ID
	res, err := b.client.request("GET", "brewery/info/"+strconv.Itoa(id), q, &v)
	if err != nil {
		return nil, res, err
	}

	return v.Response.Brewery.export(), res, nil
}
