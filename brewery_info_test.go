package untappd_test

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/mdlayher/untappd"
)

// TestClientBreweryInfoBadBrewery verifies that Client.Brewery.Info returns an error when
// an invalid brewery is queried.
func TestClientBreweryInfoBadBrewery(t *testing.T) {
	breweryID := -1
	sBreweryID := strconv.Itoa(breweryID)

	c, done := breweryInfoTestClient(t, func(t *testing.T, w http.ResponseWriter, r *http.Request) {
		path := "/v4/brewery/info/" + sBreweryID + "/"
		if p := r.URL.Path; p != path {
			t.Fatalf("unexpected URL path: %q != %q", p, path)
		}

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(invalidBreweryErrJSON)
	})
	defer done()

	_, _, err := c.Brewery.Info(breweryID, false)
	assertInvalidBreweryErr(t, err)
}

// TestClientBreweryInfoCompactOK verifies that Client.Brewery.Info properly requests compact
// brewery output.
func TestClientBreweryInfoCompactOK(t *testing.T) {
	c, done := breweryInfoTestClient(t, func(t *testing.T, w http.ResponseWriter, r *http.Request) {
		assertParameters(t, r, url.Values{
			"compact": []string{"true"},
		})

		// In the future, we may return compact canned brewery data here.
		// For now, write a mostly empty JSON object is enough to get
		// test coverage.
		w.Write([]byte(`{"response":{"brewery":{"brewery_id":1}}}`))
	})
	defer done()

	if _, _, err := c.Brewery.Info(1, true); err != nil {
		t.Fatal(err)
	}
}

// TestClientBreweryInfoOK verifies that Client.Brewery.Info returns a valid brewery when
// provided with correct input parameters.
func TestClientBreweryInfoOK(t *testing.T) {
	breweryID := 1
	sBreweryID := strconv.Itoa(breweryID)

	c, done := breweryInfoTestClient(t, func(t *testing.T, w http.ResponseWriter, r *http.Request) {
		path := "/v4/brewery/info/" + sBreweryID + "/"
		if p := r.URL.Path; p != path {
			t.Fatalf("unexpected URL path: %q != %q", p, path)
		}

		w.Write(bellsBreweryJSON)
	})
	defer done()

	b, _, err := c.Brewery.Info(breweryID, false)
	if err != nil {
		t.Fatal(err)
	}

	if id := b.ID; id != breweryID {
		t.Fatalf("unexpected ID: %d != %d", id, 1)
	}
	breweryName := "Bell's Brewery, Inc."
	if n := b.Name; n != breweryName {
		t.Fatalf("unexpected Brewery.Name: %q != %q", n, breweryName)
	}
	brewerySlug := "bells-brewery-inc"
	if n := b.Slug; n != brewerySlug {
		t.Fatalf("unexpected Brewery.Slug: %q != %q", n, brewerySlug)
	}
	breweryType := "Micro Brewery"
	if n := b.Type; n != breweryType {
		t.Fatalf("unexpected Brewery.Type: %q != %q", n, breweryType)
	}
	breweryTypeID := 2
	if n := b.TypeID; n != breweryTypeID {
		t.Fatalf("unexpected Brewery.TypeID: %q != %q", n, breweryTypeID)
	}
	breweryContactTwitter := "BellsBrewery"
	if n := b.Contact.Twitter; n != breweryContactTwitter {
		t.Fatalf("unexpected Brewery.TypeID: %q != %q", n, breweryContactTwitter)
	}
}

// breweryInfoTestClient builds upon testClient, and adds additional sanity checks
// for tests which target the brewery info API.
func breweryInfoTestClient(t *testing.T, fn func(t *testing.T, w http.ResponseWriter, r *http.Request)) (*untappd.Client, func()) {
	return testClient(t, func(t *testing.T, w http.ResponseWriter, r *http.Request) {
		// Always GET request
		method := "GET"
		if m := r.Method; m != method {
			t.Fatalf("unexpected HTTP method: %q != %q", m, method)
		}

		// Always uses specific path prefix
		prefix := "/v4/brewery/info/"
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
var bellsBreweryJSON = []byte(`
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
    "brewery": {
      "brewery_id": 1,
      "brewery_name": "Bell's Brewery, Inc.",
      "brewery_slug": "bells-brewery-inc",
      "brewery_type": "Micro Brewery",
      "brewery_type_id": 2,
      "contact": {
        "twitter": "BellsBrewery",
        "facebook": "",
        "url": ""
      }
    }
  }
}`)
