package untappd

import (
	"math"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

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
func (u *UserService) FriendsOffsetLimit(username string, offset int, limit int) ([]*User, *http.Response, error) {
	q := url.Values{
		"offset": []string{strconv.Itoa(offset)},
		"limit":  []string{strconv.Itoa(limit)},
	}

	// Temporary struct to unmarshal friends JSON
	var v struct {
		Response struct {
			Count int `json:"count"`
			Items []struct {
				User rawUser `json:"user"`
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

	return users, res, nil
}

// Badges queries for information about a User's badges.  The username
// parameter specifies the User whose badges will be returned.
//
// This method returns up to 50 of the User's most recently earned badges.
// For more granular control, and to page through the badges list, use
// BadgesOffsetLimit instead.
func (u *UserService) Badges(username string) ([]*Badge, *http.Response, error) {
	// Use default parameters as specified by API
	return u.BadgesOffsetLimit(username, 0, 50)
}

// BadgesOffsetLimit queries for information about a User's badges, but also
// accepts offset and limit parameters to enable paging through more than 50
// badges.  The username parameter specifies the User whose badges will be
// returned.
//
// 50 badges is the maximum number of badges which may be returned by one call.
func (u *UserService) BadgesOffsetLimit(username string, offset int, limit int) ([]*Badge, *http.Response, error) {
	q := url.Values{
		"offset": []string{strconv.Itoa(offset)},
		"limit":  []string{strconv.Itoa(limit)},
	}

	// Temporary struct to unmarshal badges JSON
	var v struct {
		Response struct {
			Count int         `json:"count"`
			Items []*rawBadge `json:"items"`
		} `json:"response"`
	}

	// Perform request for user badges by username
	res, err := u.client.request("GET", "user/badges/"+username, q, &v)
	if err != nil {
		return nil, res, err
	}

	// Build result slice from struct
	badges := make([]*Badge, v.Response.Count)
	for i := range v.Response.Items {
		badges[i] = v.Response.Items[i].export()
	}

	return badges, res, nil
}

// WishList queries for information about a User's wish list beers.
// The username parameter specifies the User whose beers will be returned.
//
// This method returns up to 25 of the User's wish list beers.
// For more granular control, and to page through and sort the beers list, use
// WishListOffsetLimitSort instead.
func (u *UserService) WishList(username string) ([]*Beer, *http.Response, error) {
	// Use default parameters as specified by API
	return u.WishListOffsetLimitSort(username, 0, 25, SortDate)
}

// WishListOffsetLimitSort queries for information about a User's wish list beers,
// but also accepts offset, limit, and sort parameters to enable paging and sorting
// through more than 25 beers.  The username parameter specifies the User whose
// wish list beers will be returned.  Beers may be sorted using any of the provided
// Sort constants with this package.
//
// 50 beers is the maximum number of beers which may be returned by one call.
func (u *UserService) WishListOffsetLimitSort(username string, offset int, limit int, sort Sort) ([]*Beer, *http.Response, error) {
	q := url.Values{
		"offset": []string{strconv.Itoa(offset)},
		"limit":  []string{strconv.Itoa(limit)},
		"sort":   []string{string(sort)},
	}

	// Temporary struct to unmarshal beers JSON
	var v struct {
		Response struct {
			Beers struct {
				Count int `json:"count"`
				Items []struct {
					WishListed responseTime `json:"created_at"`
					Beer       rawBeer      `json:"beer"`
					Brewery    rawBrewery   `json:"brewery"`
				} `json:"items"`
			} `json:"beers"`
		} `json:"response"`
	}

	// Perform request for user beers by username
	res, err := u.client.request("GET", "user/wishlist/"+username, q, &v)
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
		beers[i].WishListed = time.Time(v.Response.Beers.Items[i].WishListed)
	}

	return beers, res, nil
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
func (u *UserService) BeersOffsetLimitSort(username string, offset int, limit int, sort Sort) ([]*Beer, *http.Response, error) {
	q := url.Values{
		"offset": []string{strconv.Itoa(offset)},
		"limit":  []string{strconv.Itoa(limit)},
		"sort":   []string{string(sort)},
	}

	// Temporary struct to unmarshal beers JSON
	var v struct {
		Response struct {
			Beers struct {
				Count int `json:"count"`
				Items []struct {
					FirstHad   responseTime `json:"first_had"`
					UserRating float64      `json:"rating_score"`
					Count      int          `json:"count"`

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
		beers[i].FirstHad = time.Time(v.Response.Beers.Items[i].FirstHad)
		beers[i].UserRating = v.Response.Beers.Items[i].UserRating
		beers[i].Count = v.Response.Beers.Items[i].Count
	}

	return beers, res, nil
}

// Checkins queries for information about a User's checkins.
// The username parameter specifies the User whose checkins will be
// returned.
//
// This method returns up to 25 of the User's most recent checkins.
// For more granular control, and to page through the checkins list using ID
// parameters, use CheckinsMinMaxIDLimit instead.
func (u *UserService) Checkins(username string) ([]*Checkin, *http.Response, error) {
	// Use default parameters as specified by API.  Max ID is somewhat
	// arbitrary, but should provide plenty of headroom, just in case.
	return u.CheckinsMinMaxIDLimit(username, 0, math.MaxInt32, 25)
}

// CheckinsMinMaxIDLimit queries for information about a User's checkins,
// but also accepts minimum checkin ID, maximum checkin ID, and a limit
// parameter to enable paging through checkins. The username parameter
// specifies the User whose checkins will be returned.
//
// 50 checkins is the maximum number of checkins which may be returned by
// one call.
func (u *UserService) CheckinsMinMaxIDLimit(username string, minID int, maxID int, limit int) ([]*Checkin, *http.Response, error) {
	q := url.Values{
		"min_id": []string{strconv.Itoa(minID)},
		"max_id": []string{strconv.Itoa(maxID)},
		"limit":  []string{strconv.Itoa(limit)},
	}

	// Temporary struct to unmarshal checkin JSON
	var v struct {
		Response struct {
			Checkins struct {
				Count int `json:"count"`
				Items []struct {
					ID         int          `json:"checkin_id"`
					Beer       rawBeer      `json:"beer"`
					Brewery    rawBrewery   `json:"brewery"`
					User       rawUser      `json:"user"`
					UserRating float64      `json:"rating_score"`
					Comment    string       `json:"checkin_comment"`
					Created    responseTime `json:"created_at"`
				} `json:"items"`
			} `json:"checkins"`
		} `json:"response"`
	}

	// Perform request for user checkins by username
	res, err := u.client.request("GET", "user/checkins/"+username, q, &v)
	if err != nil {
		return nil, res, err
	}

	// Build result slice from struct
	checkins := make([]*Checkin, v.Response.Checkins.Count)
	for i := range v.Response.Checkins.Items {
		// Information about the beer itself
		checkin := &Checkin{
			ID:         v.Response.Checkins.Items[i].ID,
			Comment:    v.Response.Checkins.Items[i].Comment,
			UserRating: v.Response.Checkins.Items[i].UserRating,
			Created:    time.Time(v.Response.Checkins.Items[i].Created),
		}
		checkins[i] = checkin
		checkins[i].Beer = v.Response.Checkins.Items[i].Beer.export()

		// Information about the beer's brewery
		checkins[i].Brewery = v.Response.Checkins.Items[i].Brewery.export()
		checkins[i].User = v.Response.Checkins.Items[i].User.export()
	}

	return checkins, res, nil
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
