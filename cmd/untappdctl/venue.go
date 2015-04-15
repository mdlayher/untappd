package main

import (
	"log"
	"strconv"

	"github.com/codegangsta/cli"
	"github.com/mdlayher/untappd"
)

// venueCommand allows access to untappd.Client.Venue methods, such as venue
// information by ID.
func venueCommand() cli.Command {
	return cli.Command{
		Name:    "venue",
		Aliases: []string{"v"},
		Usage:   "query for venue information, by venue ID",
		Subcommands: []cli.Command{
			venueInfoCommand(),
		},
	}
}

// venueInfoCommand allows access to the untappd.Client.Venue.Info method, which
// can query for information about a venue, by ID.
func venueInfoCommand() cli.Command {
	return cli.Command{
		Name:    "info",
		Aliases: []string{"i"},
		Usage:   "query for venue information, by ID",

		Action: func(ctx *cli.Context) {
			// Check for valid integer ID
			id, err := strconv.Atoi(mustStringArg(ctx, "venue ID"))
			checkAtoiError(err)

			// Query for venue by ID, e.g. "untappdctl venue info 1"
			c := untappdClient(ctx)
			venue, res, err := c.Venue.Info(id, false)
			printRateLimit(res)
			if err != nil {
				log.Fatal(err)
			}

			// Print out venue in human-readable format
			printVenues([]*untappd.Venue{venue})
		},
	}
}
