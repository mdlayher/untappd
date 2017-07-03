package untappd

import (
	"time"
)

// Venue represents an Untappd venue, and contains information about a
// venue's name, location, categories, and various other metadata.
type Venue struct {
	// Metadata from Untappd.
	ID      int
	Name    string
	Updated time.Time

	// Category of thie venue.
	Category string

	// Is this a public venue?
	Public bool

	// Location of this venue.
	Location VenueLocation

	// Foursquare data.
	Foursquare VenueFoursquare

	// Popular beers at this venue.
	TopBeers []*Beer

	// Checkins at this venue.
	Checkins []*Checkin
}

// VenueService is a "service" which allows access to API methods involving
// venues.
type VenueService struct {
	client *Client
}

// VenueLocation represent's an Untappd venue's location, and contains
// information such as the venue's address, city, state, country, and
// latitude/longitude.
type VenueLocation struct {
	Address   string  `json:"venue_address"`
	City      string  `json:"venue_city"`
	State     string  `json:"venue_state"`
	Country   string  `json:"venue_country"`
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lng"`
}

// VenueFoursquare represents an Untappd venue's Foursquare data, and contains
// the venue's Foursquare ID and URL.
type VenueFoursquare struct {
	ID  string `json:"foursquare_id"`
	URL string `json:"foursquare_url"`
}

// rawVenue is the raw JSON representation of an Untappd venue.  Its data is
// unmarshaled from JSON and then exported to a Venue struct.
type rawVenue struct {
	ID         int             `json:"venue_id"`
	Name       string          `json:"venue_name"`
	Updated    responseTime    `json:"last_updated"`
	Category   string          `json:"primary_category"`
	Public     bool            `json:"public_venue"`
	Location   VenueLocation   `json:"location"`
	Foursquare VenueFoursquare `json:"foursquare"`
	TopBeers   struct {
		Offset int `json:"offset"`
		Limit  int `json:"limit"`
		Count  int `json:"count"`
		Items  []struct {
			Created    responseTime `json:"created_at"`
			TotalCount int          `json:"total_count"`
			YourCount  int          `json:"your_count"`

			Beer    rawBeer    `json:"beer"`
			Brewery rawBrewery `json:"brewery"`
		} `json:"items"`
	} `json:"top_beers"`
	Checkins struct {
		Count int           `json:"count"`
		Items []*rawCheckin `json:"items"`
	} `json:"checkins"`
}

// export creates an exported Venue from a rawVenue struct, allowing for
// more useful structures to be created for client consumption.
func (r *rawVenue) export() *Venue {
	beers := make([]*Beer, r.TopBeers.Count)
	for i := range r.TopBeers.Items {
		beers[i] = r.TopBeers.Items[i].Beer.export()
		beers[i].Brewery = r.TopBeers.Items[i].Brewery.export()
	}

	checkins := make([]*Checkin, r.Checkins.Count)
	for i := range r.Checkins.Items {
		checkins[i] = r.Checkins.Items[i].export()
	}

	return &Venue{
		ID:         r.ID,
		Name:       r.Name,
		Updated:    time.Time(r.Updated),
		Category:   r.Category,
		Public:     r.Public,
		Location:   r.Location,
		Foursquare: r.Foursquare,
		TopBeers:   beers,
		Checkins:   checkins,
	}
}
