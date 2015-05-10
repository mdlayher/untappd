package untappd

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"
)

// TestClientAuthCheckinOK verifies that Client.Auth.Checkin always sets the
// appropriate POST body parameters for a valid checkin.
func TestClientAuthCheckinOK(t *testing.T) {
	beerID := 1
	sBeerID := strconv.Itoa(beerID)

	timezone, offset := time.Now().Zone()
	offset = offset / 60 / 60
	sOffset := strconv.Itoa(offset)

	foursquareID := "ABCDEF"

	latitude := 1.0
	sLatitude := formatFloat(latitude)
	longitude := 1.0
	sLongitude := formatFloat(longitude)

	comment := "hello world"

	rating := 3.5
	sRating := formatFloat(rating)

	facebook := true
	twitter := true
	foursquare := true

	c, done := authCheckinTestClient(t, func(t *testing.T, w http.ResponseWriter, r *http.Request) {
		assertBodyParameters(t, r, url.Values{
			"bid":           []string{sBeerID},
			"gmt_offset":    []string{sOffset},
			"timezone":      []string{timezone},
			"foursquare_id": []string{foursquareID},
			"geolat":        []string{sLatitude},
			"geolng":        []string{sLongitude},
			"shout":         []string{comment},
			"rating":        []string{sRating},
			"facebook":      []string{"on"},
			"twitter":       []string{"on"},
			"foursquare":    []string{"on"},
		})

		// Empty JSON response since we already passed checks
		w.Write([]byte("{}"))
	})
	defer done()

	if _, _, err := c.Auth.Checkin(CheckinRequest{
		BeerID:    beerID,
		GMTOffset: offset,
		TimeZone:  timezone,

		FoursquareID: foursquareID,
		Latitude:     latitude,
		Longitude:    longitude,
		Comment:      comment,
		Rating:       rating,
		Facebook:     facebook,
		Twitter:      twitter,
		Foursquare:   foursquare,
	}); err != nil {
		t.Fatal(err)
	}
}

// TestClientAuthCheckinBadBeerID verifies that Client.Auth.Checkin returns an
// error when an invalid beer ID is checked-in.
func TestClientAuthCheckinBadBeerID(t *testing.T) {
	beerID := -1

	c, done := authCheckinTestClient(t, func(t *testing.T, w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(invalidCheckinErrJSON)
	})
	defer done()

	_, _, err := c.Auth.Checkin(CheckinRequest{
		BeerID: beerID,
	})
	assertInvalidCheckinErr(t, err)
}

// authCheckinTestClient builds upon testClient, and adds additional sanity checks
// for tests which target the Check-in API.
func authCheckinTestClient(t *testing.T, fn func(t *testing.T, w http.ResponseWriter, r *http.Request)) (*Client, func()) {
	return testClient(t, func(t *testing.T, w http.ResponseWriter, r *http.Request) {
		// Always POST request
		method := "POST"
		if m := r.Method; m != method {
			t.Fatalf("unexpected HTTP method: %q != %q", m, method)
		}

		// Always uses specific path prefix
		prefix := "/v4/checkin/add"
		if p := r.URL.Path; !strings.HasPrefix(p, prefix) {
			t.Fatalf("unexpected HTTP path prefix: %q != %q", p, prefix)
		}

		// Guard against panics
		if fn != nil {
			fn(t, w, r)
		}
	})
}
