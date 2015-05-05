package untappd

import (
	"net/http"
	"net/url"
	"strconv"
)

// Search searches for information about beers, using the specified search query.
//
// This method returns up to 25 search results.  For more granular control,
// and to page through and sort the results list, use SearchOffsetLimitSort instead.
//
// It is recommended to search using a "Brewery Name + Beer Name" query, such as
// "Dogfish 60 Minute".
func (b *BeerService) Search(query string) ([]*Beer, *http.Response, error) {
	// Use default parameters as specified by API
	return b.SearchOffsetLimitSort(query, 0, 25, SortDate)
}

// SearchOffsetLimitSort searches for information about beers, using the specified
// search query.  In addition, it accepts offset, limit, and sort parameters to
// enable paging and sorting through more than 25 beers.  Beers may be sorted using
// any of the provided Sort constants with this package.
//
// 50 beers is the maximum number of results which may be returned by one call.
//
// It is recommended to search using a "Brewery Name + Beer Name" query, such as
// "Dogfish 60 Minute".
func (b *BeerService) SearchOffsetLimitSort(query string, offset int, limit int, sort Sort) ([]*Beer, *http.Response, error) {
	q := url.Values{
		"q":      []string{query},
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
					Beer    rawBeer    `json:"beer"`
					Brewery rawBrewery `json:"brewery"`
				} `json:"items"`
			} `json:"beers"`
		} `json:"response"`
	}

	// Perform request for beer search
	res, err := b.client.request("GET", "search/beer", nil, q, &v)
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
	}

	return beers, res, nil
}
