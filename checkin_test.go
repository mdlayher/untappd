package untappd_test

import (
	"testing"

	"github.com/mdlayher/untappd"
)

// assertExpectedCheckins validates a set of mock checkins from a test function
// against the expected checkin JSON used when testing this package.
func assertExpectedCheckins(t *testing.T, checkins []*untappd.Checkin) {
	expected := []*untappd.Checkin{
		&untappd.Checkin{
			ID:      137117722,
			Comment: "When in Rome..",
			Beer: &untappd.Beer{
				Name:  "Brooklyn Bowl Pale Ale",
				Style: "American Pale Ale",
			},
			Brewery: &untappd.Brewery{
				Name: "Kelso of Brooklyn",
			},
			User: &untappd.User{
				UserName: "gregavola",
			},
			Badges: []*untappd.Badge{{
				Name: "Taste the Music",
			}},
			Toasts: []*untappd.Toast{{
				ID: 1,
				User: &untappd.User{
					UserName: "gregavola",
				},
			}},
			Comments: []*untappd.Comment{{
				ID:      1,
				Comment: "hello, world",
				User: &untappd.User{
					UserName: "gregavola",
				},
			}},
		},
	}

	for i := range checkins {
		if checkins[i].ID != expected[i].ID {
			t.Fatalf("unexpected checkin ID: %d != %d", checkins[i].ID, expected[i].ID)
		}
		if checkins[i].Comment != expected[i].Comment {
			t.Fatalf("unexpected checkin Comment: %q != %q", checkins[i].Comment, expected[i].Comment)
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
		if checkins[i].Badges[0].Name != expected[i].Badges[0].Name {
			t.Fatalf("unexpected checkin Badge.Name: %q != %q", checkins[i].Badges[0].Name, expected[i].Badges[0].Name)
		}
		if checkins[i].Toasts[0].ID != expected[i].Toasts[0].ID {
			t.Fatalf("unexpected checkin Toast.ID: %d != %d", checkins[i].Toasts[0].ID, expected[i].Toasts[0].ID)
		}
		if checkins[i].Toasts[0].User.UserName != expected[i].Toasts[0].User.UserName {
			t.Fatalf("unexpected checkin Toast.User.UserName: %q != %q", checkins[i].Toasts[0].User.UserName, expected[i].Toasts[0].User.UserName)
		}
		if checkins[i].Comments[0].ID != expected[i].Comments[0].ID {
			t.Fatalf("unexpected checkin Toast.ID: %d != %d", checkins[i].Comments[0].ID, expected[i].Comments[0].ID)
		}
		if checkins[i].Comments[0].Comment != expected[i].Comments[0].Comment {
			t.Fatalf("unexpected checkin Toast.Comment: %q != %q", checkins[i].Comments[0].Comment, expected[i].Comments[0].Comment)
		}
		if checkins[i].Comments[0].User.UserName != expected[i].Comments[0].User.UserName {
			t.Fatalf("unexpected checkin Toast.User.UserName: %q != %q", checkins[i].Comments[0].User.UserName, expected[i].Comments[0].User.UserName)
		}
	}
}

// Canned checkins JSON response, taken from documentation: https://untappd.com/api/docs#useractivityfeed
// All checkin responses are in this format, and it is used throughout various
// Checkin method tests.
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
            "count": 1,
            "items": [
              {
                "comment_id": 1,
                "comment": "hello, world",
                "user": {
                  "user_name": "gregavola"
                }
              }
            ]
          },
          "toasts": {
            "total_count": 0,
            "count": 1,
            "auth_toast": false,
            "items": [
              {
                "like_id": 1,
                "user": {
                  "user_name": "gregavola"
                }
              }
            ]
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
