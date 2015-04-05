package untappd

import (
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// Beer represents an Untappd beer, and contains information regarding its
// name, style, description, ratings, and other various metadata.
//
// If available, a beer's brewery information can be accessed via the Brewery
// member.
type Beer struct {
	// Metadata from Untappd.
	ID          int
	Name        string
	Label       url.URL
	ABV         float64
	IBU         int
	Slug        string
	Style       string
	Description string

	// Time when this beer was added to Untappd.
	Created time.Time

	// Is this beer present in the specified user's wish list?
	WishList bool

	// Global Untappd rating for this beer.
	OverallRating float64

	// If applicable, the specified user's rating for this beer.
	UserRating float64

	// If applicable, time when the specified user first checked in
	// this beer.
	FirstHad time.Time

	// If applicable, time when the specified user added this beer to
	// their wish list.
	WishListed time.Time

	// If applicable, number of times the specified user has checked
	// in this beer.
	Count int

	// If available, information regarding the brewery which created
	// this beer.
	Brewery *Brewery
}

// BeerService is a "service" which allows access to API methods involving beers.
type BeerService struct {
	client *Client
}

// Info queries for information about a Beer with the specified ID.
// If the compact parameter is set to 'true', only basic beer information will
// be populated.
func (b *BeerService) Info(id int, compact bool) (*Beer, *http.Response, error) {
	// Determine if a compact response is requested
	q := url.Values{}
	if compact {
		q.Set("compact", "true")
	}

	// Temporary struct to unmarshal raw user JSON
	var v struct {
		Response struct {
			Beer rawBeer `json:"beer"`
		} `json:"response"`
	}

	// Perform request for user information by username
	res, err := b.client.request("GET", "beer/info/"+strconv.Itoa(id), q, &v)
	if err != nil {
		return nil, res, err
	}

	// Return results
	return v.Response.Beer.export(), res, nil
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

	// For /v4/beer/info/ID, brewery is located inside the rawBeer struct.
	// This is not the case with /v4/user/beers/username, where it is
	// added by the client method.
	Brewery *rawBrewery `json:"brewery"`
}

// export creates an exported Beer from a rawBeer struct, allowing for more
// useful structures to be created for client consumption.
func (r *rawBeer) export() *Beer {
	b := &Beer{
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

	// If brewery was present inside the Beer struct, as is the case
	// with /v4/beer/info/ID, add it now.
	if r.Brewery != nil {
		b.Brewery = r.Brewery.export()
	}

	return b
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
