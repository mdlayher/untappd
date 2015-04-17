package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/mdlayher/untappd"
)

// localCommand allows access to untappd.Client.Local methods, such as local
// checkins by latitude and longitude.
func localCommand(limitFlag cli.IntFlag, minIDFlag cli.IntFlag, maxIDFlag cli.IntFlag) cli.Command {
	return cli.Command{
		Name:    "local",
		Aliases: []string{"l"},
		Usage:   "query for local area checkins, by latitude and longitude",
		Subcommands: []cli.Command{
			localCheckinsCommand(limitFlag, minIDFlag, maxIDFlag),
		},
	}
}

// localCheckinsCommand allows access to the untappd.Client.Local.Checkins method, which
// can query for information about recent checkins for a local area, by latitude, longitude,
// and several other parameters.
func localCheckinsCommand(limitFlag cli.IntFlag, minIDFlag cli.IntFlag, maxIDFlag cli.IntFlag) cli.Command {
	return cli.Command{
		Name:    "checkins",
		Aliases: []string{"c"},
		Usage:   "query for recent checkins for a local area, by latitude and longitude",
		Flags: []cli.Flag{
			limitFlag,
			minIDFlag,
			maxIDFlag,
			cli.IntFlag{
				Name:  "radius",
				Value: 25,
				Usage: "checkin radius around latitude,longitude pair",
			},
			cli.StringFlag{
				Name:  "unit",
				Value: string(untappd.DistanceMiles),
				Usage: fmt.Sprintf("units for radius, either %q or %q", untappd.DistanceMiles, untappd.DistanceKilometers),
			},
		},

		Action: func(ctx *cli.Context) {
			// Check for valid latitude and longitude pair
			pair := strings.Split(mustStringArg(ctx, "latitude,longitude pair"), ",")
			if len(pair) != 2 {
				log.Fatal("pair must in form: latitude,longitude")
			}

			// Basic semantic check for valid floating point numbers
			lat, err := strconv.ParseFloat(pair[0], 64)
			lng, err2 := strconv.ParseFloat(pair[1], 64)
			if err != nil || err2 != nil {
				log.Fatal("latitude,longitude pair must be floating point values")
			}

			// Validate units
			unit := untappd.Distance(ctx.String("unit"))
			if unit != untappd.DistanceMiles && unit != untappd.DistanceKilometers {
				log.Fatalf("unit must be %q or %q", untappd.DistanceMiles, untappd.DistanceKilometers)
			}

			// Query for local's checkins by local area with latitude,longitude
			// pair, e.g.
			// "untappdctl local checkins 42.291,-85.587"
			c := untappdClient(ctx)
			checkins, res, err := c.Local.CheckinsMinMaxIDLimitRadius(
				lat,
				lng,
				ctx.Int("min_id"),
				ctx.Int("max_id"),
				ctx.Int("limit"),
				ctx.Int("radius"),
				unit,
			)
			printRateLimit(res)
			if err != nil {
				log.Fatal(err)
			}

			// Print out checkins in human-readable format
			printCheckins(checkins)
		},
	}
}
