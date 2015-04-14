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

// Canned user checkins JSON response, taken from documentation: https://untappd.com/api/docs#useractivityfeed
var userCheckinsJSON = []byte(`{
  "meta": {
    "code": 200,
    "response_time": {
      "time": 0.841,
      "measure": "seconds"
    },
    "init_time": {
      "time": 0.001,
      "measure": "seconds"
    }
  },
  "notifications": [

  ],
  "response": {
    "pagination": {
      "since_url": "https://api.untappd.com/v4/user/checkins/gregavola?min_id=171626491",
      "next_url": "https://api.untappd.com/v4/user/checkins/gregavola?max_id=161830366",
      "max_id": 161830366
    },
    "checkins": {
      "count": 1,
      "items": [
        {
          "checkin_id": 137117722,
          "created_at": "Sat, 13 Dec 2014 19:15:38 +0000",
          "checkin_comment": "When in Rome..",
          "rating_score": 3,
          "user": {
            "uid": 1,
            "user_name": "gregavola",
            "first_name": "Greg",
            "last_name": "Avola",
            "location": "New York, NY",
            "is_supporter": 1,
            "url": "http://gregavola.com",
            "bio": "Co-Founder and CTO of Untappd, Web Developer, Beer Drinker & Community Guy",
            "relationship": "self",
            "user_avatar": "https://gravatar.com/avatar/0c6922e238dae5cccce96a32889fc911?size=100&d=htt\u202644.cloudfront.net%2Fsite%2Fassets%2Fimages%2Fdefault_avatar_v2.jpg%3Fv%3D1",
            "is_private": 0,
            "contact": {
              "foursquare": 195741,
              "twitter": "gregavola",
              "facebook": 18603076
            }
          },
          "beer": {
            "bid": 7481,
            "beer_name": "Brooklyn Bowl Pale Ale",
            "beer_label": "https://d1c8v1qci5en44.cloudfront.net/site/assets/images/temp/badge-beer-default.png",
            "beer_style": "American Pale Ale",
            "beer_abv": 0,
            "auth_rating": 0,
            "wish_list": false,
            "beer_active": 1
          },
          "brewery": {
            "brewery_id": 1954,
            "brewery_name": "Kelso of Brooklyn",
            "brewery_slug": "kelso-of-brooklyn",
            "brewery_label": "https://d1c8v1qci5en44.cloudfront.net/site/brewery_logos/brewery-KelsoofBrooklyn_1954.jpeg",
            "country_name": "United States",
            "contact": {
              "twitter": "KelsoBeer",
              "facebook": "",
              "instagram": "",
              "url": "http://www.kelsoofbrooklyn.com/"
            },
            "location": {
              "brewery_city": "Brooklyn",
              "brewery_state": "NY",
              "lat": 40.6823,
              "lng": -73.9656
            },
            "brewery_active": 1
          },
          "venue": {
            "venue_id": 2141,
            "venue_name": "Brooklyn Bowl",
            "primary_category": "Arts & Entertainment",
            "parent_category_id": "4d4b7104d754a06370d81259",
            "categories": {
              "count": 3,
              "items": [
                {
                  "category_name": "Bowling Alley",
                  "category_id": "4bf58dd8d48988d1e4931735",
                  "is_primary": true
                },
                {
                  "category_name": "Music Venue",
                  "category_id": "4bf58dd8d48988d1e5931735",
                  "is_primary": false
                },
                {
                  "category_name": "Bar",
                  "category_id": "4bf58dd8d48988d116941735",
                  "is_primary": false
                }
              ]
            },
            "location": {
              "venue_address": "61 Wythe Ave",
              "venue_city": "Brooklyn",
              "venue_state": "NY",
              "venue_country": "United States",
              "lat": 40.7219,
              "lng": -73.9575
            },
            "contact": {
              "twitter": "@brooklynbowl",
              "venue_url": "http://www.brooklynbowl.com"
            },
            "public_venue": true,
            "foursquare": {
              "foursquare_id": "4a1afeb7f964a520b77a1fe3",
              "foursquare_url": "http://4sq.com/3fjtlA"
            },
            "venue_icon": {
              "sm": "https://ss3.4sqi.net/img/categories_v2/arts_entertainment/bowling_bg_64.png",
              "md": "https://ss3.4sqi.net/img/categories_v2/arts_entertainment/bowling_bg_88.png",
              "lg": "https://ss3.4sqi.net/img/categories_v2/arts_entertainment/bowling_bg_88.png"
            }
          },
          "comments": {
            "total_count": 0,
            "count": 0,
            "items": []
          },
          "toasts": {
            "total_count": 0,
            "count": 0,
            "auth_toast": false,
            "items": []
          },
          "media": {
            "count": 0,
            "items": []
          },
          "source": {
            "app_name": "Untappd for iPhone - (V2)",
            "app_website": "http://untpd.it/iphoneapp"
          },
          "badges": {
            "count": 1,
            "items": [
              {
                "badge_id": 189,
                "user_badge_id": 39410316,
                "badge_name": "Taste the Music",
                "badge_description": "Badge Description Here",
                "created_at": "Sat, 13 Dec 2014 19:15:41 +0000",
                "badge_image": {
                  "sm": "https://d1c8v1qci5en44.cloudfront.net/badges/bdg_ConcertVenue_sm.jpg",
                  "md": "https://d1c8v1qci5en44.cloudfront.net/badges/bdg_ConcertVenue_md.jpg",
                  "lg": "https://d1c8v1qci5en44.cloudfront.net/badges/bdg_ConcertVenue_lg.jpg"
                }
              }
            ]
          }
        }
      ]
    }
  }
}`)
