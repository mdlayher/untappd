package untappd

import (
	"net/http"
	"strconv"
	"strings"
	"testing"
)

// TestClientBeerInfoBadBeer verifies that Client.Beer.Info returns an error when
// an invalid beer is queried.
func TestClientBeerInfoBadBeer(t *testing.T) {
	beerID := -1
	sBeerID := strconv.Itoa(beerID)

	c, done := beerInfoTestClient(t, func(t *testing.T, w http.ResponseWriter, r *http.Request) {
		path := "/v4/beer/info/" + sBeerID + "/"
		if p := r.URL.Path; p != path {
			t.Fatalf("unexpected URL path: %q != %q", p, path)
		}

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(invalidBeerErrJSON)
	})
	defer done()

	_, _, err := c.Beer.Info(beerID, false)
	assertInvalidBeerErr(t, err)
}

// TestClientBeerInfoCompactOK verifies that Client.Beer.Info properly requests compact
// beer output.
func TestClientBeerInfoCompactOK(t *testing.T) {
	c, done := beerInfoTestClient(t, func(t *testing.T, w http.ResponseWriter, r *http.Request) {
		if c := r.URL.Query().Get("compact"); c != "true" {
			t.Fatalf("unexpected compact query value: %q != %q", c, "true")
		}

		// In the future, we may return compact canned beer data here.
		// For now, write a mostly empty JSON object is enough to get
		// test coverage.
		w.Write([]byte(`{"response":{"beer":{"id":1}}}`))
	})
	defer done()

	if _, _, err := c.Beer.Info(1, true); err != nil {
		t.Fatal(err)
	}
}

// TestClientBeerInfoOK verifies that Client.Beer.Info returns a valid beer when
// provided with correct input parameters.
func TestClientBeerInfoOK(t *testing.T) {
	beerID := 1
	sBeerID := strconv.Itoa(beerID)

	c, done := beerInfoTestClient(t, func(t *testing.T, w http.ResponseWriter, r *http.Request) {
		path := "/v4/beer/info/" + sBeerID + "/"
		if p := r.URL.Path; p != path {
			t.Fatalf("unexpected URL path: %q != %q", p, path)
		}

		w.Write(blackNoteBeerJSON)
	})
	defer done()

	b, _, err := c.Beer.Info(beerID, false)
	if err != nil {
		t.Fatal(err)
	}

	if id := b.ID; id != beerID {
		t.Fatalf("unexpected ID: %d != %d", id, 1)
	}
	beerName := "Black Note Stout"
	if n := b.Name; n != beerName {
		t.Fatalf("unexpected Name: %q != %q", n, beerName)
	}
	breweryName := "Bell's Brewery, Inc."
	if n := b.Brewery.Name; n != breweryName {
		t.Fatalf("unexpected Brewery.Name: %q != %q", n, breweryName)
	}
}

// beerInfoTestClient builds upon testClient, and adds additional sanity checks
// for tests which target the beer info API.
func beerInfoTestClient(t *testing.T, fn func(t *testing.T, w http.ResponseWriter, r *http.Request)) (*Client, func()) {
	return testClient(t, func(t *testing.T, w http.ResponseWriter, r *http.Request) {
		// Always GET request
		method := "GET"
		if m := r.Method; m != method {
			t.Fatalf("unexpected HTTP method: %q != %q", m, method)
		}

		// Always uses specific path prefix
		prefix := "/v4/beer/info/"
		if p := r.URL.Path; !strings.HasPrefix(p, prefix) {
			t.Fatalf("unexpected HTTP path prefix: %q != %q", p, prefix)
		}

		// Guard against panics
		if fn != nil {
			fn(t, w, r)
		}
	})
}

// Canned JSON used in tests
var blackNoteBeerJSON = []byte(`
{
  "meta": {
    "code": 200,
    "response_time": {
      "time": 0,
      "measure": "seconds"
    }
  },
  "notifications": {},
  "response": {
  "beer": {
    "bid": 1,
    "beer_name": "Black Note Stout",
    "brewery": {
      "brewery_name": "Bell's Brewery, Inc."
    }
  }
  }
}`)
