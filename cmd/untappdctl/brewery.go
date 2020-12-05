package main

import (
	"log"
	"strconv"

	"github.com/urfave/cli/v2"
	"github.com/mdlayher/untappd"
)

// breweryCommand allows access to untappd.Client.Brewery methods, such as brewery
// information by ID, and query by search term.
func breweryCommand(offsetFlag, limitFlag, minIDFlag, maxIDFlag *cli.IntFlag) *cli.Command {
	return &cli.Command{
		Name:    "brewery",
		Aliases: []string{"br"},
		Usage:   "query for brewery information, by brewery ID or name",
		Subcommands: []*cli.Command{
			breweryCheckinsCommand(limitFlag, minIDFlag, maxIDFlag),
			breweryInfoCommand(),
			brewerySearchCommand(offsetFlag, limitFlag),
		},
	}
}

// breweryCheckinsCommand allows access to the untappd.Client.Brewery.Checkins method, which
// can query for information about recent checkins for beers made by a brewery, by ID.
func breweryCheckinsCommand(limitFlag, minIDFlag, maxIDFlag *cli.IntFlag) *cli.Command {
	return &cli.Command{
		Name:    "checkins",
		Aliases: []string{"c"},
		Usage:   "query for recent checkins for beers from a specified brewery, by ID",
		Flags: []cli.Flag{
			limitFlag,
			minIDFlag,
			maxIDFlag,
		},

		Action: func(ctx *cli.Context) error {
			// Check for valid integer ID
			id, err := strconv.Atoi(mustStringArg(ctx, "brewery ID"))
			checkAtoiError(err)

			minID, maxID, limit := ctx.Int("min_id"), ctx.Int("max_id"), ctx.Int("limit")

			// Query for brewery's checkins by brewery ID, e.g.
			// "untappdctl brewery checkins 1"
			c := untappdClient(ctx)
			checkins, res, err := c.Brewery.CheckinsMinMaxIDLimit(
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

// breweryInfoCommand allows access to the untappd.Client.Brewery.Info method, which
// can query for information about a brewery, by ID.
func breweryInfoCommand() *cli.Command {
	return &cli.Command{
		Name:    "info",
		Aliases: []string{"i"},
		Usage:   "query for brewery information, by ID",

		Action: func(ctx *cli.Context) error {
			// Check for valid integer ID
			id, err := strconv.Atoi(mustStringArg(ctx, "brewery ID"))
			checkAtoiError(err)

			// Query for brewery by ID, e.g. "untappdctl brewery info 1"
			c := untappdClient(ctx)
			brewery, res, err := c.Brewery.Info(id, false)
			printRateLimit(res)
			if err != nil {
				log.Fatal(err)
			}

			// Print out brewery in human-readable format
			printBreweries([]*untappd.Brewery{brewery})
			return nil
		},
	}
}

// brewerySearchCommand allows access to the untappd.Client.Brewery.Search method, which
// can search for information about breweries, by search term.
func brewerySearchCommand(offsetFlag, limitFlag *cli.IntFlag) *cli.Command {
	return &cli.Command{
		Name:    "search",
		Aliases: []string{"s"},
		Usage:   "search for breweries, by brewery name",
		Flags: []cli.Flag{
			offsetFlag,
			limitFlag,
		},

		Action: func(ctx *cli.Context) error {
			offset, limit, _ := offsetLimitSort(ctx)

			// Query for brewery's earned breweries by name, e.g.
			// "untappdctl brewery search oberon"
			c := untappdClient(ctx)
			breweries, res, err := c.Brewery.SearchOffsetLimit(
				mustStringArg(ctx, "brewery name"),
				offset,
				limit,
			)
			printRateLimit(res)
			if err != nil {
				log.Fatal(err)
			}

			// Print out breweries in human-readable format
			printBreweries(breweries)
			return nil
		},
	}
}
