package untappd

import (
	"net/url"
	"time"
)

// Beer represents an Untappd beer, and contains information regarding its
// name, style, description, ratings, and other various metadata.
//
// If available, a beer's brewery information can be accessed via the Brewery
// member.
type Beer struct {
	ID            int
	Name          string
	Label         url.URL
	ABV           float64
	IBU           int
	Slug          string
	Style         string
	Description   string
	Created       time.Time
	WishList      bool
	OverallRating float64
	UserRating    float64
	FirstHad      time.Time
	Count         int

	Brewery *Brewery
}

// rawBeer is the raw JSON representation of an Untappd beer.  Its data is
// unmarshaled from JSON and then exported to a Beer struct.
type rawBeer struct {
	ID            int          `json:"bid"`
	Name          string       `json:"beer_name"`
	Label         responseURL  `json:"beer_label"`
	ABV           float64      `json:"beer_abv"`
	IBU           int          `json:"beer_ibu"`
	Slug          string       `json:"beer_slug"`
	Style         string       `json:"beer_style"`
	Description   string       `json:"beer_description"`
	Created       responseTime `json:"created_at"`
	WishList      bool         `json:"wish_list"`
	OverallRating float64      `json:"rating_score"`
}

// export creates an exported Beer from a rawBeer struct, allowing for more
// useful structures to be created for client consumption.
func (r *rawBeer) export() *Beer {
	return &Beer{
		ID:            r.ID,
		Name:          r.Name,
		Label:         url.URL(r.Label),
		ABV:           r.ABV,
		IBU:           r.IBU,
		Slug:          r.Slug,
		Style:         r.Style,
		Description:   r.Description,
		Created:       time.Time(r.Created),
		WishList:      r.WishList,
		OverallRating: r.OverallRating,
	}
}

// Brewery represents an Untappd brewery, and contains information about a
// brewery's name, location, logo, and various other metadata.
type Brewery struct {
	ID       int
	Name     string
	Slug     string
	Label    url.URL
	Country  string
	Active   bool
	Location BreweryLocation
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
	Label    responseURL     `json:"brewery_label"`
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
		Label:    url.URL(r.Label),
		Country:  r.Country,
		Active:   bool(r.Active),
		Location: r.Location,
	}
}
