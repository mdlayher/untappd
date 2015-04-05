package untappd

import (
	"net/http"
	"strconv"
	"strings"
	"testing"
)

// TestClientUserWishListOK verifies that Client.User.WishList always sets the
// appropriate default offset and limit values.
func TestClientUserWishListOK(t *testing.T) {
	offset := "0"
	limit := "25"
	sort := "date"

	c, done := userWishListTestClient(t, func(t *testing.T, w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()

		if o := q.Get("offset"); o != offset {
			t.Fatalf("unexpected offset parameter: %s != %s", o, offset)
		}
		if l := q.Get("limit"); l != limit {
			t.Fatalf("unexpected limit parameter: %s != %s", l, limit)
		}
		if s := q.Get("sort"); s != sort {
			t.Fatalf("unexpected sort parameter: %q != %q", s, sort)
		}

		// Empty JSON response since we already passed checks
		w.Write([]byte("{}"))
	})
	defer done()

	if _, _, err := c.User.WishList("foo"); err != nil {
		t.Fatal(err)
	}
}

// TestClientUserWishListOffsetLimitSortBadUser verifies that
// Client.User.WishListOffsetLimitSort returns an error when an invalid user
// is queried.
func TestClientUserWishListOffsetLimitBadUser(t *testing.T) {
	c, done := userWishListTestClient(t, func(t *testing.T, w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write(invalidUserErrJSON)
	})
	defer done()

	_, _, err := c.User.WishListOffsetLimitSort("foo", 0, 25, SortDate)
	assertInvalidUserErr(t, err)
}

// TestClientUserWishListOffsetLimitSortOK verifies that Client.User.WishListOffsetLimitSort
// returns a valid beers list, when used with correct parameters.
func TestClientUserWishListOffsetLimitSortOK(t *testing.T) {
	var offset int
	sOffset := strconv.Itoa(offset)

	var limit = 25
	sLimit := strconv.Itoa(limit)

	var sort = SortDate

	username := "mdlayher"
	c, done := userWishListTestClient(t, func(t *testing.T, w http.ResponseWriter, r *http.Request) {
		path := "/v4/user/wishlist/" + username + "/"
		if p := r.URL.Path; p != path {
			t.Fatalf("unexpected URL path: %q != %q", p, path)
		}

		q := r.URL.Query()

		if o := q.Get("offset"); o != sOffset {
			t.Fatalf("unexpected offset parameter: %s != %s", o, sOffset)
		}
		if l := q.Get("limit"); l != sLimit {
			t.Fatalf("unexpected limit parameter: %s != %s", l, sLimit)
		}
		if s := q.Get("sort"); s != string(sort) {
			t.Fatalf("unexpected sort parameter: %q != %q", s, sort)
		}

		w.Write(userWishListJSON)
	})
	defer done()

	beers, _, err := c.User.WishListOffsetLimitSort(username, offset, limit, sort)
	if err != nil {
		t.Fatal(err)
	}

	expected := []*Beer{
		&Beer{
			ID:    1,
			Name:  "Rare Bourbon County Brand Stout",
			Style: "American Imperial / Double Stout",
			Brewery: &Brewery{
				Name: "Goose Island Beer Co.",
			},
		},
		&Beer{
			ID:    2,
			Name:  "Double Barrel Hunahpu's",
			Style: "American Imperial / Double Stout",
			Brewery: &Brewery{
				Name: "Cigar City Brewing",
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

// userWishListTestClient builds upon testClient, and adds additional sanity checks
// for tests which target the user wishlist API.
func userWishListTestClient(t *testing.T, fn func(t *testing.T, w http.ResponseWriter, r *http.Request)) (*Client, func()) {
	return testClient(t, func(t *testing.T, w http.ResponseWriter, r *http.Request) {
		// Always GET request
		method := "GET"
		if m := r.Method; m != method {
			t.Fatalf("unexpected HTTP method: %q != %q", m, method)
		}

		// Always uses specific path prefix
		prefix := "/v4/user/wishlist/"
		if p := r.URL.Path; !strings.HasPrefix(p, prefix) {
			t.Fatalf("unexpected HTTP path prefix: %q != %q", p, prefix)
		}

		// Guard against panics
		if fn != nil {
			fn(t, w, r)
		}
	})
}

// Canned user wishlist JSON response, taken from documentation: https://untappd.com/api/docs#userwishlist
// Slight modifications made to add multiple beers to items list
var userWishListJSON = []byte(`{
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
        "beer_name": "Rare Bourbon County Brand Stout",
        "beer_style": "American Imperial / Double Stout"
      },
      "brewery": {
        "brewery_name": "Goose Island Beer Co."
      }
    },
    {
      "beer": {
        "bid": 2,
        "beer_name": "Double Barrel Hunahpu's",
        "beer_style": "American Imperial / Double Stout"
      },
      "brewery": {
        "brewery_name": "Cigar City Brewing"
      }
    }
    ]
  }
}}`)
