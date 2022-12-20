package usps

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/littlehawk93/columba/service"
	"github.com/littlehawk93/columba/tracking"
)

const (
	urlFormat                     string = "https://tools.usps.com/go/TrackConfirmAction?tLabels=%s"
	id                            string = "USPS"
	trackingEventItemSelector     string = "div.tb-step"
	trackingEventDateSelector     string = "p.tb-date"
	trackingEventStatusSelector   string = "p.tb-status-detail"
	trackingEventLocationSelector string = "p.tb-location"
	trackingEventDateFormat       string = "January 2, 2006, 3:04 pm"
	trackingEventShortDateFormat  string = "January 2, 2006"
)

// Provider extends the service.Provider interface for USPS
type Provider struct {
}

// GetID returns the USPS provider ID
func (me Provider) GetID() string {
	return id
}

// GetTrackingURL returns the USPS tracking URL for a given tracking number
func (me Provider) GetTrackingURL(trackingNumber string) string {
	return fmt.Sprintf(urlFormat, url.QueryEscape(trackingNumber))
}

// GetTrackingEvents get all tracking events for a given tracking number
func (me Provider) GetTrackingEvents(trackingNumber string) ([]tracking.Event, error) {

	return service.ParseTrackingEvents(me.GetTrackingURL(trackingNumber), trackingEventItemSelector, trackingEventParser)
}

func trackingEventParser(s *goquery.Selection) (tracking.Event, error) {

	event := tracking.Event{}

	if dateElem := s.Find(trackingEventDateSelector).First(); dateElem != nil && len(strings.TrimSpace(dateElem.Text())) > 0 {
		date, err := cleanAndParseDate(dateElem.Text())

		if err != nil {
			return event, err
		}
		event.Timestamp = &date
	}

	if locationElem := s.Find(trackingEventLocationSelector).First(); locationElem != nil && len(strings.TrimSpace(locationElem.Text())) > 0 {
		location, err := cleanAndParseLocation(locationElem.Text())

		if err != nil {
			return event, err
		}
		event.Location = location
	}

	if statusElem := s.Find(trackingEventStatusSelector).First(); statusElem != nil && len(strings.TrimSpace(statusElem.Text())) > 0 {
		event.EventText = strings.TrimSpace(statusElem.Text())
	}

	if attr, ok := s.Attr("class"); ok && strings.Contains(attr, "current-step") {
		event.IsCurrent = true
	} else {
		event.IsCurrent = false
	}

	return event, nil
}

func cleanAndParseDate(dateStr string) (time.Time, error) {

	dateStr = regexp.MustCompile(`\s\s+`).ReplaceAllString(strings.TrimSpace(dateStr), " ")

	result, err := time.Parse(trackingEventDateFormat, dateStr)

	if err == nil {
		return result, err
	}

	return time.Parse(trackingEventShortDateFormat, dateStr)
}

func cleanAndParseLocation(locationStr string) (*tracking.Location, error) {

	location := &tracking.Location{}

	words := strings.Fields(locationStr)

	for i := 0; i < len(words); i++ {
		words[i] = cleanLocationWord(words[i])
	}

	stateIndex := -1
	zipIndex := -1

	for i, word := range words {
		if tracking.IsStateAbbreviation(word) && stateIndex == -1 {
			stateIndex = i
		} else if regexp.MustCompile(`^[0-9]{5}(\-[0-9]{4})?$`).MatchString(word) && zipIndex == -1 {
			zipIndex = i
		}

		if stateIndex != -1 && zipIndex == -1 {
			break
		}
	}

	if stateIndex > -1 {
		location.City = strings.Join(words[:stateIndex], " ")
		location.State = strings.ToUpper(words[stateIndex])

		if zipIndex > -1 {
			location.Zip = words[zipIndex]
			if zipIndex+1 < len(words) {
				location.Facility = strings.Join(words[zipIndex+1:], " ")
			}
		} else if stateIndex+1 < len(words) {
			location.Facility = strings.Join(words[stateIndex+1:], " ")
		}
	} else {
		location.Facility = strings.Join(words, " ")
	}

	return location, nil
}

func cleanLocationWord(word string) string {

	return regexp.MustCompile(`(^[^A-Za-z]+)|([^A-Za-z]+$)`).ReplaceAllString(word, "")
}
