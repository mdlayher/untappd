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

// rawVenue is the raw JSON representation of an Untappd venue.  Its data is
// unmarshaled from JSON and then exported to a Venue struct.
type rawVenue struct {
	ID       int           `json:"venue_id"`
	Name     string        `json:"venue_name"`
	Updated  responseTime  `json:"last_updated"`
	Category string        `json:"primary_category"`
	Public   bool          `json:"public_venue"`
	Location VenueLocation `json:"location"`
}

// export creates an exported Venue from a rawVenue struct, allowing for
// more useful structures to be created for client consumption.
func (r *rawVenue) export() *Venue {
	return &Venue{
		ID:       r.ID,
		Name:     r.Name,
		Updated:  time.Time(r.Updated),
		Category: r.Category,
		Public:   r.Public,
		Location: r.Location,
	}
}
