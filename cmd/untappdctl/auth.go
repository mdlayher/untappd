package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/codegangsta/cli"
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

func authLoginCommand() cli.Command {
	return cli.Command{
		Name:    "login",
		Aliases: []string{"l"},
		Usage:   "authenticate using OAuth to Untappd APIv4",

		Action: func(ctx *cli.Context) {
			// 8338 looks kinda like "BEER", right?
			const host = ":8338"
			const endpoint = "/auth"

			// Set up redirect URL, which will use our HTTP server
			redirectURL := fmt.Sprintf("http://localhost%s%s", host, endpoint)

			// Start listening for TCP connections
			l, err := net.Listen("tcp", host)
			if err != nil {
				log.Fatal(err)
			}

			// Block until authentication flow is complete
			doneC := make(chan struct{}, 0)

			// Handle HTTP connections for OAuth authentication
			mux := http.NewServeMux()
			mux.HandleFunc(endpoint, func(w http.ResponseWriter, r *http.Request) {
				// Set up HTTP request to Untappd APIv4, with intermediary token
				// passed with redirect URL to this HTTP server
				u := fmt.Sprintf(
					"https://untappd.com/oauth/authorize/?client_id=%s&client_secret=%s&response_type=code&redirect_url=%s&code=%s",
					ctx.GlobalString("client_id"),
					ctx.GlobalString("client_secret"),
					redirectURL,
					r.URL.Query().Get("code"),
				)

				// Perform HTTP GET request
				res, err := http.Get(u)
				if err != nil {
					log.Fatal(err)
				}

				// Temporary struct for JSON body
				var v struct {
					Response struct {
						AccessToken string `json:"access_token"`
					} `json:"response"`
				}

				// Decode JSON body to retrieve token
				if err := json.NewDecoder(res.Body).Decode(&v); err != nil {
					log.Fatal(err)
				}

				// Provide token to user
				t := v.Response.AccessToken
				log.Println("token:", t)
				fmt.Fprintln(w, "token:", t)

				// Close HTTP listener so no further connections can be made,
				// and stop blocking untappdctl from exiting
				_ = l.Close()
				close(doneC)
			})

			// Start HTTP server in background, using our custom authentication handler
			go func() {
				if err := (&http.Server{
					Handler: mux,
				}).Serve(l); err != nil {
					// Ignore this error on shutdown
					if !strings.Contains(err.Error(), "use of closed network connection") {
						log.Println(err)
					}
				}
			}()

			// Provide link for user to open to start authentication flow
			log.Printf(
				"https://untappd.com/oauth/authenticate/?client_id=%s&response_type=code&redirect_url=%s",
				ctx.GlobalString("client_id"),
				redirectURL,
			)

			// Block until authentication flow is complete
			<-doneC
		},
	}
}
