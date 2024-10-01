package cmd

import (
	"fmt"
	"github.com/kritibb/ktz/tzdata"
	"strings"
	"time"
)

type zoneInfo struct {
	formattedTime string
	timezoneName  string
	abbreviation  string
}

type locationInfo struct {
	country       string // The country name or code (e.g., USA, IN)
	city          string // The city name (e.g., New York, London)
	timezone      string // The full timezone name (e.g., America/New_York)
	formattedTime string // The time formatted according to the timezone
}

// ResolveTimeZone prints the current time in the specified location.
// The 'city' parameter should be a prefix or a complete city.
// The 'zone' parameter should be a timezone like 'Asia/Kathmandu' or 'PST'
// The 'country' parameter should be a prefix/ complete country or a country code.
// Time is displayed based on either location or zone
// If the location is invalid, an error message is printed, else
// timezone is displayed based on it.
func ResolveTimezone(city, country, zone string) {
	if zone != "" {
		zoneData, err := getDataFromZone(zone)
		if err != nil {
			fmt.Printf("\nError: %v\n", err)
			return
		}
		renderZoneInfoTable(zoneData)
		return
	}
	//Get potential location based on the provided city/country string
	locationList, err := getMatchingLocation(city, country)
	if err != nil {
		fmt.Println("\nError:", err)
		return
	}
	currentLocationData, errLocation := getDataFromLocation(locationList)
	if errLocation != nil {
		fmt.Printf("\nError: %v\n", err)
		return
	}
	renderDateTimeTableFromLocation(currentLocationData)
	return
}

// getDataFromLocation gives locationInfo based on given locationList (city/country)
//
// Parameters:
//   - locationList: A country or city list
//
// Returns:
//   - locationInfo:
//   - error: any error message if zone does not exist

func getDataFromLocation(locationList []string) (locationInfo, error) {
	//if locationList consists only one location, set timezone, city and country based on that
	if len(locationList) == 1 {
		location := locationList[0]
		if val, ok := tzdata.CityToIanaTimezone[location]; ok {
			locationData.timezone = val["tz"]
			locationData.city = location
			locationData.country = val["country"]
		} else if val, ok := tzdata.CountryToIanaTimezone[location]; ok {
			if len(val) > 1 {
				locationData.country = location
				listViewTz(val)
			}
			locationData.timezone = val[0]
			locationData.country = location
		}
	} else {
		listViewTz(locationList)
	}
	datetime, err := formatTime(locationData.timezone)
	if err != nil {
		return locationData, err
	}
	locationData.formattedTime = datetime
	return locationData, err
}

// getDataFromZone gives zoneInfo based on given zone name
//
// Parameters:
//   - zone: The timezone in string format;
//          could either be abbreviation like pst or full name like "Asia/Kathmandu"
//
// Returns:
//   - zoneInfo:
//   - error: any error message if zone does not exist

func getDataFromZone(zone string) (zoneInfo, error) {
	var zoneData zoneInfo
	var err error
	if len(zone) < 6 {
		zoneData.abbreviation = zone
		if val, ok := tzdata.AbbToIanaTimezone[strings.ToUpper(zone)]; !ok {
			err = fmt.Errorf("Zone abbreviation '%v' not found.\n\n", zone)
			return zoneData, err
		} else {
			zoneData.timezoneName = val
		}
	} else {
		zoneData.timezoneName = zone
	}
	datetime, err := formatTime(zoneData.timezoneName)
	if err != nil {
		return zoneData, err
	}
	zoneData.formattedTime = datetime
	return zoneData, err
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
