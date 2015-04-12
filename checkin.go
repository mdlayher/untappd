package untappd

import (
	"net/url"
	"time"
)

// Checkin represents an Untappd checkin, and contains information regarding its
// name, style, description, ratings, and other various metadata.
//
// If available, a beer's brewery information can be accessed via the Brewery
// member.
type Checkin struct {
	// Metadata from Untappd.
	ID int

	// Time when this checkin was added to Untappd.
	Created time.Time

	Comment string

	// If applicable, the specified user's rating for this beer.
	UserRating float64

	// The user checking in
	User *User

	// The checkin beer
	Beer *Beer

	// If available, information regarding the brewery which created
	// this beer.
	Brewery *Brewery
}

// rawCheckinBeer is the raw JSON representation of an Untappd beer as used when
// describing a checkin.  Its data is unmarshaled from JSON and then exported
// to a Beer struct.
type rawCheckinBeer struct {
	ID       int         `json:"bid"`
	Name     string      `json:"beer_name"`
	Label    responseURL `json:"beer_label"`
	ABV      float64     `json:"beer_abv"`
	Style    string      `json:"beer_style"`
	WishList bool        `json:"wish_list"`
}

// export creates an exported Beer from a rawCheckinBeer struct, allowing for more
// useful structures to be created for client consumption.
func (r *rawCheckinBeer) export() *Beer {
	b := &Beer{
		ID:       r.ID,
		Name:     r.Name,
		Label:    url.URL(r.Label),
		ABV:      r.ABV,
		Style:    r.Style,
		WishList: r.WishList,
	}

	return b
}
