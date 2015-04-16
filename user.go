package untappd

import (
	"net/url"
)

// UserService is a "service" which allows access to API methods involving users.
type UserService struct {
	client *Client
}

// User represents an Untappd user, and contains information regarding a user's
// username, first and last name, avatar, cover photo, and various other attributes.
type User struct {
	// Metadata from Untappd.
	UID       int
	ID        int
	UserName  string
	FirstName string
	LastName  string
	Location  string
	Bio       string
	Supporter bool

	// Links to the user's avatar, cover photo, custom URL, and Untappd profile.
	Avatar     url.URL
	CoverPhoto url.URL
	URL        url.URL
	UntappdURL url.URL

	// Struct containing this user's total badges, friends, checkins,
	// and other various totals.
	Stats UserStats
}

// UserStats is a struct which contains various statistics regarding an Untappd
// user.
type UserStats struct {
	TotalBadges       int `json:"total_badges"`
	TotalFriends      int `json:"total_friends"`
	TotalCheckins     int `json:"total_checkins"`
	TotalBeers        int `json:"total_beers"`
	TotalCreatedBeers int `json:"total_created_beers"`
	TotalFollowings   int `json:"total_followings"`
	TotalPhotos       int `json:"total_photos"`
}

// rawUser is the raw JSON representation of an Untappd user.  Its data is
// unmarshaled from JSON and then exported to a User struct.
type rawUser struct {
	UID        int          `json:"uid"`
	ID         int          `json:"id"`
	UserName   string       `json:"user_name"`
	FirstName  string       `json:"first_name"`
	LastName   string       `json:"last_name"`
	Avatar     responseURL  `json:"user_avatar"`
	AvatarHD   responseURL  `json:"user_avatar_hd"`
	CoverPhoto responseURL  `json:"user_cover_photo"`
	Location   string       `json:"location"`
	URL        responseURL  `json:"url"`
	Bio        string       `json:"bio"`
	Supporter  responseBool `json:"is_supporter"`
	UntappdURL responseURL  `json:"untappd_url"`
	Stats      UserStats    `json:"stats"`
}

// export creates an exported User from a rawUser struct, allowing for more
// useful structures to be created for client consumption.
func (r *rawUser) export() *User {
	u := &User{
		UID:        r.UID,
		ID:         r.ID,
		UserName:   r.UserName,
		FirstName:  r.FirstName,
		LastName:   r.LastName,
		Avatar:     url.URL(r.Avatar),
		CoverPhoto: url.URL(r.CoverPhoto),
		Location:   r.Location,
		URL:        url.URL(r.URL),
		Bio:        r.Bio,
		Supporter:  bool(r.Supporter),
		UntappdURL: url.URL(r.UntappdURL),
		Stats:      r.Stats,
	}

	// If high resolution avatar is available, use it instead
	if a := url.URL(r.AvatarHD); a.String() != "" {
		u.Avatar = a
	}

	return u
}
