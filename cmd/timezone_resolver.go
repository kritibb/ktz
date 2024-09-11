package cmd

import (
	"fmt"
	"time"
)

// ResolveTimeZone prints the current time in the specified location.
// The 'city' parameter should be a prefix or a complete city.
// The 'zone' parameter should be a timezone like 'Asia/Kathmandu' or 'PST'
// The 'country' parameter should be a prefix/ complete country or a country code.
// Time is displayed based on either location or zone
// If the location is invalid, an error message is printed, else
// timezone is displayed based on it.
func ResolveTimezone(city, country, zone string) {
	if zone != "" {
		tableViewDateTime("", zone)
        return
	}
	//Get potential location based on the provided city/country string
    locationList, err := getMatchingLocation(city, country)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	//if locationList consists only one location, display time for that location, else
	//display list of potential timezones that came from the location string (prefix)
	if len(locationList) == 1 {
		tableViewDateTime(locationList[0], zone)
	} else {
		listViewTz(locationList)
	}

}

// formatTime displays time for a given timeZone in a specified fromat.
//
// Parameters:
//   - tz: The timezone in string format.
//
// Returns:
//   - string: time of a particular tz in a certain format
//   - error: any error message if time.LoadLocation does not find the given tz
func formatTime(tz string) (string, error) {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return "", err
	}
	// Get current time in UTC
	utcTime := time.Now().UTC()

	// Convert UTC time to local time of the specified location
	localTime := utcTime.In(loc)

	// Print the local time
	const customFormat = "Mon, 02 Jan 2006 03:04:05 PM"
	return localTime.Format(customFormat), nil
}
