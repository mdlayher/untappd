package untappd

import (
	"net/http"
	"net/url"
	"strconv"
	"time"
)

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
