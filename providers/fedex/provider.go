package fedex

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/littlehawk93/columba/providers/utils"
	"github.com/littlehawk93/columba/tracking"
)

var trackingEventDateFormats []string = []string{
	"1/02/2003 at 3:04 PM",
	"1/02/2003 3:04 PM",
	"1/02/2003",
}

const (
	urlFormat                            string = "https://www.fedex.com/fedextrack/?tracknumbers=%s"
	id                                   string = "FedEx"
	trackingEventItemSelector            string = "div.shipment-status-progress-container div.shipment-status-progress-step"
	trackingEventStatusSelector          string = "span.shipment-status-progress-step-label"
	trackingEventInfoSelector            string = "h5.shipment-status-progress-step-label-info"
	trackingEventLocationAndDateSelector string = "div.shipment-status-progress-step-content span.shipment-status-progress-step-label-content"
)

// Provider extends the service.Provider interface for FedEx
type Provider struct {
}

// GetID returns the FedEx provider ID
func (me *Provider) GetID() string {
	return id
}

// GetTrackingURL returns the FedEx tracking URL for a given tracking number
func (me *Provider) GetTrackingURL(trackingNumber string) string {
	return fmt.Sprintf(urlFormat, url.QueryEscape(trackingNumber))
}

// GetTrackingEvents get all tracking events for a given tracking number
func (me *Provider) GetTrackingEvents(trackingNumber string) ([]tracking.Event, error) {

	return utils.ParseTrackingEvents(me.GetTrackingURL(trackingNumber), trackingEventItemSelector, trackingEventParser)
}

func trackingEventParser(s *goquery.Selection) (tracking.Event, error) {

	event := tracking.Event{}

	if statusElem, ok := utils.ElemExistsWithData(s, trackingEventStatusSelector); ok {
		event.EventText = strings.TrimSpace(statusElem.Text())
	} else if statusElem, ok := utils.ElemExistsWithData(s, trackingEventInfoSelector); ok {
		event.EventText = strings.TrimSpace(statusElem.Text())
	}

	s.Find(trackingEventLocationAndDateSelector).Each(func(idx int, item *goquery.Selection) {

		if event.Timestamp == nil {
			if date, err := cleanAndParseDate(item.Text()); err == nil {
				event.Timestamp = &date
				return
			}
		}

		if event.Location == nil {
			if location, err := cleanAndParseLocation(item.Text()); err == nil {
				event.Location = location
				return
			}
		}
	})

	if attr, ok := s.Attr("class"); ok && strings.Contains(attr, "active") {
		event.IsCurrent = true
	} else {
		event.IsCurrent = false
	}

	return event, nil
}

func cleanAndParseDate(dateStr string) (time.Time, error) {

	dateStr = regexp.MustCompile(`\s+`).ReplaceAllString(strings.TrimSpace(dateStr), " ")

	if idx := strings.Index(dateStr, " am "); idx != -1 && len(dateStr) > idx+4 {
		dateStr = dateStr[:idx+3]
	} else if idx := strings.Index(dateStr, " pm "); idx != -1 && len(dateStr) > idx+4 {
		dateStr = dateStr[:idx+3]
	}

	var result time.Time
	var err error

	for _, format := range trackingEventDateFormats {
		if result, err = time.Parse(format, dateStr); err == nil {
			break
		}
	}

	return result, err
}

func cleanAndParseLocation(locationStr string) (*tracking.Location, error) {

	locationStr = strings.TrimSpace(locationStr)

	location := &tracking.Location{
		Facility: "",
	}

	words := strings.Fields(locationStr)

	for i := 0; i < len(words); i++ {
		words[i] = cleanLocationWord(words[i])
	}

	stateIndex := -1
	zipIndex := -1

	for i, word := range words {
		if tracking.IsStateAbbreviation(word) && stateIndex == -1 {
			stateIndex = i
		} else if tracking.IsZipCode(word) && zipIndex == -1 {
			zipIndex = i
		}

		if stateIndex != -1 && zipIndex != -1 {
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
		}
	}

	return location, nil
}

func cleanLocationWord(word string) string {

	return regexp.MustCompile(`(^[^A-Za-z0-9]+)|([^A-Za-z0-9]+$)`).ReplaceAllString(word, "")
}
