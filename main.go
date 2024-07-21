// Package main is the entry point for the timezone application.
// It retrieves the timezone for a given city name 
// and displays the current local time for that timezone.
package main

import (
	"fmt"
	"github.com/kritibb/ktz/cmd"
	"os"
    "strings"
)

// main is the entry point of the application.
// It expects the city name as a command-line argument, retrieves the corresponding timezone,
// and displays the current local time in that timezone.
func main() {
    // Check if the city name is provided as a command-line argument
	if len(os.Args) < 2 {
		fmt.Println("Please provide a city name!")
		return
	}
    city := os.Args[1:]
    // Join the command-line arguments to form the city name
    cityString:=strings.Join(city, " ")


    // Display the current local time based on the given cityString
    cmd.DisplayTime(cityString)


}
