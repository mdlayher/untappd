package main

import (
	"log"
	"strconv"

	"github.com/codegangsta/cli"
	"github.com/mdlayher/untappd"
)

// venueCommand allows access to untappd.Client.Venue methods, such as venue
// information by ID.
func venueCommand(limitFlag, minIDFlag, maxIDFlag *cli.IntFlag) *cli.Command {
	return &cli.Command{
		Name:    "venue",
		Aliases: []string{"v"},
		Usage:   "query for venue information, by venue ID",
		Subcommands: []*cli.Command{
			venueCheckinsCommand(limitFlag, minIDFlag, maxIDFlag),
			venueInfoCommand(),
		},
	}
}

// venueCheckinsCommand allows access to the untappd.Client.Venue.Checkins method, which
// can query for information about recent checkins at a specified venue, by ID.
func venueCheckinsCommand(limitFlag, minIDFlag, maxIDFlag *cli.IntFlag) *cli.Command {
	return &cli.Command{
		Name:    "checkins",
		Aliases: []string{"c"},
		Usage:   "query for recent checkins at a specified venue, by ID",
		Flags: []cli.Flag{
			limitFlag,
			minIDFlag,
			maxIDFlag,
		},

		Action: func(ctx *cli.Context) error {
			// Check for valid integer ID
			id, err := strconv.Atoi(mustStringArg(ctx, "venue ID"))
			checkAtoiError(err)

			minID, maxID, limit := ctx.Int("min_id"), ctx.Int("max_id"), ctx.Int("limit")

			// Query for venue's checkins by venue ID, e.g.
			// "untappdctl venue checkins 1"
			c := untappdClient(ctx)
			checkins, res, err := c.Venue.CheckinsMinMaxIDLimit(
				id,
				minID,
				maxID,
				limit,
			)
			printRateLimit(res)
			if err != nil {
				log.Fatal(err)
			}

			// Print out checkins in human-readable format
			printCheckins(checkins)
			return nil
		},
	}
}

// venueInfoCommand allows access to the untappd.Client.Venue.Info method, which
// can query for information about a venue, by ID.
func venueInfoCommand() *cli.Command {
	return &cli.Command{
		Name:    "info",
		Aliases: []string{"i"},
		Usage:   "query for venue information, by ID",

		Action: func(ctx *cli.Context) error {
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
			return nil
		},
	}
}
