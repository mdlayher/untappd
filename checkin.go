package untappd

import (
	"time"
)

// Checkin represents an Untappd checkin, and contains metadata regarding the
// checkin, including the checkin ID, comment, when the checkin occurred, and
// information about the user, beer, and brewery for a given checkin.
type Checkin struct {
	// Metadata from Untappd.
	ID int

	// Time when this checkin was added to Untappd.
	Created time.Time

	// User comment for this checkin.  May be blank.
	Comment string

	// If applicable, the specified user's rating for this beer.
	UserRating float64

	// The user checking in.
	User *User

	// The checkin beer.
	Beer *Beer

	// If available, information regarding the brewery which created
	// this beer.
	Brewery *Brewery
}
