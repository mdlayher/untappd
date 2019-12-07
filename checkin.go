package untappd

import (
	"net/url"
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

	// If available, information regarding the venue where this checkin
	// occurred.  If a venue was not added to the checkin, this member
	// will be nil.
	Venue *Venue

	// Badges earned when this checkin was submitted.
	Badges []*Badge

	// Toasts by Untappd users for this checkin.
	Toasts []*Toast

	// Comments by Untappd users about this checkin.
	Comments []*Comment

	// Media uploaded by Untappd users about this checkin
	// If the slice has zero length, no media exists for this checkin.
	Media []*CheckinMedia
}

// CheckinMedia contains links to media regarding a Checkin.  Included are links
// to a small, medium, large, and original photos for a given Checkin.
type CheckinMedia struct {
	PhotoID       int
	SmallPhoto    url.URL
	MediumPhoto   url.URL
	LargePhoto    url.URL
	OriginalPhoto url.URL
}

// rawCheckin is the raw JSON representation of an Untappd checkin.  Its data is
// unmarshaled from JSON and then exported to a Checkin struct.
type rawCheckin struct {
	ID         int           `json:"checkin_id"`
	Beer       rawBeer       `json:"beer"`
	Brewery    rawBrewery    `json:"brewery"`
	User       rawUser       `json:"user"`
	Venue      responseVenue `json:"venue"`
	UserRating float64       `json:"rating_score"`
	Comment    string        `json:"checkin_comment"`
	Created    responseTime  `json:"created_at"`

	Badges struct {
		Count int         `json:"count"`
		Items []*rawBadge `json:"items"`
	} `json:"badges"`

	Toasts struct {
		Count int         `json:"count"`
		Items []*rawToast `json:"items"`
	} `json:"toasts"`

	Comments struct {
		Count int           `json:"count"`
		Items []*rawComment `json:"items"`
	} `json:"comments"`

	Media struct {
		Count int                `json:"count"`
		Items []*rawCheckinMedia `json:"items"`
	} `json:"media"`
}

type rawCheckinMedia struct {
	PhotoID int `json:"photo_id"`
	Photo   struct {
		SmallPhoto    responseURL `json:"photo_img_sm"`
		MediumPhoto   responseURL `json:"photo_img_med"`
		LargePhoto    responseURL `json:"photo_img_lg"`
		OriginalPhoto responseURL `json:"photo_img_og"`
	} `json:"photo"`
}

// export creates an exported CheckinMedia from a rawCheckinMedia struct, allowing
// for more useful structures to be created for client consumption.
func (r *rawCheckinMedia) export() *CheckinMedia {
	return &CheckinMedia{
		PhotoID:       r.PhotoID,
		SmallPhoto:    url.URL(r.Photo.SmallPhoto),
		MediumPhoto:   url.URL(r.Photo.MediumPhoto),
		LargePhoto:    url.URL(r.Photo.LargePhoto),
		OriginalPhoto: url.URL(r.Photo.OriginalPhoto),
	}
}

// export creates an exported Checkin from a rawCheckin struct, allowing for more
// useful structures to be created for client consumption.
func (r *rawCheckin) export() *Checkin {
	c := &Checkin{
		ID:         r.ID,
		Comment:    r.Comment,
		UserRating: r.UserRating,
		Created:    time.Time(r.Created),
		Beer:       r.Beer.export(),
		Brewery:    r.Brewery.export(),
		User:       r.User.export(),
	}

	// If no venue was set in the response JSON, venue will be nil
	if r.Venue.ID != 0 && r.Venue.Name != "" {
		// Since venue was not empty, add it to the struct
		rv := rawVenue(r.Venue)
		c.Venue = rv.export()
	}

	badges := make([]*Badge, r.Badges.Count)
	for i := range r.Badges.Items {
		badges[i] = r.Badges.Items[i].export()
	}
	c.Badges = badges

	toasts := make([]*Toast, r.Toasts.Count)
	for i := range r.Toasts.Items {
		toasts[i] = r.Toasts.Items[i].export()
	}
	c.Toasts = toasts

	comments := make([]*Comment, r.Comments.Count)
	for i := range r.Comments.Items {
		comments[i] = r.Comments.Items[i].export()
	}
	c.Comments = comments

	media := make([]*CheckinMedia, r.Media.Count)
	for i := range r.Media.Items {
		media[i] = r.Media.Items[i].export()
	}
	c.Media = media

	return c
}
