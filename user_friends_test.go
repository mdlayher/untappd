package untappd

import (
	"net/http"
	"strconv"
	"strings"
	"testing"
)

// TestClientUserFriendsOK verifies that Client.User.Friends always sets the
// appropriate default offset and limit values.
func TestClientUserFriendsOK(t *testing.T) {
	offset := "0"
	limit := "25"

	c, done := userFriendsTestClient(t, func(t *testing.T, w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()

		if o := q.Get("offset"); o != offset {
			t.Fatalf("unexpected offset parameter: %s != %s", o, offset)
		}
		if l := q.Get("limit"); l != limit {
			t.Fatalf("unexpected limit parameter: %s != %s", l, limit)
		}

		// Empty JSON response since we already passed checks
		w.Write([]byte("{}"))
	})
	defer done()

	if _, _, err := c.User.Friends("foo"); err != nil {
		t.Fatal(err)
	}
}

// TestClientUserFriendsOffsetLimitBadUser verifies that Client.User.FriendsOffsetLimit
// returns an error when an invalid user is queried.
func TestClientUserFriendsOffsetLimitBadUser(t *testing.T) {
	c, done := userFriendsTestClient(t, func(t *testing.T, w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write(invalidUserErrJSON)
	})
	defer done()

	_, _, err := c.User.FriendsOffsetLimit("foo", 0, 25)
	assertInvalidUserErr(t, err)
}

// TestClientUserFriendsOffsetLimitOK verifies that Client.User.FriendsOffsetLimit
// returns a valid friends list, when used with correct parameters.
func TestClientUserFriendsOffsetLimitOK(t *testing.T) {
	var offset int
	sOffset := strconv.Itoa(offset)

	var limit = 25
	sLimit := strconv.Itoa(limit)

	username := "mdlayher"
	c, done := userFriendsTestClient(t, func(t *testing.T, w http.ResponseWriter, r *http.Request) {
		path := "/v4/user/friends/" + username + "/"
		if p := r.URL.Path; p != path {
			t.Fatalf("unexpected URL path: %q != %q", p, path)
		}

		q := r.URL.Query()

		if o := q.Get("offset"); o != sOffset {
			t.Fatalf("unexpected offset parameter: %s != %s", o, sOffset)
		}
		if l := q.Get("limit"); l != sLimit {
			t.Fatalf("unexpected limit parameter: %s != %s", l, sLimit)
		}

		w.Write(userFriendsJSON)
	})
	defer done()

	friends, _, err := c.User.FriendsOffsetLimit(username, offset, limit)
	if err != nil {
		t.Fatal(err)
	}

	expected := []*User{
		&User{
			UID:      123456,
			UserName: "XXXXXX",
		},
		&User{
			UID:      789123,
			UserName: "YYYYYY",
		},
	}

	for i := range friends {
		if friends[i].UID != expected[i].UID {
			t.Fatalf("unexpected friend UID: %d != %d", friends[i].UID, expected[i].UID)
		}
		if friends[i].UserName != expected[i].UserName {
			t.Fatalf("unexpected friend UserName: %q != %q", friends[i].UserName, expected[i].UserName)
		}
	}
}

// userFriendsTestClient builds upon testClient, and adds additional sanity checks
// for tests which target the user friends API.
func userFriendsTestClient(t *testing.T, fn func(t *testing.T, w http.ResponseWriter, r *http.Request)) (*Client, func()) {
	return testClient(t, func(t *testing.T, w http.ResponseWriter, r *http.Request) {
		// Always GET request
		method := "GET"
		if m := r.Method; m != method {
			t.Fatalf("unexpected HTTP method: %q != %q", m, method)
		}

		// Always uses specific path prefix
		prefix := "/v4/user/friends/"
		if p := r.URL.Path; !strings.HasPrefix(p, prefix) {
			t.Fatalf("unexpected HTTP path prefix: %q != %q", p, prefix)
		}

		// Guard against panics
		if fn != nil {
			fn(t, w, r)
		}
	})
}

// Canned user friends JSON response, taken from documentation: https://untappd.com/api/docs#userfriends
// Slight modifications made to add multiple users to items list
var userFriendsJSON = []byte(`{
  "meta": {
    "code": 200,
    "response_time": {
      "time": 0,
      "measure": "seconds"
    }
  },
  "notifications": {},
  "response": {
  "count": 2,
  "items": [{
    "friendship_hash": "143242342325453",
    "created_at": "Sun, 23 Nov 2014 04:33:12 +0000",
    "user": {
      "uid": 123456,
      "user_name": "XXXXXX",
      "location": "XXXXX",
      "bio": "BioHere",
      "is_supporter": 1,
      "first_name": "XXXXXX",
      "last_name": "XXXXX",
      "relationship": "friends",
      "user_avatar": "https://d1c8v1qci5en44.cloudfront.net/profile/844124b9ff349b226018dd7bf549f052_thumb.jpg"
    },
    "mutual_friends": {
      "count": 0,
      "items": []
    }
  },
  {
    "friendship_hash": "143242342325453",
    "created_at": "Sun, 23 Nov 2014 04:33:12 +0000",
    "user": {
      "uid": 789123,
      "user_name": "YYYYYY",
      "location": "YYYYY",
      "bio": "BioHere",
      "is_supporter": 1,
      "first_name": "YYYYYY",
      "last_name": "YYYYY",
      "relationship": "friends",
      "user_avatar": "https://d1c8v1qci5en44.cloudfront.net/profile/844124b9ff349b226018dd7bf549f052_thumb.jpg"
    },
    "mutual_friends": {
      "count": 0,
      "items": []
    }
  }]
}}`)
