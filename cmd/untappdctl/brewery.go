package main

import (
	"log"
	"strconv"

	"github.com/codegangsta/cli"
	"github.com/mdlayher/untappd"
)

// breweryCommand allows access to untappd.Client.Brewery methods, such as brewery
// information by ID, and query by search term.
func breweryCommand(offsetFlag cli.IntFlag, limitFlag cli.IntFlag) cli.Command {
	return cli.Command{
		Name:    "brewery",
		Aliases: []string{"br"},
		Usage:   "query for brewery information, by brewery ID or name",
		Subcommands: []cli.Command{
			breweryInfoCommand(),
			brewerySearchCommand(offsetFlag, limitFlag),
		},
	}
}

// breweryInfoCommand allows access to the untappd.Client.Brewery.Info method, which
// can query for information about a brewery, by ID.
func breweryInfoCommand() cli.Command {
	return cli.Command{
		Name:    "info",
		Aliases: []string{"i"},
		Usage:   "query for brewery information, by ID",

		Action: func(ctx *cli.Context) {
			// Check for valid integer ID
			id, err := strconv.Atoi(mustStringArg(ctx, "brewery ID"))
			if err != nil {
				nErr, ok := err.(*strconv.NumError)
				if !ok {
					log.Fatal(err)
				}

				if nErr.Err == strconv.ErrSyntax {
					log.Fatal("invalid integer ID")
				}

				log.Fatal(err)
			}

			// Query for brewery by ID, e.g. "untappdctl brewery info 1"
			c := untappdClient(ctx)
			brewery, res, err := c.Brewery.Info(id, false)
			printRateLimit(res)
			if err != nil {
				log.Fatal(err)
			}

			// Print out brewery in human-readable format
			printBreweries([]*untappd.Brewery{brewery})
		},
	}
}

// brewerySearchCommand allows access to the untappd.Client.Brewery.Search method, which
// can search for information about breweries, by search term.
func brewerySearchCommand(offsetFlag cli.IntFlag, limitFlag cli.IntFlag) cli.Command {
	return cli.Command{
		Name:    "search",
		Aliases: []string{"s"},
		Usage:   "search for breweries, by brewery name",
		Flags: []cli.Flag{
			offsetFlag,
			limitFlag,
		},

		Action: func(ctx *cli.Context) {
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
		},
	}
}
