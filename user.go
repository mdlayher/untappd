package untappd

import (
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// User represents an Untappd user.
type User struct {
	UID        int64
	ID         int64
	UserName   string
	FirstName  string
	LastName   string
	Avatar     url.URL
	CoverPhoto url.URL
	Location   string
	URL        url.URL
	Bio        string
	Supporter  bool
	UntappdURL url.URL
	Stats      UserStats
}

// UserStats is a struct which contains various statistics regarding an Untappd
// user.
type UserStats struct {
	TotalBadges       int64 `json:"total_badges"`
	TotalFriends      int64 `json:"total_friends"`
	TotalCheckins     int64 `json:"total_checkins"`
	TotalBeers        int64 `json:"total_beers"`
	TotalCreatedBeers int64 `json:"total_created_beers"`
	TotalFollowings   int64 `json:"total_followings"`
	TotalPhotos       int64 `json:"total_photos"`
}

// UserService is a "service" which allows access to API methods involving users.
type UserService struct {
	client *Client
}

// Info queries for information about a User with the specified username.
// If the compact parameter is set to 'true', only basic user information will
// be populated.
func (u *UserService) Info(username string, compact bool) (*User, *http.Response, error) {
	// Determine if a compact response is requested
	q := url.Values{}
	if compact {
		q.Set("compact", "true")
	}

	// Temporary struct to unmarshal raw user JSON
	var v struct {
		Response struct {
			User rawUser `json:"user"`
		} `json:"response"`
	}

	// Perform request for user information by username
	res, err := u.client.request("GET", "user/info/"+username, q, &v)
	if err != nil {
		return nil, res, err
	}

	// Return results
	return v.Response.User.export(), res, nil
}

// Friends queries for information about a User's friends.  The username
// parameter specifies the User whose friends will be returned.
//
// This method returns up to a maximum of 25 friends.  For more granular
// control, and to page through the friends list, use FriendsOffsetLimit
// instead.
//
// The resulting slice of User structs contains a more limited set of user
// information than a call to Info would.  However, basic information such as
// user ID, username, first name, last name, bio, etc. is available.
func (u *UserService) Friends(username string) ([]*User, *http.Response, error) {
	// Use default parameters as specified by API
	return u.FriendsOffsetLimit(username, 0, 25)
}

// FriendsOffsetLimit queries for information about a User's friends, but also
// accepts offset and limit parameters to enable paging through more than 25
// friends.  The username parameter specifies the User whose friends will be
// returned.
//
// 25 friends is the maximum number of friends which may be returned by one call.
func (u *UserService) FriendsOffsetLimit(username string, offset uint, limit uint) ([]*User, *http.Response, error) {
	q := url.Values{
		"offset": []string{strconv.Itoa(int(offset))},
		"limit":  []string{strconv.Itoa(int(limit))},
	}

	// Temporary struct to unmarshal friends JSON
	var v struct {
		Response struct {
			Count int `json:"count"`
			Items []struct {
				User *rawUser `json:"user"`
			} `json:"items"`
		} `json:"response"`
	}

	// Perform request for user friends by username
	res, err := u.client.request("GET", "user/friends/"+username, q, &v)
	if err != nil {
		return nil, res, err
	}

	// Build result slice from struct
	users := make([]*User, v.Response.Count)
	for i := range v.Response.Items {
		users[i] = v.Response.Items[i].User.export()
	}

	// Return results
	return users, res, nil
}

// Badge represents an Untappd badge.
//
// BUG(mdlayher): write out fields to access more badge information.
type Badge struct {
	ID          int64  `json:"badge_id"`
	CheckinID   int64  `json:"checkin_id"`
	Name        string `json:"badge_name"`
	Description string `json:"badge_description"`
}

// Badges queries for information about a User's badges.  The username
// parameter specifies the User whose badges will be returned.
//
// This method returns up to 50 of the User's most recently earned badges.
// For more granular control, and to page through the badges list, use
// BadgesOffset instead.
func (u *UserService) Badges(username string) ([]*Badge, *http.Response, error) {
	// Use default parameters as specified by API
	return u.BadgesOffset(username, 0)
}

// BadgesOffset queries for information about a User's badges, but also
// accepts an offset parameter to enable paging through more than 50
// badges.  The username parameter specifies the User whose badges will be
// returned.
//
// 50 badges is the maximum number of badges which may be returned by one call.
func (u *UserService) BadgesOffset(username string, offset uint) ([]*Badge, *http.Response, error) {
	q := url.Values{
		"offset": []string{strconv.Itoa(int(offset))},
	}

	// Temporary struct to unmarshal badges JSON
	var v struct {
		Response struct {
			Count int      `json:"count"`
			Items []*Badge `json:"items"`
		} `json:"response"`
	}

	// Perform request for user badges by username
	res, err := u.client.request("GET", "user/badges/"+username, q, &v)
	if err != nil {
		return nil, res, err
	}

	// Return results
	return v.Response.Items, res, nil
}

// Beers queries for information about a User's checked-in beers.
// The username parameter specifies the User whose beers will be returned.
//
// This method returns up to 25 of the User's most recently checked-in beerss.
// For more granular control, and to page through and sort the beers list, use
// BeersOffsetLimitSort instead.
func (u *UserService) Beers(username string) ([]*Beer, *http.Response, error) {
	// Use default parameters as specified by API
	return u.BeersOffsetLimitSort(username, 0, 25, SortDate)
}

// BeersOffsetLimitSort queries for information about a User's checked-in beers,
// but also accepts offset, limit, and sort parameters to enable paging and sorting
// through more than 25 beers.  The username parameter specifies the User whose
// checked-in beers will be returned.  Beers may be sorted using any of the provided
// Sort constants with this package.
//
// 50 beers is the maximum number of beers which may be returned by one call.
func (u *UserService) BeersOffsetLimitSort(username string, offset uint, limit uint, sort Sort) ([]*Beer, *http.Response, error) {
	q := url.Values{
		"offset": []string{strconv.Itoa(int(offset))},
		"limit":  []string{strconv.Itoa(int(limit))},
		"sort":   []string{string(sort)},
	}

	// Temporary struct to unmarshal beers JSON
	var v struct {
		Response struct {
			Beers struct {
				Count int `json:"count"`
				Items []struct {
					FirstCheckinID int64        `json:"first_checkin_id"`
					FirstCheckin   responseTime `json:"first_checkin_at"`
					UserRating     float64      `json:"rating_score"`

					Beer    rawBeer    `json:"beer"`
					Brewery rawBrewery `json:"brewery"`
				} `json:"items"`
			} `json:"beers"`
		} `json:"response"`
	}

	// Perform request for user beers by username
	res, err := u.client.request("GET", "user/beers/"+username, q, &v)
	if err != nil {
		return nil, res, err
	}

	// Build result slice from struct
	beers := make([]*Beer, v.Response.Beers.Count)
	for i := range v.Response.Beers.Items {
		// Information about the beer itself
		beers[i] = v.Response.Beers.Items[i].Beer.export()

		// Information about the beer's brewery
		beers[i].Brewery = v.Response.Beers.Items[i].Brewery.export()

		// Information related to this user and this beer
		beers[i].FirstCheckinID = v.Response.Beers.Items[i].FirstCheckinID
		beers[i].FirstCheckin = time.Time(v.Response.Beers.Items[i].FirstCheckin)
		beers[i].UserRating = v.Response.Beers.Items[i].UserRating
	}

	// Return results
	return beers, res, nil
}

// rawUser is the raw JSON representation of an Untappd user.  Its data is
// unmarshaled from JSON and then exported to a User struct.
type rawUser struct {
	UID        int64        `json:"uid"`
	ID         int64        `json:"id"`
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
