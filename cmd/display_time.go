package cmd

import (
	"fmt"
	"time"
)

// DisplayTime prints the current time in the specified location.
// The 'location' parameter should be a prefix or a complete city.
// The 'zone' parameter should be a timezone like 'Asia/Kathmandu' or 'pst'
// Time is displayed based on either location or zone
// If the location is invalid, an error message is printed, else 
// timezone is displayed based on it.
func DisplayTime(location, zone string) {
	if zone == "" {
		//Get potential cities based on the provided location string
		cityList, err := getMatchingCities(location)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		//if cityList consists only one timezone, display time for that timezone, else
		//display list of potential timezones that came from the location string (prefix)
		if len(cityList) == 1 {
			tableViewDateTime(cityList, zone)
		} else {
			listViewTz(cityList)
		}
	} else {
		tableViewDateTime([]string{}, zone)

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
