package untappd_test

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/mdlayher/untappd"
)

// TestClientBrewerySearchOK verifies that Client.Brewery.Search always sets the
// appropriate default offset and limit values.
func TestClientBrewerySearchOK(t *testing.T) {
	query := "foo"
	offset := "0"
	limit := "25"

	c, done := brewerySearchTestClient(t, func(t *testing.T, w http.ResponseWriter, r *http.Request) {
		assertParameters(t, r, url.Values{
			"q":      []string{query},
			"offset": []string{offset},
			"limit":  []string{limit},
		})

		// Empty JSON response since we already passed checks
		w.Write([]byte("{}"))
	})
	defer done()

	if _, _, err := c.Brewery.Search(query); err != nil {
		t.Fatal(err)
	}
}

// TestClientBrewerySearchOffsetLimitBadQuery verifies that
// Client.Brewery.SearchOffsetLimit returns an error when an invalid query
// parameter is used.
func TestClientBrewerySearchOffsetLimitBadQuery(t *testing.T) {
	c, done := brewerySearchTestClient(t, func(t *testing.T, w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(invalidQueryErrJSON)
	})
	defer done()

	_, _, err := c.Brewery.SearchOffsetLimit("", 0, 25)
	assertInvalidQueryErr(t, err)
}

// TestClientBrewerySearchOffsetLimitOK verifies that Client.Brewery.SearchOffsetLimit
// returns a valid results list, when used with correct parameters.
func TestClientBrewerySearchOffsetLimitOK(t *testing.T) {
	var offset int
	sOffset := strconv.Itoa(offset)

	var limit = 25
	sLimit := strconv.Itoa(limit)

	query := "russian river"
	c, done := brewerySearchTestClient(t, func(t *testing.T, w http.ResponseWriter, r *http.Request) {
		assertParameters(t, r, url.Values{
			"q":      []string{query},
			"offset": []string{sOffset},
			"limit":  []string{sLimit},
		})

		w.Write(brewerySearchJSON)
	})
	defer done()

	breweries, _, err := c.Brewery.SearchOffsetLimit(query, offset, limit)
	if err != nil {
		t.Fatal(err)
	}

	expected := []*untappd.Brewery{
		&untappd.Brewery{
			ID:      1,
			Name:    "Russian River Brewing Company",
			Country: "United States",
		},
	}

	for i := range breweries {
		if breweries[i].ID != expected[i].ID {
			t.Fatalf("unexpected brewery ID: %d != %d", breweries[i].ID, expected[i].ID)
		}
		if breweries[i].Name != expected[i].Name {
			t.Fatalf("unexpected brewery Name: %q != %q", breweries[i].Name, expected[i].Name)
		}
		if breweries[i].Country != expected[i].Country {
			t.Fatalf("unexpected brewery Country: %q != %q", breweries[i].Country, expected[i].Country)
		}
	}
}

// brewerySearchTestClient builds upon testClient, and adds additional sanity checks
// for tests which target the user breweries API.
func brewerySearchTestClient(t *testing.T, fn func(t *testing.T, w http.ResponseWriter, r *http.Request)) (*untappd.Client, func()) {
	return testClient(t, func(t *testing.T, w http.ResponseWriter, r *http.Request) {
		// Always GET request
		method := "GET"
		if m := r.Method; m != method {
			t.Fatalf("unexpected HTTP method: %q != %q", m, method)
		}

		// Always uses specific path prefix
		prefix := "/v4/search/brewery/"
		if p := r.URL.Path; !strings.HasPrefix(p, prefix) {
			t.Fatalf("unexpected HTTP path prefix: %q != %q", p, prefix)
		}

		// Guard against panics
		if fn != nil {
			fn(t, w, r)
		}
	})
}

// Canned brewery search JSON response, taken from documentation: https://untappd.com/api/docs#brewerysearch
var brewerySearchJSON = []byte(`{
  "meta": {
    "code": 200,
    "response_time": {
      "time": 0,
      "measure": "seconds"
    }
  },
  "notifications": {},
  "response": {
  "brewery": {
    "count": 1,
    "items": [
    {
      "brewery": {
        "brewery_id": 1,
        "brewery_name": "Russian River Brewing Company",
        "country_name": "United States"
      }
    }
    ]
  }
}}`)
