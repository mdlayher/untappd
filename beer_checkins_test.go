package untappd

import (
	"math"
	"net/http"
	"strconv"
	"strings"
	"testing"
)

// TestClientBeerCheckinsOK verifies that Client.Beer.Checkins always sets the
// appropriate default minimum ID, maximum ID, and limit values.
func TestClientBeerCheckinsOK(t *testing.T) {
	minID := "0"
	maxID := strconv.Itoa(math.MaxInt32)
	limit := "25"

	c, done := beerCheckinsTestClient(t, func(t *testing.T, w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()

		if i := q.Get("min_id"); i != minID {
			t.Fatalf("unexpected min_id parameter: %s != %s", i, minID)
		}
		if i := q.Get("max_id"); i != maxID {
			t.Fatalf("unexpected max_id parameter: %s != %s", i, maxID)
		}
		if l := q.Get("limit"); l != limit {
			t.Fatalf("unexpected limit parameter: %s != %s", l, limit)
		}

		// Empty JSON response since we already passed checks
		w.Write([]byte("{}"))
	})
	defer done()

	if _, _, err := c.Beer.Checkins(1); err != nil {
		t.Fatal(err)
	}
}

// TestClientBeerCheckinsMinMaxIDLimitBadBeer verifies that
// Client.Beer.CheckinsMinMaxIDLimit returns an error when an invalid beer
// is queried.
func TestClientBeerCheckinsMinMaxIDLimitBadBeer(t *testing.T) {
	c, done := beerCheckinsTestClient(t, func(t *testing.T, w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(invalidBeerErrJSON)
	})
	defer done()

	_, _, err := c.Beer.CheckinsMinMaxIDLimit(-1, 0, math.MaxInt32, 25)
	assertInvalidBeerErr(t, err)
}

// TestClientBeerCheckinsMinMaxIDLimitOK verifies that Client.Beer.CheckinsMinMaxIDLimit
// returns a valid checkins list, when used with correct parameters.
func TestClientBeerCheckinsMinMaxIDLimitOffsetLimitOK(t *testing.T) {
	var minID int
	sMinID := strconv.Itoa(minID)

	var maxID = math.MaxInt32
	sMaxID := strconv.Itoa(maxID)

	var limit = 25
	sLimit := strconv.Itoa(limit)

	id := 1
	sID := strconv.Itoa(id)
	c, done := beerCheckinsTestClient(t, func(t *testing.T, w http.ResponseWriter, r *http.Request) {
		path := "/v4/beer/checkins/" + sID + "/"
		if p := r.URL.Path; p != path {
			t.Fatalf("unexpected URL path: %q != %q", p, path)
		}

		q := r.URL.Query()

		if i := q.Get("min_id"); i != sMinID {
			t.Fatalf("unexpected min_id parameter: %s != %s", i, sMinID)
		}
		if i := q.Get("max_id"); i != sMaxID {
			t.Fatalf("unexpected max_id parameter: %s != %s", i, sMaxID)
		}
		if l := q.Get("limit"); l != sLimit {
			t.Fatalf("unexpected limit parameter: %s != %s", l, sLimit)
		}

		// JSON is in same format as /v4/user/checkins, so we can
		// reuse it here
		w.Write(userCheckinsJSON)
	})
	defer done()

	checkins, _, err := c.Beer.CheckinsMinMaxIDLimit(id, minID, maxID, limit)
	if err != nil {
		t.Fatal(err)
	}

	// Check data against expected set of checkins
	assertExpectedCheckins(t, checkins)
}

// beerCheckinTestClient builds upon testClient, and adds additional sanity checks
// for tests which target the beer checkin API.
func beerCheckinsTestClient(t *testing.T, fn func(t *testing.T, w http.ResponseWriter, r *http.Request)) (*Client, func()) {
	return testClient(t, func(t *testing.T, w http.ResponseWriter, r *http.Request) {
		// Always GET request
		method := "GET"
		if m := r.Method; m != method {
			t.Fatalf("unexpected HTTP method: %q != %q", m, method)
		}

		// Always uses specific path prefix
		prefix := "/v4/beer/checkins/"
		if p := r.URL.Path; !strings.HasPrefix(p, prefix) {
			t.Fatalf("unexpected HTTP path prefix: %q != %q", p, prefix)
		}

		// Guard against panics
		if fn != nil {
			fn(t, w, r)
		}
	})
}
