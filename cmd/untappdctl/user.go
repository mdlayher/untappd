package main

import (
	"log"

	"github.com/codegangsta/cli"
	"github.com/mdlayher/untappd"
)

// userCommand allows access to untappd.Client.User methods, such as user
// information, checked in beers, friends, badges, and wish list.
func userCommand(offsetFlag cli.IntFlag, limitFlag cli.IntFlag, sortFlag cli.StringFlag) cli.Command {
	return cli.Command{
		Name:    "user",
		Aliases: []string{"u"},
		Usage:   "query for user information, by username",
		Subcommands: []cli.Command{
			userBadgesCommand(offsetFlag, limitFlag),
			userBeersCommand(offsetFlag, limitFlag, sortFlag),
			userFriendsCommand(offsetFlag, limitFlag),
			userInfoCommand(),
			userWishListCommand(offsetFlag, limitFlag, sortFlag),
		},
	}
}

// userBadgesCommand allows access to the untappd.Client.User.Badges method, which
// can query for information about a user's badges, by username.
func userBadgesCommand(offsetFlag cli.IntFlag, limitFlag cli.IntFlag) cli.Command {
	return cli.Command{
		Name:    "badges",
		Aliases: []string{"ba"},
		Usage:   "query for badges a user has earned, by username",
		Flags: []cli.Flag{
			offsetFlag,
			limitFlag,
		},

		Action: func(ctx *cli.Context) {
			offset, limit, _ := offsetLimitSort(ctx)

			// Query for user's earned badges by username, e.g.
			// "untappdctl user badges mdlayher"
			c := untappdClient(ctx)
			badges, res, err := c.User.BadgesOffsetLimit(
				mustStringArg(ctx, "username"),
				offset,
				limit,
			)
			printRateLimit(res)
			if err != nil {
				log.Fatal(err)
			}

			// Print out badges in human-readable format
			printBadges(badges)
		},
	}
}

// userBeersCommand allows access to the untappd.Client.User.Beers method, which
// can query for information about a user's checked in beers, by username.
func userBeersCommand(offsetFlag cli.IntFlag, limitFlag cli.IntFlag, sortFlag cli.StringFlag) cli.Command {
	return cli.Command{
		Name:    "beers",
		Aliases: []string{"be"},
		Usage:   "query for beers a user has checked in, by username",
		Flags: []cli.Flag{
			offsetFlag,
			limitFlag,
			sortFlag,
		},

		Action: func(ctx *cli.Context) {
			offset, limit, sort := offsetLimitSort(ctx)

			// Query for user's checked in beers by username, e.g.
			// "untappdctl user beers mdlayher"
			c := untappdClient(ctx)
			beers, res, err := c.User.BeersOffsetLimitSort(
				mustStringArg(ctx, "username"),
				offset,
				limit,
				untappd.Sort(sort),
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

// userFriendsCommand allows access to the untappd.Client.User.Friends method, which
// can query for information about a user's friends, by username.
func userFriendsCommand(offsetFlag cli.IntFlag, limitFlag cli.IntFlag) cli.Command {
	return cli.Command{
		Name:    "friends",
		Aliases: []string{"f"},
		Usage:   "query for a user's friends, by username",
		Flags: []cli.Flag{
			offsetFlag,
			limitFlag,
		},

		Action: func(ctx *cli.Context) {
			offset, limit, _ := offsetLimitSort(ctx)

			// Query for user's friends by username, e.g.
			// "untappdctl user friends mdlayher"
			c := untappdClient(ctx)
			friends, res, err := c.User.FriendsOffsetLimit(
				mustStringArg(ctx, "username"),
				offset,
				limit,
			)
			printRateLimit(res)
			if err != nil {
				log.Fatal(err)
			}

			// Print out users in human-readable format
			printUsers(friends, false)
		},
	}
}

// userInfoCommand allows access to the untappd.Client.User.Info method, which
// can query for information about a user, by username.
func userInfoCommand() cli.Command {
	return cli.Command{
		Name:    "info",
		Aliases: []string{"i"},
		Usage:   "query for user information, such as ID, real name, etc. by username",

		Action: func(ctx *cli.Context) {
			// Query for user by username, e.g. "untappdctl user info mdlayher"
			c := untappdClient(ctx)
			user, res, err := c.User.Info(mustStringArg(ctx, "username"), false)
			printRateLimit(res)
			if err != nil {
				log.Fatal(err)
			}

			// Print out user in human-readable format
			printUsers([]*untappd.User{user}, true)
		},
	}
}

// userWishListCommand allows access to the untappd.Client.User.WishList method,
// which can query for information about a user's wish list beers, by username.
func userWishListCommand(offsetFlag cli.IntFlag, limitFlag cli.IntFlag, sortFlag cli.StringFlag) cli.Command {
	return cli.Command{
		Name:    "wishlist",
		Aliases: []string{"w"},
		Usage:   "query for beers a user has on their wishlist, by username",
		Flags: []cli.Flag{
			offsetFlag,
			limitFlag,
			sortFlag,
		},

		Action: func(ctx *cli.Context) {
			offset, limit, sort := offsetLimitSort(ctx)

			// Query for user wishlist beers by username,
			// e.g. "untappdctl user wishlist mdlayher"
			c := untappdClient(ctx)
			beers, res, err := c.User.WishListOffsetLimitSort(
				mustStringArg(ctx, "username"),
				offset,
				limit,
				untappd.Sort(sort),
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
