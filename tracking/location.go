package tracking

import (
	"fmt"
	"strings"
)

type Location struct {
	Facility string
	City     string
	State    string
	Zip      string
}

// String serialize this location to a string
func (me Location) String() string {

	result := fmt.Sprintf("%s|%s|%s|%s", strings.TrimSpace(strings.ReplaceAll(me.City, "|", "")), strings.TrimSpace(strings.ReplaceAll(me.State, "|", "")), strings.TrimSpace(strings.ReplaceAll(me.Zip, "|", "")), strings.TrimSpace(strings.ReplaceAll(me.Facility, "|", "")))

	return strings.Trim(result, "|")
}

// GetOriginLocationFromEvents get the originating location from a set of events
func GetOriginLocationFromEvents(events []Event) *Location {

	if len(events) == 0 {
		return nil
	}

	return events[0].Location
}

// NewLocationFromString parses a serialized location string into a location struct, returns nil if the string is empty
func NewLocationFromString(locationString string) *Location {

	locationString = strings.TrimSpace(locationString)

	if locationString == "" {
		return nil
	}

	words := strings.Split(locationString, "|")

	loc := &Location{}

	for i, word := range words {

		switch i {
		case 0:
			loc.City = word
		case 1:
			loc.State = word
		case 2:
			if IsZipCode(word) {
				loc.Zip = word
			} else {
				loc.Facility = word
			}
		default:
			if IsZipCode(word) {
				loc.Zip = word
			} else {
				loc.Facility = word
			}
		}
	}

	return loc
}
