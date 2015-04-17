package untappd

// Distance is a distance unit accepted by the Untappd APIv4.
// A set of Distance constants are provided for ease of use.
type Distance string

const (
	// DistanceMiles requests a radius in miles for local checkins.
	DistanceMiles Distance = "m"

	// DistanceKilometers requests a radius in kilometers for local checkins.
	DistanceKilometers Distance = "km"
)

// LocalService is a "service" which allows access to API methods involving checkins
// in a localized area.
type LocalService struct {
	client *Client
}
