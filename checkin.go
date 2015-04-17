package untappd

import (
	"net/http"
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
}

// getCheckins is the backing method for both Client.User.Checkins and
// Client.Beer.Checkins.  It handles performing the necessary HTTP request
// with the correct parameters, and returns a list of Checkins.
func getCheckins(c *Client, endpoint string, q url.Values) ([]*Checkin, *http.Response, error) {
	// Temporary struct to unmarshal checkin JSON
	var v struct {
		Response struct {
			Checkins struct {
				Count int           `json:"count"`
				Items []*rawCheckin `json:"items"`
			} `json:"checkins"`
		} `json:"response"`
	}

	// Perform request for user checkins by ID
	res, err := c.request("GET", endpoint, q, &v)
	if err != nil {
		return nil, res, err
	}

	// Build result slice from struct
	checkins := make([]*Checkin, v.Response.Checkins.Count)
	for i := range v.Response.Checkins.Items {
		checkins[i] = v.Response.Checkins.Items[i].export()
	}

	return checkins, res, nil
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

	return c
}
