//
package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"path"
)

// While we could have incorporated a full featured CLI library, such as github.com/urfave/cli or
// github.com/spf13/cobra, the added dependency is likely not warranted at this point. A full featured library would,
// for example, simplify the definition of global and per sub-command flags, or allow defining short and long versions
// for each option... Since the standard flag package does not support these, we resort to using positional arguments
// for sub-commands.
//
// Quoting https://go-proverbs.github.io: "A little copying is better than a little dependency" ;-)
//

func main() {
	// global flags shared by all sub-commands
	//
	// see https://github.com/gameontext/gameon-room-go#get-game-on-id-and-shared-secret for directions
	// on generating and retrieving the shared secret and id values
	var verbose bool        // CLI verbosity, boolean flags (default: false)
	var sharedSecret string // GameOn! shared secret, required. Set in command line or GAMEON_SECRET env-var
	var identity string     // GameOn! identity, required. Set in command line or GAMEON_ID env-var
	var server urlFlag      // GameOn! server URL

	addCommonFlags(&verbose, &sharedSecret, &identity, &server)

	//-- Sub-commands
	subcommands := []subcommand{deleteRoomSubcommand(), listRoomsSubcommand(), registerRoomSubcommand()}

	flag.Parse()
	err := validateRequiredCommonFlags(verbose, sharedSecret, identity, server.Url())

	if err != nil {
		usage(err.Error(), subcommands)
		os.Exit(1)
	}

	tail := flag.Args() // left-over, unprocessed, positional args

	// verify that a sub-command has been provided (tail[0] is the sub-command)
	if len(tail) < 1 {
		usage("A sub-command is required", subcommands)
		os.Exit(1)
	}

	for _, sc := range subcommands {
		if sc.Keyword() == tail[0] {
			err := sc.Parse(os.Args[1:]) // tail[1:] will be all arguments following sub-command
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			// create the GameOn client
			err = sc.Process(nil)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			return // all done, command completed successfully
		}
	}

	// if we got this far, tail[0] did not match any sub-command keyword
	fmt.Println(os.Args[0], "Unknown sub-command:", tail[0])
	fmt.Print("Expecting one of ")
	for _, sc := range subcommands {
		fmt.Print(sc.Keyword(), " ")
	}
	fmt.Println()
	os.Exit(1)
}

func addCommonFlags(verbose *bool, secret *string, identity *string, serverUrl *urlFlag) {
	flag.StringVar(identity, "id", os.Getenv("GAMEON_ID"),
		"GameOn! identity (required, default: $GAMEON_ID environment variable)")
	flag.StringVar(secret, "secret", os.Getenv("GAMEON_SECRET"),
		"GameOn! shared secret (required, default: $GAMEON_SECRET environment variable)")
	flag.BoolVar(verbose, "d", false, "Enable debug (optional, default: false)")
	flag.Var(serverUrl, "g", "GameOn! server URL (required)")
}

// validate that (required) flags have been successfully parsed
func validateRequiredCommonFlags(verbose bool, secret string, identity string, gos *url.URL) error {
	if identity == "" {
		return fmt.Errorf("identity not provided in command line or $GAMEON_ID environment variable")
	} else if secret == "" {
		return fmt.Errorf("shared secret not provided in command line or $GAMEON_SECRET environment variable")
	} else if gos == nil {
		return fmt.Errorf("no valid URL for GameOn! server")
	}

	if verbose {
		fmt.Printf("id:%s secret:%s server:%s verbose:%t\r\n", identity, secret, gos.String(), verbose)
	}

	return nil
}

func usage(message string, subcommands []subcommand) {
	fmt.Println("Error:", message, "\r\n")
	fmt.Println("Usage:\r\n\t", path.Base(os.Args[0]), "<global flags>", "<sub-command> [sub-command flags]")
	fmt.Println("\r\nGlobal flags:")
	flag.PrintDefaults()
	fmt.Println("\r\nSub-commands:")
	for _, sc := range subcommands {
		fmt.Println("  ", sc.Usage())
	}
}
