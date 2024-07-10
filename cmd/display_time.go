package cmd

import (
	"fmt"
    "time"
)



func DisplayTime(location string){
    loc, err := time.LoadLocation(location)
    if err != nil {
		fmt.Printf("Error loading time zone: %v\n", err)
		return
	}
    // Get current time in UTC
	utcTime := time.Now().UTC()

	// Convert UTC time to local time of the specified location
	localTime := utcTime.In(loc)

	// Print the local time
    const customFormat = "Mon, 02 Jan 2006 03:04:05 PM"
	fmt.Printf("Current time in %s: %s\n", location, localTime.Format(customFormat))

}
