package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/codegangsta/cli"
	"github.com/mdlayher/untappd"
)

// authCommand allows a user to easily authenticate to the Untappd APIv4, and
// perform actions which require authentication, such as checking in beers.
func authCommand(limitFlag, minIDFlag, maxIDFlag *cli.IntFlag) *cli.Command {
	return &cli.Command{
		Name:    "auth",
		Aliases: []string{"a"},
		Usage:   "access authenticated Untappd APIv4 methods",
		Subcommands: []*cli.Command{
			authCheckinCommand(),
			authCheckinsCommand(limitFlag, minIDFlag, maxIDFlag),
			authLoginCommand(),
		},
	}
}

// authCheckinCommand allows access to the untappd.Client.Beer.Checkin method, which
// can check in a beer, by ID.
func authCheckinCommand() *cli.Command {
	return &cli.Command{
		Name:  "checkin",
		Usage: "[auth] check-in a beer, by ID",
		Flags: []cli.Flag{
			&cli.Float64Flag{
				Name:  "rating",
				Usage: "optional rating, 0.5-5.0, for this checkin",
			},
			&cli.StringFlag{
				Name:  "comment",
				Usage: "optional comment for this checkin",
			},
		},

		Action: func(ctx *cli.Context) error {
			// Check for valid integer ID
			id, err := strconv.Atoi(mustStringArg(ctx, "beer ID"))
			checkAtoiError(err)

			// Use system's timezone and offset for request,
			// dividing to get a single digit offset
			// Thanks: https://github.com/cmar/untappd/blob/master/lib/untappd/checkin.rb#L50
			timezone, offset := time.Now().Zone()
			offset = offset / 60 / 60

			// Attempt to perform checkin
			c := untappdClient(ctx)
			checkin, res, err := c.Auth.Checkin(untappd.CheckinRequest{
				BeerID:    id,
				GMTOffset: offset,
				TimeZone:  timezone,
				Comment:   ctx.String("comment"),
				Rating:    ctx.Float64("rating"),
			})
			printRateLimit(res)
			if err != nil {
				log.Fatal(err)
			}

			// Print out checkin in human-readable format
			printCheckins([]*untappd.Checkin{checkin})
			return nil
		},
	}
}

// authCheckinsCommand allows access to the untappd.Client.Beer.Checkins method, which
// can query for information about recent checkins for a beer, by ID.
func authCheckinsCommand(limitFlag, minIDFlag, maxIDFlag *cli.IntFlag) *cli.Command {
	return &cli.Command{
		Name:  "checkins",
		Usage: "[auth] query for recent checkins from friends",
		Flags: []cli.Flag{
			limitFlag,
			minIDFlag,
			maxIDFlag,
		},

		Action: func(ctx *cli.Context) error {
			// Query for checkins by beername, e.g.
			// "untappdctl beer checkins mdlayher"
			c := untappdClient(ctx)
			checkins, res, err := c.Auth.CheckinsMinMaxIDLimit(
				ctx.Int("min_id"),
				ctx.Int("max_id"),
				ctx.Int("limit"),
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

// authLoginCommand performs the OAuth Authentication process required to retrieve
// an Access Token for the Untappd APIv4.
func authLoginCommand() *cli.Command {
	return &cli.Command{
		Name:    "login",
		Aliases: []string{"l"},
		Usage:   "authenticate using OAuth to Untappd APIv4",

		Action: func(ctx *cli.Context) error {
			// 8338 looks kinda like "BEER", right?
			const host = ":8338"

			// Set up redirect URL, which will use our HTTP server
			redirectURL := fmt.Sprintf("http://localhost%s", host)

			// Start listening for TCP connections
			l, err := net.Listen("tcp", host)
			if err != nil {
				log.Fatal(err)
			}

			// Wait for a single token to arrive, then cancel listener
			doneC := make(chan struct{})

			// Handle response token by providing to to both HTTP response
			// and terminal output
			tokenFn := func(token string, w http.ResponseWriter, r *http.Request) {
				// Print token in terminal and to HTTP response body
				log.Println("token:", token)
				if _, err := w.Write([]byte(token)); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				// Close HTTP listener to prevent further requests
				_ = l.Close()
				close(doneC)
			}

			// Set up http.Handler which allows easy OAuth authentication
			// with Untappd APIv4
			h, clientURL, err := untappd.NewAuthHandler(
				ctx.String("client_id"),
				ctx.String("client_secret"),
				redirectURL,
				tokenFn,
				nil,
			)
			if err != nil {
				log.Fatal(err)
			}

			// Start HTTP server in background, using our custom authentication handler
			go func() {
				if err := (&http.Server{
					Handler: h,
				}).Serve(l); err != nil {
					// Ignore this error on shutdown
					if !strings.Contains(err.Error(), "use of closed network connection") {
						log.Println(err)
					}
				}
			}()

			// Provide link for user to open to start authentication flow
			log.Println(clientURL.String())

			// Block until one authentication completes
			<-doneC
			return nil
		},
	}
}
