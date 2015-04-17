package untappd

import (
	"math"
	"net/http"
	"strconv"
	"strings"
	"testing"
)

// TestClientLocalCheckinsOK verifies that Client.Local.Checkins always sets the
// appropriate default values.
func TestClientLocalCheckinsOK(t *testing.T) {
	lat := "0"
	lng := "0"
	minID := "0"
	maxID := strconv.Itoa(math.MaxInt32)
	limit := "25"
	radius := "25"
	distance := DistanceMiles

	c, done := localCheckinsTestClient(t, func(t *testing.T, w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()

		if l := q.Get("lat"); l != lat {
			t.Fatalf("unexpected lat parameter: %s != %s", l, lat)
		}
		if l := q.Get("lng"); l != lng {
			t.Fatalf("unexpected lng parameter: %s != %s", l, lng)
		}
		if i := q.Get("min_id"); i != minID {
			t.Fatalf("unexpected min_id parameter: %s != %s", i, minID)
		}
		if i := q.Get("max_id"); i != maxID {
			t.Fatalf("unexpected max_id parameter: %s != %s", i, maxID)
		}
		if l := q.Get("limit"); l != limit {
			t.Fatalf("unexpected limit parameter: %s != %s", l, limit)
		}
		if r := q.Get("radius"); r != radius {
			t.Fatalf("unexpected radius parameter: %s != %s", r, radius)
		}
		if d := q.Get("dist_pref"); d != string(distance) {
			t.Fatalf("unexpected dist_pref parameter: %s != %s", d, distance)
		}

		// Empty JSON response since we already passed checks
		w.Write([]byte("{}"))
	})
	defer done()

	if _, _, err := c.Local.Checkins(0, 0); err != nil {
		t.Fatal(err)
	}
}

// TestClientLocalCheckinsMinMaxIDLimitRadiusBadCoordinates verifies that
// Client.Local.CheckinsMinMaxIDLimitRadius returns an error when zero
// latitude and longitude coordinates are queried.
func TestClientLocalCheckinsMinMaxIDLimitRadiusBadLocal(t *testing.T) {
	c, done := localCheckinsTestClient(t, func(t *testing.T, w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(invalidLocalErrJSON)
	})
	defer done()

	_, _, err := c.Local.CheckinsMinMaxIDLimitRadius(0.0, 0.0, 0, math.MaxInt32, 25, 25, DistanceMiles)
	assertInvalidLocalErr(t, err)
}

// TestClientLocalCheckinsMinMaxIDLimitRadiusOK verifies that Client.Local.CheckinsMinMaxIDLimitRadius
// returns a valid checkins list, when used with correct parameters.
func TestClientLocalCheckinsMinMaxIDLimitRadiusOffsetLimitOK(t *testing.T) {
	var lat = 1.00
	sLat := strconv.FormatFloat(lat, 'f', -1, 64)

	var lng = -1.00
	sLng := strconv.FormatFloat(lng, 'f', -1, 64)

	var minID int
	sMinID := strconv.Itoa(minID)

	var maxID = math.MaxInt32
	sMaxID := strconv.Itoa(maxID)

	var limit = 25
	sLimit := strconv.Itoa(limit)

	var radius = 25
	sRadius := strconv.Itoa(radius)

	var distance = DistanceMiles

	c, done := localCheckinsTestClient(t, func(t *testing.T, w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()

		if l := q.Get("lat"); l != sLat {
			t.Fatalf("unexpected lat parameter: %s != %s", l, sLat)
		}
		if l := q.Get("lng"); l != sLng {
			t.Fatalf("unexpected lng parameter: %s != %s", l, sLng)
		}
		if i := q.Get("min_id"); i != sMinID {
			t.Fatalf("unexpected min_id parameter: %s != %s", i, sMinID)
		}
		if i := q.Get("max_id"); i != sMaxID {
			t.Fatalf("unexpected max_id parameter: %s != %s", i, sMaxID)
		}
		if l := q.Get("limit"); l != sLimit {
			t.Fatalf("unexpected limit parameter: %s != %s", l, sLimit)
		}
		if r := q.Get("radius"); r != sRadius {
			t.Fatalf("unexpected radius parameter: %s != %s", r, sRadius)
		}
		if d := q.Get("dist_pref"); d != string(distance) {
			t.Fatalf("unexpected dist_pref parameter: %s != %s", d, distance)
		}

		// JSON is in same format as /v4/user/checkins, so we can
		// reuse it here
		w.Write(userCheckinsJSON)
	})
	defer done()

	checkins, _, err := c.Local.CheckinsMinMaxIDLimitRadius(
		lat,
		lng,
		minID,
		maxID,
		limit,
		radius,
		distance,
	)
	if err != nil {
		t.Fatal(err)
	}

	expected := []*Checkin{
		&Checkin{
			ID:      137117722,
			Comment: "When in Rome..",
			Beer: &Beer{
				Name:  "Brooklyn Bowl Pale Ale",
				Style: "American Pale Ale",
			},
			Brewery: &Brewery{
				Name: "Kelso of Brooklyn",
			},
			User: &User{
				UserName: "gregavola",
			},
		},
	}

	for i := range checkins {
		if checkins[i].ID != expected[i].ID {
			t.Fatalf("unexpected checkin ID: %d != %d", checkins[i].ID, expected[i].ID)
		}
		if checkins[i].Comment != expected[i].Comment {
			t.Fatalf("unexpected checkin Comment: %d != %d", checkins[i].Comment, expected[i].Comment)
		}
		if checkins[i].Beer.Name != expected[i].Beer.Name {
			t.Fatalf("unexpected beer Name: %q != %q", checkins[i].Beer.Name, expected[i].Beer.Name)
		}
		if checkins[i].Beer.Style != expected[i].Beer.Style {
			t.Fatalf("unexpected beer Style: %q != %q", checkins[i].Beer.Style, expected[i].Beer.Style)
		}
		if checkins[i].Brewery.Name != expected[i].Brewery.Name {
			t.Fatalf("unexpected checkin Brewery.Name: %q != %q", checkins[i].Brewery.Name, expected[i].Brewery.Name)
		}
		if checkins[i].User.UserName != expected[i].User.UserName {
			t.Fatalf("unexpected checkin User.Name: %q != %q", checkins[i].User.UserName, expected[i].User.UserName)
		}
	}
}

// localCheckinTestClient builds upon testClient, and adds additional sanity checks
// for tests which target the local checkin API.
func localCheckinsTestClient(t *testing.T, fn func(t *testing.T, w http.ResponseWriter, r *http.Request)) (*Client, func()) {
	return testClient(t, func(t *testing.T, w http.ResponseWriter, r *http.Request) {
		// Always GET request
		method := "GET"
		if m := r.Method; m != method {
			t.Fatalf("unexpected HTTP method: %q != %q", m, method)
		}

		// Always uses specific path prefix
		prefix := "/v4/thepub/local/"
		if p := r.URL.Path; !strings.HasPrefix(p, prefix) {
			t.Fatalf("unexpected HTTP path prefix: %q != %q", p, prefix)
		}

		// Guard against panics
		if fn != nil {
			fn(t, w, r)
		}
	})
}
