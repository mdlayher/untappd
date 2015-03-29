package untappd

import (
	"net/http"
	"strconv"
	"strings"
	"testing"
)

// TestClientUserBadgesOK verifies that Client.User.Badges always sets the
// appropriate default offset value.
func TestClientUserBadgesOK(t *testing.T) {
	offset := "0"

	c, done := userBadgesTestClient(t, func(t *testing.T, w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()

		if o := q.Get("offset"); o != offset {
			t.Fatalf("unexpected offset parameter: %s != %s", o, offset)
		}

		// Empty JSON response since we already passed checks
		w.Write([]byte("{}"))
	})
	defer done()

	if _, _, err := c.User.Badges("foo"); err != nil {
		t.Fatal(err)
	}
}

// TestClientUserBadgesOffsetBadUser verifies that Client.User.BadgesOffset
// returns an error when an invalid user is queried.
func TestClientUserBadgesOffsetBadUser(t *testing.T) {
	c, done := userBadgesTestClient(t, func(t *testing.T, w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write(invalidUserErrJSON)
	})
	defer done()

	_, _, err := c.User.BadgesOffset("foo", 50)
	assertInvalidUserErr(t, err)
}

// TestClientUserBadgesOffsetOK verifies that Client.User.BadgesOffset
// returns a valid badges list, when used with correct parameters.
func TestClientUserBadgesOffsetOK(t *testing.T) {
	var offset int
	sOffset := strconv.Itoa(int(offset))

	username := "mdlayher"
	c, done := userBadgesTestClient(t, func(t *testing.T, w http.ResponseWriter, r *http.Request) {
		path := "/v4/user/badges/" + username + "/"
		if p := r.URL.Path; p != path {
			t.Fatalf("unexpected URL path: %q != %q", p, path)
		}

		q := r.URL.Query()

		if o := q.Get("offset"); o != sOffset {
			t.Fatalf("unexpected offset parameter: %s != %s", o, sOffset)
		}

		w.Write(userBadgesJSON)
	})
	defer done()

	badges, _, err := c.User.BadgesOffset(username, offset)
	if err != nil {
		t.Fatal(err)
	}

	expected := []*Badge{
		&Badge{
			ID:   189,
			Name: "Taste the Music",
		},
		&Badge{
			ID:   190,
			Name: "Oberon (2015)",
		},
	}

	for i := range badges {
		if badges[i].ID != expected[i].ID {
			t.Fatalf("unexpected badge ID: %d != %d", badges[i].ID, expected[i].ID)
		}
		if badges[i].Name != expected[i].Name {
			t.Fatalf("unexpected badge Name: %q != %q", badges[i].Name, expected[i].Name)
		}
	}
}

// userBadgesTestClient builds upon testClient, and adds additional sanity checks
// for tests which target the user badges API.
func userBadgesTestClient(t *testing.T, fn func(t *testing.T, w http.ResponseWriter, r *http.Request)) (*Client, func()) {
	return testClient(t, func(t *testing.T, w http.ResponseWriter, r *http.Request) {
		// Always GET request
		method := "GET"
		if m := r.Method; m != method {
			t.Fatalf("unexpected HTTP method: %q != %q", m, method)
		}

		// Always uses specific path prefix
		prefix := "/v4/user/badges/"
		if p := r.URL.Path; !strings.HasPrefix(p, prefix) {
			t.Fatalf("unexpected HTTP path prefix: %q != %q", p, prefix)
		}

		// Guard against panics
		if fn != nil {
			fn(t, w, r)
		}
	})
}

// Canned user badges JSON response, taken from documentation: https://untappd.com/api/docs#userbadges
// Slight modifications made to add multiple badges to items list
var userBadgesJSON = []byte(`{
"response": {
  "type": "earned",
  "sort": "all",
  "count": 2,
  "items": [
  {
    "user_badge_id": 39410316,
    "badge_id": 189,
    "checkin_id": 137117722,
    "badge_name": "Taste the Music",
    "badge_description": "Description Here",
    "badge_active_status": 1,
    "media": {
      "badge_image_sm": "https://d1c8v1qci5en44.cloudfront.net/badges/bdg_ConcertVenue_sm.jpg",
      "badge_image_md": "https://d1c8v1qci5en44.cloudfront.net/badges/bdg_ConcertVenue_md.jpg",
      "badge_image_lg": "https://d1c8v1qci5en44.cloudfront.net/badges/bdg_ConcertVenue_lg.jpg"
    },
    "created_at": "Sat, 13 Dec 2014 19:15:41 +0000",
    "is_level": true,
    "category_id": 2,
    "levels": {
      "count": 1,
      "items": [
        {
          "actual_badge_id": 189,
          "badge_id": 39410316,
          "checkin_id": 137117722,
          "badge_name": "Taste the Music",
          "badge_description": "Descriptio  here",
          "media": {
            "badge_image_sm": "https://d1c8v1qci5en44.cloudfront.net/badges/bdg_ConcertVenue_sm.jpg",
            "badge_image_md": "https://d1c8v1qci5en44.cloudfront.net/badges/bdg_ConcertVenue_md.jpg",
            "badge_image_lg": "https://d1c8v1qci5en44.cloudfront.net/badges/bdg_ConcertVenue_lg.jpg"
          },
          "created_at": "Sat, 13 Dec 2014 19:15:41 +0000"
        }
      ]
    }
  },
  {
    "user_badge_id": 39410316,
    "badge_id": 190,
    "checkin_id": 137117722,
    "badge_name": "Oberon (2015)",
    "badge_description": "Description Here",
    "badge_active_status": 1,
    "media": {
      "badge_image_sm": "https://d1c8v1qci5en44.cloudfront.net/badges/bdg_ConcertVenue_sm.jpg",
      "badge_image_md": "https://d1c8v1qci5en44.cloudfront.net/badges/bdg_ConcertVenue_md.jpg",
      "badge_image_lg": "https://d1c8v1qci5en44.cloudfront.net/badges/bdg_ConcertVenue_lg.jpg"
    },
    "created_at": "Sat, 13 Dec 2014 19:15:41 +0000",
    "is_level": true,
    "category_id": 2,
    "levels": {
      "count": 1,
      "items": [
        {
          "actual_badge_id": 189,
          "badge_id": 39410316,
          "checkin_id": 137117722,
          "badge_name": "Taste the Music",
          "badge_description": "Descriptio  here",
          "media": {
            "badge_image_sm": "https://d1c8v1qci5en44.cloudfront.net/badges/bdg_ConcertVenue_sm.jpg",
            "badge_image_md": "https://d1c8v1qci5en44.cloudfront.net/badges/bdg_ConcertVenue_md.jpg",
            "badge_image_lg": "https://d1c8v1qci5en44.cloudfront.net/badges/bdg_ConcertVenue_lg.jpg"
          },
          "created_at": "Sat, 13 Dec 2014 19:15:41 +0000"
        }
      ]
    }
  }
  ]
}
}`)
