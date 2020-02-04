package untappd_test

import (
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/mdlayher/untappd"
)

// TestClientLocalCheckinsOK verifies that Client.Local.Checkins always sets the
// appropriate default values.
func TestClientLocalCheckinsOK(t *testing.T) {
	lat := "0"
	lng := "0"
	limit := "25"
	radius := "25"
	distance := untappd.DistanceMiles

	c, done := localCheckinsTestClient(t, func(t *testing.T, w http.ResponseWriter, r *http.Request) {
		assertParameters(t, r, url.Values{
			"lat":       []string{lat},
			"lng":       []string{lng},
			"limit":     []string{limit},
			"radius":    []string{radius},
			"dist_pref": []string{string(distance)},
		})

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

	_, _, err := c.Local.CheckinsMinMaxIDLimitRadius(untappd.LocalCheckinsRequest{
		Latitude:  0.0,
		Longitude: 0.0,
	})
	assertInvalidLocalErr(t, err)
}

// TestClientLocalCheckinsMinMaxIDLimitRadiusOK verifies that Client.Local.CheckinsMinMaxIDLimitRadius
// returns a valid checkins list, when used with correct parameters.
func TestClientLocalCheckinsMinMaxIDLimitRadiusOffsetLimitOK(t *testing.T) {
	var lat = 1.00
	sLat := untappd.FormatFloat(lat)

	var lng = -1.00
	sLng := untappd.FormatFloat(lng)

	var minID = 1
	sMinID := strconv.Itoa(minID)

	var maxID = math.MaxInt32
	sMaxID := strconv.Itoa(maxID)

	var limit = 25
	sLimit := strconv.Itoa(limit)

	var radius = 25
	sRadius := strconv.Itoa(radius)

	var distance = untappd.DistanceMiles

	c, done := localCheckinsTestClient(t, func(t *testing.T, w http.ResponseWriter, r *http.Request) {
		assertParameters(t, r, url.Values{
			"lat":       []string{sLat},
			"lng":       []string{sLng},
			"min_id":    []string{sMinID},
			"max_id":    []string{sMaxID},
			"limit":     []string{sLimit},
			"radius":    []string{sRadius},
			"dist_pref": []string{string(distance)},
		})

		// JSON is in same format as /v4/user/checkins, so we can
		// reuse it here
		w.Write(userCheckinsJSON)
	})
	defer done()

	checkins, _, err := c.Local.CheckinsMinMaxIDLimitRadius(untappd.LocalCheckinsRequest{
		Latitude:  lat,
		Longitude: lng,

		MinID: minID,
		MaxID: maxID,

		Limit: limit,

		Radius: radius,
		Units:  distance,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Check data against expected set of checkins
	assertExpectedCheckins(t, checkins)
}

// localCheckinTestClient builds upon testClient, and adds additional sanity checks
// for tests which target the local checkin API.
func localCheckinsTestClient(t *testing.T, fn func(t *testing.T, w http.ResponseWriter, r *http.Request)) (*untappd.Client, func()) {
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
