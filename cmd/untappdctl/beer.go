package main

import (
	"log"
	"strconv"

	"github.com/codegangsta/cli"
	"github.com/mdlayher/untappd"
)

// beerCommand allows access to untappd.Client.Beer methods, such as beer
// information by ID, and query by search term.
func beerCommand(offsetFlag cli.IntFlag, limitFlag cli.IntFlag, sortFlag cli.StringFlag) cli.Command {
	return cli.Command{
		Name:    "beer",
		Aliases: []string{"be"},
		Usage:   "query for beer information, by beer ID or name",
		Subcommands: []cli.Command{
			beerInfoCommand(),
			beerSearchCommand(offsetFlag, limitFlag, sortFlag),
		},
	}
}

// beerInfoCommand allows access to the untappd.Client.Beer.Info method, which
// can query for information about a beer, by ID.
func beerInfoCommand() cli.Command {
	return cli.Command{
		Name:    "info",
		Aliases: []string{"i"},
		Usage:   "query for beer information, by ID",

		Action: func(ctx *cli.Context) {
			// Check for valid integer ID
			id, err := strconv.Atoi(mustStringArg(ctx, "beer ID"))
			checkAtoiError(err)

			// Query for beer by ID, e.g. "untappdctl beer info 1"
			c := untappdClient(ctx)
			beer, res, err := c.Beer.Info(id, false)
			printRateLimit(res)
			if err != nil {
				log.Fatal(err)
			}

			// Print out beer in human-readable format
			printBeers([]*untappd.Beer{beer})
		},
	}
}

// beerSearchCommand allows access to the untappd.Client.Beer.Search method, which
// can search for information about beers, by search term.
func beerSearchCommand(offsetFlag cli.IntFlag, limitFlag cli.IntFlag, sortFlag cli.StringFlag) cli.Command {
	return cli.Command{
		Name:    "search",
		Aliases: []string{"s"},
		Usage:   "search for beers, by brewery and/or beer name",
		Flags: []cli.Flag{
			offsetFlag,
			limitFlag,
		},

		Action: func(ctx *cli.Context) {
			offset, limit, sort := offsetLimitSort(ctx)

			// Query for beer's earned beers by name, e.g.
			// "untappdctl beer search oberon"
			c := untappdClient(ctx)
			beers, res, err := c.Beer.SearchOffsetLimitSort(
				mustStringArg(ctx, "beer name (optionally, with brewery name)"),
				offset,
				limit,
				sort,
			)
			printRateLimit(res)
			if err != nil {
				log.Fatal(err)
			}

			// Print out beers in human-readable format
			printBeers(beers)
		},
	}
}
