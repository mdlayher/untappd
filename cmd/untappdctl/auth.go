package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/mdlayher/untappd"
)

// authCommand allows a user to easily authenticate to the Untappd APIv4 using
// their client ID and client secret.  A temporary HTTP server is started, and
// authentication flow is handled automatically once the user clicks the initial
// URL.
func authCommand(limitFlag cli.IntFlag, minIDFlag cli.IntFlag, maxIDFlag cli.IntFlag) cli.Command {
	return cli.Command{
		Name:    "auth",
		Aliases: []string{"a"},
		Usage:   "access authenticated Untappd APIv4 methods",
		Subcommands: []cli.Command{
			authCheckinsCommand(limitFlag, minIDFlag, maxIDFlag),
			authLoginCommand(),
		},
	}
}

// authCheckinsCommand allows access to the untappd.Client.Beer.Checkins method, which
// can query for information about recent checkins for a beer, by ID.
func authCheckinsCommand(limitFlag cli.IntFlag, minIDFlag cli.IntFlag, maxIDFlag cli.IntFlag) cli.Command {
	return cli.Command{
		Name:    "checkins",
		Aliases: []string{"c"},
		Usage:   "[auth] query for recent checkins from friends",
		Flags: []cli.Flag{
			limitFlag,
			minIDFlag,
			maxIDFlag,
		},

		Action: func(ctx *cli.Context) {
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
		},
	}
}

// authLoginCommand performs the OAuth Authentication process required to retrieve
// an Access Token for the Untappd APIv4.
func authLoginCommand() cli.Command {
	return cli.Command{
		Name:    "login",
		Aliases: []string{"l"},
		Usage:   "authenticate using OAuth to Untappd APIv4",

		Action: func(ctx *cli.Context) {
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
			doneC := make(chan struct{}, 0)

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
				ctx.GlobalString("client_id"),
				ctx.GlobalString("client_secret"),
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
		},
	}
}
