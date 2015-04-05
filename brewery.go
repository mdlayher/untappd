package untappd

import (
	"net/http"
	"net/url"
	"strconv"
)

// Brewery represents an Untappd brewery, and contains information about a
// brewery's name, location, logo, and various other metadata.
type Brewery struct {
	ID       int
	Name     string
	Slug     string
	Logo     url.URL
	Country  string
	Active   bool
	Location BreweryLocation
}

// BreweryService is a "service" which allows access to API methods involving
// breweries.
type BreweryService struct {
	client *Client
}

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

	// Return results
	return v.Response.Brewery.export(), res, nil
}

// BreweryLocation represent's an Untappd brewery's location, and contains
// information such as the brewery's city, state, and latitude/longitude.
type BreweryLocation struct {
	City      string  `json:"brewery_city"`
	State     string  `json:"brewery_state"`
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lng"`
}

// rawBrewery is the raw JSON representation of an Untappd brewery.  Its data is
// unmarshaled from JSON and then exported to a Brewery struct.
type rawBrewery struct {
	ID       int             `json:"brewery_id"`
	Name     string          `json:"brewery_name"`
	Slug     string          `json:"brewery_slug"`
	Logo     responseURL     `json:"brewery_label"`
	Country  string          `json:"country_name"`
	Active   responseBool    `json:"brewery_active"`
	Location BreweryLocation `json:"location"`
}

// export creates an exported Brewery from a rawBrewery struct, allowing for
// more useful structures to be created for client consumption.
func (r *rawBrewery) export() *Brewery {
	return &Brewery{
		ID:       r.ID,
		Name:     r.Name,
		Slug:     r.Slug,
		Logo:     url.URL(r.Logo),
		Country:  r.Country,
		Active:   bool(r.Active),
		Location: r.Location,
	}
}
