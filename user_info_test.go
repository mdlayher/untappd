package untappd_test

import (
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/mdlayher/untappd"
)

// TestClientUserInfoBadUser verifies that Client.User.Info returns an error when
// an invalid user is queried.
func TestClientUserInfoBadUser(t *testing.T) {
	username := "mdlayher"
	c, done := userInfoTestClient(t, func(t *testing.T, w http.ResponseWriter, r *http.Request) {
		path := "/v4/user/info/" + username + "/"
		if p := r.URL.Path; p != path {
			t.Fatalf("unexpected URL path: %q != %q", p, path)
		}

		w.WriteHeader(http.StatusNotFound)
		w.Write(invalidUserErrJSON)
	})
	defer done()

	_, _, err := c.User.Info(username, false)
	assertInvalidUserErr(t, err)
}

// TestClientUserInfoCompactOK verifies that Client.User.Info properly requests compact
// user output.
func TestClientUserInfoCompactOK(t *testing.T) {
	c, done := userInfoTestClient(t, func(t *testing.T, w http.ResponseWriter, r *http.Request) {
		assertParameters(t, r, url.Values{
			"compact": []string{"true"},
		})

		// In the future, we may return compact canned user data here.
		// For now, write a mostly empty JSON object is enough to get
		// test coverage.
		w.Write([]byte(`{"response":{"user":{"id":1}}}`))
	})
	defer done()

	if _, _, err := c.User.Info("test", true); err != nil {
		t.Fatal(err)
	}
}

// TestClientUserInfoOK verifies that Client.User.Info returns a valid user when
// provided with correct input parameters.
func TestClientUserInfoOK(t *testing.T) {
	username := "gregavola"
	c, done := userInfoTestClient(t, func(t *testing.T, w http.ResponseWriter, r *http.Request) {
		path := "/v4/user/info/" + username + "/"
		if p := r.URL.Path; p != path {
			t.Fatalf("unexpected URL path: %q != %q", p, path)
		}

		w.Write(gregavolaUserJSON)
	})
	defer done()

	u, _, err := c.User.Info(username, false)
	if err != nil {
		t.Fatal(err)
	}

	if id := u.UID; id != 1 {
		t.Fatalf("unexpected UID: %d != %d", id, 1)
	}
	if id := u.ID; id != 1 {
		t.Fatalf("unexpected ID: %d != %d", id, 1)
	}
	if u := u.UserName; u != username {
		t.Fatalf("unexpected username: %q != %q", u, username)
	}
}

// userInfoTestClient builds upon testClient, and adds additional sanity checks
// for tests which target the user info API.
func userInfoTestClient(t *testing.T, fn func(t *testing.T, w http.ResponseWriter, r *http.Request)) (*untappd.Client, func()) {
	return testClient(t, func(t *testing.T, w http.ResponseWriter, r *http.Request) {
		// Always GET request
		method := "GET"
		if m := r.Method; m != method {
			t.Fatalf("unexpected HTTP method: %q != %q", m, method)
		}

		// Always uses specific path prefix
		prefix := "/v4/user/info/"
		if p := r.URL.Path; !strings.HasPrefix(p, prefix) {
			t.Fatalf("unexpected HTTP path prefix: %q != %q", p, prefix)
		}

		// Guard against panics
		if fn != nil {
			fn(t, w, r)
		}
	})
}

// Canned user JSON response, taken from documentation: https://untappd.com/api/docs#userinfo
var gregavolaUserJSON = []byte(`
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
  "user": {
    "uid": 1,
    "id": 1,
    "user_name": "gregavola",
    "first_name": "Greg",
    "last_name": "Avola",
    "user_avatar": "https://gravatar.com/avatar/0c6922e238dae5cccce96a32889fc911?size=100&d=htt…44.cloudfront.net%2Fsite%2Fassets%2Fimages%2Fdefault_avatar_v2.jpg%3Fv%3D1",
    "user_avatar_hd": "https://gravatar.com/avatar/0c6922e238dae5cccce96a32889fc911?size=125&d=htt…44.cloudfront.net%2Fsite%2Fassets%2Fimages%2Fdefault_avatar_v2.jpg%3Fv%3D1",
    "user_cover_photo": "https://untappd.s3.amazonaws.com/coverphoto/933f9eebffb9151299188512cbd5981b.jpg",
    "user_cover_photo_offset": 214,
    "is_private": 0,
    "location": "New York, NY",
    "url": "http://gregavola.com",
    "bio": "Co-Founder and CTO of Untappd, Web Developer, Beer Drinker & Community Guy",
    "is_supporter": 1,
    "relationship": "self",
    "untappd_url": "http://untappd.com/user/gregavola",
    "account_type": "user",
    "stats": {
      "total_badges": 379,
      "total_friends": 1723,
      "total_checkins": 2197,
      "total_beers": 1187,
      "total_created_beers": 65,
      "total_followings": 176,
      "total_photos": 325
    },
    "recent_brews": {
      "count": 1,
      "items": {
        "beer": {
          "bid": 7481,
          "beer_name": "Brooklyn Bowl Pale Ale",
          "beer_label": "https://d1c8v1qci5en44.cloudfront.net/site/assets/images/temp/badge-beer-default.png",
          "beer_abv": 0,
          "beer_description": "",
          "beer_style": "American Pale Ale",
          "auth_rating": 0,
          "wish_list": false
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
        }
      }
    },
    "media": {
      "count": 1,
      "items": {
        "photo_id": 24739915,
        "photo": {
          "photo_img_sm": "https://d1c8v1qci5en44.cloudfront.net/photo/2014_11_28/0a418014e88d2bc841b0f4f688714ce7_100x100.jpg",
          "photo_img_md": "https://d1c8v1qci5en44.cloudfront.net/photo/2014_11_28/0a418014e88d2bc841b0f4f688714ce7_320x320.jpg",
          "photo_img_lg": "https://d1c8v1qci5en44.cloudfront.net/photo/2014_11_28/0a418014e88d2bc841b0f4f688714ce7_640x640.jpg",
          "photo_img_og": "https://d1c8v1qci5en44.cloudfront.net/photo/2014_11_28/0a418014e88d2bc841b0f4f688714ce7_raw.jpg"
        },
        "created_at": "Fri, 28 Nov 2014 22:21:05 +0000",
        "checkin_id": 133319903,
        "user": {
          "uid": 1,
          "user_name": "gregavola",
          "location": "New York, NY",
          "bio": "Co-Founder and CTO of Untappd, Web Developer, Beer Drinker & Community Guy",
          "first_name": "Greg",
          "last_name": "Avola",
          "user_avatar": "https://gravatar.com/avatar/0c6922e238dae5cccce96a32889fc911?size=100&d=htt…44.cloudfront.net%2Fsite%2Fassets%2Fimages%2Fdefault_avatar_v2.jpg%3Fv%3D1",
          "account_type": "user",
          "url": "http://gregavola.com"
        },
        "beer": {
          "bid": 273820,
          "beer_name": "Holiday Ale",
          "beer_label": "https://d1c8v1qci5en44.cloudfront.net/site/beer_logos/beer-_273820_sm_f8f53be8552dfe14bb712208a387be.jpeg",
          "beer_abv": 7.3,
          "beer_style": "Bière de Garde",
          "beer_description": "Beer made in the traditional Biere de Garde style. Lots of grains, malty and fuller body. ",
          "auth_rating": 0,
          "wish_list": false
        },
        "brewery": {
          "brewery_id": 45815,
          "brewery_name": "Two Roads Brewing Company",
          "brewery_slug": "two-roads-brewing-company",
          "brewery_label": "https://d1c8v1qci5en44.cloudfront.net/site/brewery_logos/brewery-tworoadsbrewing_45815.jpeg",
          "country_name": "United States",
          "contact": {
            "twitter": "2RoadsBrewing",
            "facebook": "http://www.facebook.com/TwoRoadsBrewing",
            "instagram": "",
            "url": "http://www.tworoadsbrewing.com"
          },
          "location": {
            "brewery_city": "Stratford",
            "brewery_state": "CT",
            "lat": 41.1855,
            "lng": -73.1419
          },
          "brewery_active": 1
        },
        "venue": []
      }
    },
    "contact": {
      "foursquare": 195741,
      "twitter": "gregavola",
      "facebook": 18603076
    },
    "date_joined": "Wed, 07 Jul 2010 05:51:10 +0000"
  }
}}`)
