package untappd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"
)

var (
	// errInvalidBool is returned when the Untappd API returns a
	// non 0 or 1 integer for a boolean value.
	errInvalidBool = errors.New("invalid boolean value")

	// errInvalidTimeUnit is returned when the Untappd API returns an
	// unrecognized time unit.
	errInvalidTimeUnit = errors.New("invalid time unit")
)

// responseDuration implements json.Unmarshaler, so that duration responses
// in the Untappd APIv4 can be decoded directly into Go time.Duration structs.
type responseDuration time.Duration

// UnmarshalJSON implements json.Unmarshaler.
func (r *responseDuration) UnmarshalJSON(data []byte) error {
	var v struct {
		Time    float64 `json:"time"`
		Measure string  `json:"measure"`
	}

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	// Known measure strings mapped to Go parse-able equivalents
	timeUnits := map[string]string{
		"milliseconds": "ms",
		"seconds":      "s",
		"minutes":      "m",
	}

	// Parse a Go time.Duration from string
	d, err := time.ParseDuration(fmt.Sprintf("%f%s", v.Time, timeUnits[v.Measure]))
	if err != nil && strings.Contains(err.Error(), "time: missing unit in duration") {
		return errInvalidTimeUnit
	}

	*r = responseDuration(d)
	return err
}

// responseTime implements json.Unmarshaler, so that timestamp responses
// in the Untappd APIv4 can be decoded directly into Go time.Time structs.
type responseTime time.Time

// UnmarshalJSON implements json.Unmarshaler.
func (r *responseTime) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	// Parse a Go time.Time from string
	t, err := time.Parse(time.RFC1123Z, v)
	if err != nil {
		return err
	}

	*r = responseTime(t)
	return nil
}

// responseURL implements json.Unmarshaler, so that URL string responses
// in the Untappd APIv4 can be decoded directly into Go *url.URL structs.
type responseURL url.URL

// UnmarshalJSON implements json.Unmarshaler.
func (r *responseURL) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	u, err := url.Parse(v)
	if err != nil {
		return err
	}

	*r = responseURL(*u)
	return nil
}

// responseBool implements json.Unmarshaler, so that integer 0 or 1 responses
// in the Untappd APIv4 can be decoded directly into Go boolean values.
type responseBool bool

// UnmarshalJSON implements json.Unmarshaler.
func (r *responseBool) UnmarshalJSON(data []byte) error {
	var v int
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	switch v {
	case 0:
		*r = false
	case 1:
		*r = true
	default:
		return errInvalidBool
	}

	return nil
}

// responseBadgeLevels implements json.Unmarshaler, so that an empty array on
// a badge with no levels can be appropriately handled.
type responseBadgeLevels struct {
	Count int
	Items []*rawBadge
}

// UnmarshalJSON implements json.Unmarshaler.
func (r *responseBadgeLevels) UnmarshalJSON(data []byte) error {
	// If no levels exist for a badge, the API returns an empty array instead
	// of a nil or empty object.  This method works around that.
	if bytes.Equal(data, []byte("[]")) {
		return nil
	}

	var v struct {
		Count int         `json:"count"`
		Items []*rawBadge `json:"items"`
	}

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	r.Count = v.Count
	r.Items = v.Items

	return nil
}
