package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"

	"github.com/urfave/cli/v2"
	"github.com/mdlayher/untappd"
)

const (
	// appName is the name of this binary.
	appName = "untappdctl"
)

func main() {
	// Initialize new CLI app
	app := cli.NewApp()
	app.Name = appName
	app.Usage = "query and display information from Untappd APIv4"
	app.Version = "0.0.1"
	app.Authors = []*cli.Author{
		&cli.Author{
			Name:  "Matt Layher",
			Email: "mdlayher@gmail.com",
		},
	}

	// Add global flags for Untappd API client ID, client secret, and
	// authenticated access token
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "client_id",
			Usage:   "client ID parameter for Untappd APIv4",
			EnvVars: []string{"UNTAPPD_ID"},
		},
		&cli.StringFlag{
			Name:    "client_secret",
			Usage:   "client secret parameter for Untappd APIv4",
			EnvVars: []string{"UNTAPPD_SECRET"},
		},
		&cli.StringFlag{
			Name:    "access_token",
			Usage:   "authenticated access token for Untappd APIv4",
			EnvVars: []string{"UNTAPPD_TOKEN"},
		},
	}

	// Frequently used flags for paging and sorting results, with their
	// default Untappd API values
	offsetFlag := &cli.IntFlag{
		Name:  "offset",
		Value: 0,
		Usage: "starting offset for API query results",
	}
	limitFlag := &cli.IntFlag{
		Name:  "limit",
		Value: 25,
		Usage: "maximum number of API query results",
	}
	sortFlag := &cli.StringFlag{
		Name:  "sort",
		Value: string(untappd.SortDate),
		Usage: fmt.Sprintf("sort type for API query results (options: %s)", untappd.Sorts()),
	}

	// Flags used to specify minimum and maximum checkin IDs
	minIDFlag := &cli.IntFlag{
		Name:  "min_id",
		Value: 0,
		Usage: "minimum checkin ID for API query results",
	}
	maxIDFlag := &cli.IntFlag{
		Name:  "max_id",
		Value: math.MaxInt32,
		Usage: "maximum checkin ID for API query results",
	}

	// Add commands mirroring available untappd.Client services
	app.Commands = []*cli.Command{
		authCommand(limitFlag, minIDFlag, maxIDFlag),
		beerCommand(offsetFlag, limitFlag, sortFlag, minIDFlag, maxIDFlag),
		breweryCommand(offsetFlag, limitFlag, minIDFlag, maxIDFlag),
		localCommand(limitFlag, minIDFlag, maxIDFlag),
		userCommand(offsetFlag, limitFlag, sortFlag, minIDFlag, maxIDFlag),
		venueCommand(limitFlag, minIDFlag, maxIDFlag),
	}

	// Print all log output to stderr, so stdout only contains Untappd data
	log.SetFlags(0)
	log.SetOutput(os.Stderr)
	log.SetPrefix(appName + "> ")

	app.Run(os.Args)
}

// untappdClient creates an initialized *untappd.Client using either the
// access token, or client ID and secret from global CLI context.
func untappdClient(ctx *cli.Context) *untappd.Client {
	var c *untappd.Client
	var err error

	// Always prefer authenticated access token, if available
	token := ctx.String("access_token")
	if token != "" {
		c, err = untappd.NewAuthenticatedClient(token, nil)
	} else {
		c, err = untappd.NewClient(
			ctx.String("client_id"),
			ctx.String("client_secret"),
			nil,
		)
	}
	if err != nil {
		log.Fatal(err)
	}

	return c
}

// printRateLimit is a helper method which displays the remaining rate limit
// header for each HTTP request.
func printRateLimit(res *http.Response) {
	const header = "X-Ratelimit-Remaining"
	if v := res.Header.Get(header); v != "" {
		log.Printf("%s: %s", header, v)
	}
}

// mustStringArg is a helper method which checks for a string argument in the
// CLI context, and prints a help message if it is not found.
func mustStringArg(ctx *cli.Context, name string) string {
	a := ctx.Args().First()
	if a == "" {
		log.Fatalf("missing argument: %s", name)
	}

	return a
}

// offsetLimitSort retrieves a triple of offset, limit, and sort parameters
// from CLI context, as accepted by the Untappd API.
func offsetLimitSort(ctx *cli.Context) (int, int, untappd.Sort) {
	offset, limit := ctx.Int("offset"), ctx.Int("limit")

	// If no sort found, ignore sanity checks
	sort := ctx.String("sort")
	if sort == "" {
		return offset, limit, untappd.Sort("")
	}

	// Ensure sort type is valid
	for _, s := range untappd.Sorts() {
		// Return on valid sort
		if sort == string(s) {
			return offset, limit, s
		}
	}

	// Die on invalid sort, and show options
	log.Fatalf("invalid sort type %q (options: %s)", sort, untappd.Sorts())
	return offset, limit, untappd.Sort("")
}

// checkAtoiError reduces error-checking code duplication for functions
// which require valid integer IDs.
func checkAtoiError(err error) {
	if err == nil {
		return
	}

	nErr, ok := err.(*strconv.NumError)
	if !ok {
		log.Fatal(err)
	}

	if nErr.Err == strconv.ErrSyntax {
		log.Fatal("invalid integer ID")
	}

	log.Fatal(err)
}
