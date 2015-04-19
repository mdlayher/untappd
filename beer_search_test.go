package untappd

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"
)

// TestClientBeerSearchOK verifies that Client.Beer.Search always sets the
// appropriate default offset, limit, and sort values.
func TestClientBeerSearchOK(t *testing.T) {
	query := "foo"
	offset := "0"
	limit := "25"
	sort := "date"

	c, done := beerSearchTestClient(t, func(t *testing.T, w http.ResponseWriter, r *http.Request) {
		assertParameters(t, r, url.Values{
			"q":      []string{query},
			"offset": []string{offset},
			"limit":  []string{limit},
			"sort":   []string{sort},
		})

		// Empty JSON response since we already passed checks
		w.Write([]byte("{}"))
	})
	defer done()

	if _, _, err := c.Beer.Search(query); err != nil {
		t.Fatal(err)
	}
}

// TestClientBeerSearchOffsetLimitSortBadQuery verifies that
// Client.Beer.SearchOffsetLimitSort returns an error when an invalid query
// parameter is used.
func TestClientBeerSearchOffsetLimitBadQuery(t *testing.T) {
	c, done := beerSearchTestClient(t, func(t *testing.T, w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(invalidQueryErrJSON)
	})
	defer done()

	_, _, err := c.Beer.SearchOffsetLimitSort("", 0, 25, SortDate)
	assertInvalidQueryErr(t, err)
}

// TestClientBeerSearchOffsetLimitSortOK verifies that Client.Beer.SearchOffsetLimitSort
// returns a valid results list, when used with correct parameters.
func TestClientBeerSearchOffsetLimitOK(t *testing.T) {
	var offset int
	sOffset := strconv.Itoa(offset)

	var limit = 25
	sLimit := strconv.Itoa(limit)

	var sort = SortDate

	query := "russian river pliny"
	c, done := beerSearchTestClient(t, func(t *testing.T, w http.ResponseWriter, r *http.Request) {
		assertParameters(t, r, url.Values{
			"q":      []string{query},
			"offset": []string{sOffset},
			"limit":  []string{sLimit},
			"sort":   []string{string(sort)},
		})

		w.Write(beerSearchJSON)
	})
	defer done()

	beers, _, err := c.Beer.SearchOffsetLimitSort(query, offset, limit, sort)
	if err != nil {
		t.Fatal(err)
	}

	expected := []*Beer{
		&Beer{
			ID:    1,
			Name:  "Pliny the Elder",
			Style: "Imperial / Double IPA",
			Brewery: &Brewery{
				Name: "Russian River Brewing Company",
			},
		},
		&Beer{
			ID:    2,
			Name:  "Pliny the Younger",
			Style: "Triple IPA",
			Brewery: &Brewery{
				Name: "Russian River Brewing Company",
			},
		},
	}

	for i := range beers {
		if beers[i].ID != expected[i].ID {
			t.Fatalf("unexpected beer ID: %d != %d", beers[i].ID, expected[i].ID)
		}
		if beers[i].Name != expected[i].Name {
			t.Fatalf("unexpected beer Name: %q != %q", beers[i].Name, expected[i].Name)
		}
		if beers[i].Style != expected[i].Style {
			t.Fatalf("unexpected beer Style: %q != %q", beers[i].Style, expected[i].Style)
		}
		if beers[i].Brewery.Name != expected[i].Brewery.Name {
			t.Fatalf("unexpected beer Brewery.Name: %q != %q", beers[i].Brewery.Name, expected[i].Brewery.Name)
		}
	}
}

// beerSearchTestClient builds upon testClient, and adds additional sanity checks
// for tests which target the user beers API.
func beerSearchTestClient(t *testing.T, fn func(t *testing.T, w http.ResponseWriter, r *http.Request)) (*Client, func()) {
	return testClient(t, func(t *testing.T, w http.ResponseWriter, r *http.Request) {
		// Always GET request
		method := "GET"
		if m := r.Method; m != method {
			t.Fatalf("unexpected HTTP method: %q != %q", m, method)
		}

		// Always uses specific path prefix
		prefix := "/v4/search/beer/"
		if p := r.URL.Path; !strings.HasPrefix(p, prefix) {
			t.Fatalf("unexpected HTTP path prefix: %q != %q", p, prefix)
		}

		// Guard against panics
		if fn != nil {
			fn(t, w, r)
		}
	})
}

// Canned beer search JSON response, taken from documentation: https://untappd.com/api/docs#beersearch
// Slight modifications made to add multiple beers to items list
var beerSearchJSON = []byte(`{
  "meta": {
    "code": 200,
    "response_time": {
      "time": 0,
      "measure": "seconds"
    }
  },
  "notifications": {},
  "response": {
  "beers": {
    "count": 2,
    "items": [
    {
      "beer": {
        "bid": 1,
        "beer_name": "Pliny the Elder",
        "beer_style": "Imperial / Double IPA"
      },
      "brewery": {
        "brewery_name": "Russian River Brewing Company"
      }
    },
    {
      "beer": {
        "bid": 2,
        "beer_name": "Pliny the Younger",
        "beer_style": "Triple IPA"
      },
      "brewery": {
        "brewery_name": "Russian River Brewing Company"
      }
    }
    ]
  }
}}`)
