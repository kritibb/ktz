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
// It expects command-line arguments in the form of ktz lookup [options] [city],
// retrieves the corresponding timezone, and displays the current local time in that timezone.
func main() {

	//define subcommand `lookup` and its flags
	lookupCmd := flag.NewFlagSet("lookup", flag.ExitOnError)
	lookupC := lookupCmd.String("c", "", "country name/code like `Nepal` or `NP`")
	lookupZ := lookupCmd.String("z", "", "`timezones` like `Asia/Kathmandu` or `PST`")

	// Check if subcommands like "lookup" is provided
	if len(os.Args) < 2 {
		printError("", "\n Error: Incomplete command")
		return
	}

	switch os.Args[1] {
	case "lookup":
		lookupCmd.Parse(os.Args[2:])
		//no flags or positional argument provided
		if *lookupZ == "" && *lookupC == "" && len(lookupCmd.Args()) == 0 {
			printError("lookup", "\n Error: Incomplete command")
			return
		}
		// Invalid combination: more than one flag or both flags and positional argument
		if (*lookupZ != "" && *lookupC != "") ||
			(*lookupZ != "" && len(lookupCmd.Args()) != 0) ||
			(*lookupC != "" && len(lookupCmd.Args()) != 0) {
			printError("lookup", "\nError: Use only a flag [-z] or [-c] or [city/country]")
			return
		}

		//handle -z flag
		if *lookupZ != "" {
			cmd.ResolveTimezone("", "", *lookupZ)
		} else if *lookupC != "" { //handle -c flag
			cmd.ResolveTimezone("", *lookupC, "")
		} else if len(lookupCmd.Args()) != 0 {
			// combine all non-flag arguments to create a city name
			city := strings.Join(lookupCmd.Args(), " ")
			// Display the current local time based on the given city
			cmd.ResolveTimezone(city, "", "")
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
		fmt.Println(" Usage: ktz lookup [options] [city]")
	default:
	}
	fmt.Println(" For more information, try 'ktz help'")
	fmt.Println("")

}

func printLookupHelp() {
	fmt.Println("Usage: ktz lookup [options] [city/country]")
	fmt.Println()
	fmt.Println("Look up the current time for a city,zone or country")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -tz string  Specify the timezone for a city or a zone")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  ktz lookup \"New York\"")
	fmt.Println("  ktz lookup -z=America/New_York or -z=PST")
	fmt.Println("  ktz lookup -c=NP or -c='Nepal'")
}

// ktz lookup -z="America/New_York"/ "pst"
// ktz lookup nepal
// ktz lookup kathmandu

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

//replace fuzzy with prefix
