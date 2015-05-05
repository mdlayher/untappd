package untappd

import (
	"net/http"
	"net/url"
	"strconv"
)

// Search searches for information about breweries, using the specified search query.
//
// This method returns up to 25 search results.  For more granular control,
// and to page through the results list, use SearchOffsetLimit instead.
func (b *BreweryService) Search(query string) ([]*Brewery, *http.Response, error) {
	// Use default parameters as specified by API
	return b.SearchOffsetLimit(query, 0, 25)
}

// SearchOffsetLimit searches for information about breweries, using the specified
// search query.  In addition, it accepts offset and limit parameters to enable
// paging through more than 25 breweries.
//
// 50 breweries is the maximum number of results which may be returned by one call.
func (b *BreweryService) SearchOffsetLimit(query string, offset int, limit int) ([]*Brewery, *http.Response, error) {
	q := url.Values{
		"q":      []string{query},
		"offset": []string{strconv.Itoa(offset)},
		"limit":  []string{strconv.Itoa(limit)},
	}

	// Temporary struct to unmarshal breweries JSON
	var v struct {
		Response struct {
			Brewery struct {
				Count int `json:"count"`
				Items []struct {
					Brewery rawBrewery `json:"brewery"`
				} `json:"items"`
			} `json:"brewery"`
		} `json:"response"`
	}

	// Perform request for brewery search
	res, err := b.client.request("GET", "search/brewery", nil, q, &v)
	if err != nil {
		return nil, res, err
	}

	// Build result slice from struct
	breweries := make([]*Brewery, v.Response.Brewery.Count)
	for i := range v.Response.Brewery.Items {
		breweries[i] = v.Response.Brewery.Items[i].Brewery.export()
	}

	return breweries, res, nil
}
