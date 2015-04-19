package untappd

import (
	"math"
	"net/http"
	"strconv"
	"strings"
	"testing"
)

// TestClientUserCheckinsOK verifies that Client.User.Checkins always sets the
// appropriate default minimum ID, maximum ID, and limit values.
func TestClientUserCheckinsOK(t *testing.T) {
	minID := "0"
	maxID := strconv.Itoa(math.MaxInt32)
	limit := "25"

	c, done := userCheckinsTestClient(t, func(t *testing.T, w http.ResponseWriter, r *http.Request) {
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

	if _, _, err := c.User.Checkins("foo"); err != nil {
		t.Fatal(err)
	}
}

// TestClientUserCheckinsMinMaxIDLimitBadUser verifies that
// Client.User.CheckinsMinMaxIDLimit returns an error when an invalid user
// is queried.
func TestClientUserCheckinsMinMaxIDLimitBadUser(t *testing.T) {
	c, done := userCheckinsTestClient(t, func(t *testing.T, w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(invalidUserErrJSON)
	})
	defer done()

	_, _, err := c.User.CheckinsMinMaxIDLimit("foo", 0, math.MaxInt32, 25)
	assertInvalidUserErr(t, err)
}

// TestClientUserCheckinsMinMaxIDLimitOK verifies that Client.User.CheckinsMinMaxIDLimit
// returns a valid checkins list, when used with correct parameters.
func TestClientUserCheckinsMinMaxIDLimitOffsetLimitOK(t *testing.T) {
	var minID int
	sMinID := strconv.Itoa(minID)

	var maxID = math.MaxInt32
	sMaxID := strconv.Itoa(maxID)

	var limit = 25
	sLimit := strconv.Itoa(limit)

	username := "kriben"
	c, done := userCheckinsTestClient(t, func(t *testing.T, w http.ResponseWriter, r *http.Request) {
		path := "/v4/user/checkins/" + username + "/"
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

		w.Write(userCheckinsJSON)
	})
	defer done()

	checkins, _, err := c.User.CheckinsMinMaxIDLimit(username, minID, maxID, limit)
	if err != nil {
		t.Fatal(err)
	}

	// Check data against expected set of checkins
	assertExpectedCheckins(t, checkins)
}

// userCheckinTestClient builds upon testClient, and adds additional sanity checks
// for tests which target the user checkin API.
func userCheckinsTestClient(t *testing.T, fn func(t *testing.T, w http.ResponseWriter, r *http.Request)) (*Client, func()) {
	return testClient(t, func(t *testing.T, w http.ResponseWriter, r *http.Request) {
		// Always GET request
		method := "GET"
		if m := r.Method; m != method {
			t.Fatalf("unexpected HTTP method: %q != %q", m, method)
		}

		// Always uses specific path prefix
		prefix := "/v4/user/checkins/"
		if p := r.URL.Path; !strings.HasPrefix(p, prefix) {
			t.Fatalf("unexpected HTTP path prefix: %q != %q", p, prefix)
		}

		// Guard against panics
		if fn != nil {
			fn(t, w, r)
		}
	})
}
