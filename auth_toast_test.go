package untappd

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"
)

// TestClientAuthToastOK verifies that Client.Auth.Toast always sets the
// appropriate POST body parameters for a valid toast.
func TestClientAuthToastOK(t *testing.T) {
	checkinID := 1
	sCheckinID := strconv.Itoa(checkinID)

	c, done := authToastTestClient(t, func(t *testing.T, w http.ResponseWriter, r *http.Request) {
		assertBodyParameters(t, r, url.Values{
			"checkin_id": []string{sCheckinID},
		})

		// Empty JSON response since we already passed checks
		w.Write([]byte("{}"))
	})
	defer done()

	if _, err := c.Auth.Toast(ToastRequest{
		CheckinID: checkinID,
	}); err != nil {
		t.Fatal(err)
	}
}

// authToastTestClient builds upon testClient, and adds additional sanity checks
// for tests which target the Check-in API.
func authToastTestClient(t *testing.T, fn func(t *testing.T, w http.ResponseWriter, r *http.Request)) (*Client, func()) {
	return testClient(t, func(t *testing.T, w http.ResponseWriter, r *http.Request) {
		// Always POST request
		method := "POST"
		if m := r.Method; m != method {
			t.Fatalf("unexpected HTTP method: %q != %q", m, method)
		}

		// Always uses specific path prefix
		prefix := "/v4/checkin/toast"
		if p := r.URL.Path; !strings.HasPrefix(p, prefix) {
			t.Fatalf("unexpected HTTP path prefix: %q != %q", p, prefix)
		}

		// Guard against panics
		if fn != nil {
			fn(t, w, r)
		}
	})
}
