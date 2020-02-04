package untappd_test

import (
	"fmt"
	"testing"

	"github.com/mdlayher/untappd"
)

// TestSorts verifies that every Sort type is present in the output of Sorts.
func TestSorts(t *testing.T) {
	for _, s := range []untappd.Sort{
		untappd.SortDate,
		untappd.SortCheckin,
		untappd.SortHighestRated,
		untappd.SortLowestRated,
		untappd.SortUserHighestRated,
		untappd.SortUserLowestRated,
		untappd.SortHighestABV,
		untappd.SortLowestABV,
	} {
		t.Run(fmt.Sprintf("%q", s), func(t *testing.T) {
			for _, ss := range untappd.Sorts() {
				if s == ss {
					return
				}
			}
			t.Fatalf("unknown Sort type: %q", s)
		})
	}
}
