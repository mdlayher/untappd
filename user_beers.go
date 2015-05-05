package untappd

import (
	"net/http"
	"net/url"
	"strconv"
	"time"
)

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
	res, err := u.client.request("GET", "user/beers/"+username, nil, q, &v)
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
