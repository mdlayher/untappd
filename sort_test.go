package untappd

import "testing"

// TestSorts verifies that every Sort type is present in the output of Sorts.
func TestSorts(t *testing.T) {
	for _, s := range []Sort{
		SortDate,
		SortCheckin,
		SortHighestRated,
		SortLowestRated,
		SortUserHighestRated,
		SortUserLowestRated,
		SortHighestABV,
		SortLowestABV,
	} {
		var found bool
		for _, ss := range Sorts() {
			if s == ss {
				found = true
				break
			}
		}
		if found {
			continue
		}

		t.Fatalf("unknown Sort type: %q", s)
	}
}
