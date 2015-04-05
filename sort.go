package untappd

// Sort is a sorting method accepted by the Untappd APIv4.
// A set of Sort constants are provided for ease of use.
type Sort string

// Constants that define various methods that the Untappd APIv4 can use to
// sort Beer results.
const (
	// SortDate sorts a list of beers by most recent date checked in.
	SortDate Sort = "date"

	// SortCheckin sorts a list of beers by highest number of checkins.
	SortCheckin Sort = "checkin"

	// SortHighestRated sorts a list of beers by highest rated overall on Untappd.
	SortHighestRated Sort = "highest_rated"

	// SortLowestRated sorts a list of beers by lowest rated overall on Untappd.
	SortLowestRated Sort = "lowest_rated"

	// SortUserHighestRated sorts a list of beers by highest rated by this user
	// on Untappd.
	SortUserHighestRated Sort = "highest_rated_you"

	// SortUserLowestRated sorts a list of beers by lowest rated by this user
	// on Untappd.
	SortUserLowestRated Sort = "lowest_rated_you"

	// SortHighestABV sorts a list of beers by highest alcohol by volume on Untappd.
	SortHighestABV Sort = "highest_abv"

	// SortLowestABV sorts a list of beers by lowest alcohol by volume on Untappd.
	SortLowestABV Sort = "lowest_abv"
)

// Sorts returns a slice of all available Sort constants.
func Sorts() []Sort {
	return []Sort{
		SortDate,
		SortCheckin,
		SortHighestRated,
		SortLowestRated,
		SortUserHighestRated,
		SortUserLowestRated,
		SortHighestABV,
		SortLowestABV,
	}
}
