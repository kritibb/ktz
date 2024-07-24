// Package main is the entry point for the timezone application.
// It retrieves the timezone for a given city name
// and displays the current local time for that timezone.
package main

import (
	"flag"
	"fmt"
	"github.com/kritibb/ktz/cmd"
	"os"
	"strings"
)

// main is the entry point of the application.
// It expects command-line arguments in the form of ktz lookup [options] [city/country],
// retrieves the corresponding timezone, and displays the current local time in that timezone.
func main() {

	//define subcommand `lookup` and its flags
	lookupCmd := flag.NewFlagSet("lookup", flag.ExitOnError)
	lookupTz := lookupCmd.String("tz", "", "timezone names like `Asia/Kathmandu`")
	lookupZ := lookupCmd.String("z", "", "zone names like `pst`")

	// Check if subcommands like "lookup" is provided
	if len(os.Args) < 2 {
		printError("", "\n Error: Incomplete command")
		return
	}

	switch os.Args[1] {
	case "lookup":
		lookupCmd.Parse(os.Args[2:])
		if *lookupTz == "" && *lookupZ == "" && len(lookupCmd.Args()) == 0 {
			printError("lookup", "\n Error: Incomplete command")
			return
		}
		if (*lookupTz != "" && *lookupZ != "") || (*lookupZ != "" && len(lookupCmd.Args()) != 0) ||
			(*lookupTz != "" && len(lookupCmd.Args()) != 0) {
			printError("lookup", "\n Error: Use only one flag [-tz or -z] or [city/country]")
			return
		}
		if *lookupTz != "" {
			cmd.DisplayTime("",*lookupTz)
		} else if *lookupZ != "" {
			// fmt.Println("Display time based on z:", *lookupZ)
			cmd.DisplayTime("",*lookupZ)
		} else if len(lookupCmd.Args()) != 0 {
            // combine all non-flag arguments to create a city name
			city := strings.Join(lookupCmd.Args(), " ")
			// Display the current local time based on the given city
			cmd.DisplayTime(city,"")
		}
	default:
		printError("", "\n Error: Unknown command")

	}

	// switch os.Args[1] {
	// case "add", "lookup":
	// 	if len(os.Args) < 3 {

	// 		if os.Args[1] == "lookup" {
	// 			printLookupHelp()
	// 			return
	// 		} else if os.Args[1] == "add" {
	// 			printError("Error: Incomplete command")
	// 			return
	// 		}

	// 	}
	// 	city := strings.Join(os.Args[2:], " ")
	// 	if os.Args[1] == "lookup" {
	// 		// Display the current local time based on the given city
	// 		cmd.DisplayTime(city)
	// 	} else if os.Args[1] == "add" {
	// 		// Add tz to local storage ~/.ktz_cli
	// 		fmt.Println("Added tz", city)
	// 	}
	// case "help":
	// 	fmt.Println("Show help")
	// 	return

	// case "view":
	// 	fmt.Println("Show all saved timezones")
	// 	return
	// default:
	// 	printError("Error: Unknown command")
	// }

}

func printError(errCmd, errMsg string) {
	fmt.Println(errMsg)
	switch errCmd {
	case "lookup":
		fmt.Println(" Usage: ktz lookup [options] [city/country]")
	default:
	}
	fmt.Println(" For more information, try 'app help'")
	fmt.Println("")

}

func printLookupHelp() {
	fmt.Println("Usage: ktz lookup [options] [city/country]")
	fmt.Println()
	fmt.Println("Look up the current time for a city,zone or country")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -tz string  Specify the timezone for a city")
	fmt.Println("  -z  string  Specify the zone for a region")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  ktz lookup \"New York\"")
	fmt.Println("  ktz lookup -tz=America/New_York")
	fmt.Println("  ktz lookup -z=pst")
}

// ktz lookup -tz="America/New_York"
// ktz lookup nepal
// ktz lookup kathmandu
// ktz lookup -z="pst"

//ktz add -tz="America/New_York" --- only add the tz you know
//ktz remove -tz="America/New_York" --- only remove the tz you know
//ktz remove -interactive mode

//ktz view-all

// Usage: myapp.py [OPTIONS] [THINGS]...

//   Combines things into a single element

// Options:
//   -t TEXT  Data type  [default: int]
//   -o TEXT  Output format
//   --help   Show this message and exit.
