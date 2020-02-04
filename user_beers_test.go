package untappd_test

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/mdlayher/untappd"
)

// TestClientUserBeersOK verifies that Client.User.Beers always sets the
// appropriate default offset and limit values.
func TestClientUserBeersOK(t *testing.T) {
	offset := "0"
	limit := "25"
	sort := "date"

	c, done := userBeersTestClient(t, func(t *testing.T, w http.ResponseWriter, r *http.Request) {
		assertParameters(t, r, url.Values{
			"offset": []string{offset},
			"limit":  []string{limit},
			"sort":   []string{sort},
		})

		// Empty JSON response since we already passed checks
		w.Write([]byte("{}"))
	})
	defer done()

	if _, _, err := c.User.Beers("foo"); err != nil {
		t.Fatal(err)
	}
}

// TestClientUserBeersOffsetLimitSortBadUser verifies that
// Client.User.BeersOffsetLimitSort returns an error when an invalid user
// is queried.
func TestClientUserBeersOffsetLimitBadUser(t *testing.T) {
	c, done := userBeersTestClient(t, func(t *testing.T, w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write(invalidUserErrJSON)
	})
	defer done()

	_, _, err := c.User.BeersOffsetLimitSort("foo", 0, 25, untappd.SortDate)
	assertInvalidUserErr(t, err)
}

// TestClientUserBeersOffsetLimitOK verifies that Client.User.BeersOffsetLimit
// returns a valid beers list, when used with correct parameters.
func TestClientUserBeersOffsetLimitOK(t *testing.T) {
	var offset int
	sOffset := strconv.Itoa(offset)

	var limit = 25
	sLimit := strconv.Itoa(limit)

	var sort = untappd.SortDate

	username := "mdlayher"
	c, done := userBeersTestClient(t, func(t *testing.T, w http.ResponseWriter, r *http.Request) {
		path := "/v4/user/beers/" + username + "/"
		if p := r.URL.Path; p != path {
			t.Fatalf("unexpected URL path: %q != %q", p, path)
		}

		assertParameters(t, r, url.Values{
			"offset": []string{sOffset},
			"limit":  []string{sLimit},
			"sort":   []string{string(sort)},
		})

		w.Write(userBeersJSON)
	})
	defer done()

	beers, _, err := c.User.BeersOffsetLimitSort(username, offset, limit, sort)
	if err != nil {
		t.Fatal(err)
	}

	expected := []*untappd.Beer{
		&untappd.Beer{
			ID:         1,
			Name:       "Oberon Ale",
			Style:      "American Pale Wheat Ale",
			FirstHad:   time.Date(2016, 12, 26, 1, 2, 3, 0, time.FixedZone("-0500", -5*60*60)),
			RecentHad:  time.Date(2016, 12, 31, 19, 48, 38, 0, time.FixedZone("-0500", -5*60*60)),
			UserRating: 3.75,
			Count:      1,
			Brewery: &untappd.Brewery{
				Name: "Bell's Brewery, Inc.",
			},
		},
		&untappd.Beer{
			ID:         2,
			Name:       "Two Hearted Ale",
			Style:      "American IPA",
			FirstHad:   time.Date(2016, 12, 26, 4, 5, 6, 0, time.FixedZone("-0500", -5*60*60)),
			RecentHad:  time.Date(2016, 12, 27, 19, 48, 38, 0, time.FixedZone("-0500", -5*60*60)),
			UserRating: 4.25,
			Count:      1,
			Brewery: &untappd.Brewery{
				Name: "Bell's Brewery, Inc.",
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
		if !beers[i].FirstHad.Equal(expected[i].FirstHad) {
			t.Fatalf("unexpected beer FirstHad: %q != %q", beers[i].FirstHad, expected[i].FirstHad)
		}
		if !beers[i].RecentHad.Equal(expected[i].RecentHad) {
			t.Fatalf("unexpected beer RecentHad: %q != %q", beers[i].RecentHad, expected[i].RecentHad)
		}
		if beers[i].UserRating != expected[i].UserRating {
			t.Fatalf("unexpected beer UserRating: %f != %f", beers[i].UserRating, expected[i].UserRating)
		}
		if beers[i].Count != expected[i].Count {
			t.Fatalf("unexpected beer Count: %q != %q", beers[i].Count, expected[i].Count)
		}
	}
}

// userBeersTestClient builds upon testClient, and adds additional sanity checks
// for tests which target the user beers API.
func userBeersTestClient(t *testing.T, fn func(t *testing.T, w http.ResponseWriter, r *http.Request)) (*untappd.Client, func()) {
	return testClient(t, func(t *testing.T, w http.ResponseWriter, r *http.Request) {
		// Always GET request
		method := "GET"
		if m := r.Method; m != method {
			t.Fatalf("unexpected HTTP method: %q != %q", m, method)
		}

		// Always uses specific path prefix
		prefix := "/v4/user/beers/"
		if p := r.URL.Path; !strings.HasPrefix(p, prefix) {
			t.Fatalf("unexpected HTTP path prefix: %q != %q", p, prefix)
		}

		// Guard against panics
		if fn != nil {
			fn(t, w, r)
		}
	})
}

// Canned user beers JSON response, taken from documentation: https://untappd.com/api/docs#userbeers
// Slight modifications made to add multiple beers to items list
var userBeersJSON = []byte(`{
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
      "first_checkin_id": 401400204,
      "first_created_at": "Mon, 26 Dec 2016 01:02:03 -0500",
      "recent_checkin_id": 401400204,
      "recent_created_at": "Sat, 31 Dec 2016 19:48:38 -0500",
      "recent_created_at_timezone": "-5",
      "rating_score": 3.75,
      "first_had": "Mon, 26 Dec 2016 01:02:03 -0500",
      "count": 1,
      "beer": {
        "bid": 1,
        "beer_name": "Oberon Ale",
        "beer_style": "American Pale Wheat Ale"
      },
      "brewery": {
        "brewery_name": "Bell's Brewery, Inc."
      }
    },
    {
      "first_checkin_id": 401400204,
      "first_created_at": "Mon, 26 Dec 2016 04:05:06 -0500",
      "recent_checkin_id": 401400204,
      "recent_created_at": "Tue, 27 Dec 2016 19:48:38 -0500",
      "recent_created_at_timezone": "-5",
      "rating_score": 4.25,
      "first_had": "Mon, 26 Dec 2016 04:05:06 -0500",
      "count": 1,
      "beer": {
        "bid": 2,
        "beer_name": "Two Hearted Ale",
        "beer_style": "American IPA"
      },
      "brewery": {
        "brewery_name": "Bell's Brewery, Inc."
      }
    }
    ]
  }
}}`)
