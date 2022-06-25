package untappd

import (
	"net/url"
	"time"
)

// BeerService is a "service" which allows access to API methods involving beers.
type BeerService struct {
	client *Client
}

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
	LabelHD     url.URL
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

	// For beer search requests this is the global checkin count, for beer info
	// requests this is the rating count.
	OverallCount int

	// If applicable, the specified user's rating for this beer.
	UserRating float64

	// If applicable, time when the specified user first, or most recently
	// checked in this beer.
	FirstHad  time.Time
	RecentHad time.Time

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

// rawBeer is the raw JSON representation of an Untappd beer.  Its data is
// unmarshaled from JSON and then exported to a Beer struct.
type rawBeer struct {
	ID            int          `json:"bid"`
	Name          string       `json:"beer_name"`
	Label         responseURL  `json:"beer_label"`
	LabelHD       responseURL  `json:"beer_label_hd"`
	ABV           float64      `json:"beer_abv"`
	IBU           int          `json:"beer_ibu"`
	Slug          string       `json:"beer_slug"`
	Style         string       `json:"beer_style"`
	Description   string       `json:"beer_description"`
	Created       responseTime `json:"created_at"`
	WishList      bool         `json:"wish_list"`
	OverallRating float64      `json:"rating_score"`
	OverallCount  int          `json:"rating_count"`

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
		LabelHD:       url.URL(r.LabelHD),
		ABV:           r.ABV,
		IBU:           r.IBU,
		Slug:          r.Slug,
		Style:         r.Style,
		Description:   r.Description,
		Created:       time.Time(r.Created),
		WishList:      r.WishList,
		OverallRating: r.OverallRating,
		OverallCount:  r.OverallCount,
	}

	// If brewery was present inside the Beer struct, as is the case
	// with /v4/beer/info/ID, add it now.
	if r.Brewery != nil {
		b.Brewery = r.Brewery.export()
	}

	return b
}
