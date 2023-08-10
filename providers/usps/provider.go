package usps

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/littlehawk93/columba/cdputils"
	"github.com/littlehawk93/columba/providers/utils"
	"github.com/littlehawk93/columba/tracking"
)

const (
	urlFormat                      string = "https://tools.usps.com/go/TrackConfirmAction?tLabels=%s"
	id                             string = "USPS"
	trackingShowMoreButtonSelector string = "a.expand-collapse-history"
	trackingEventItemSelector      string = "div.tb-step"
	trackingEventDateClass         string = "tb-date"
	trackingEventStatusClass       string = "tb-status-detail"
	trackingEventLocationClass     string = "tb-location"
	trackingEventDateFormat        string = "January 2, 2006, 3:04 pm"
	trackingEventShortDateFormat   string = "January 2, 2006"
)

// Provider extends the service.Provider interface for USPS
type Provider struct {
}

// GetID returns the USPS provider ID
func (me *Provider) GetID() string {
	return id
}

// GetTrackingURL returns the USPS tracking URL for a given tracking number
func (me *Provider) GetTrackingURL(trackingNumber string) string {
	return fmt.Sprintf(urlFormat, url.QueryEscape(trackingNumber))
}

// GetTrackingEvents get all tracking events for a given tracking number
func (me *Provider) GetTrackingEvents(trackingNumber string, options tracking.Options) ([]tracking.Event, error) {

	actions := []chromedp.Action{
		chromedp.WaitVisible(trackingShowMoreButtonSelector),
		chromedp.Click(trackingShowMoreButtonSelector),
	}

	return utils.ParseTrackingEventsHeadlessChrome(me.GetTrackingURL(trackingNumber), trackingEventItemSelector, options, actions, trackingEventParserChromeDp)
}

func trackingEventParserChromeDp(n *cdp.Node) ([]tracking.Event, error) {

	event := tracking.Event{}

	if dateNode := cdputils.FirstChildByClass(n, trackingEventDateClass, true); dateNode != nil {
		date, err := cleanAndParseDate(cdputils.NodeText(dateNode))

		if err != nil {
			return []tracking.Event{event}, fmt.Errorf("[trackingEventParserChromeDp] error parsing event timestamp: %w", err)
		}
		event.Timestamp = &date
	}

	if locationNode := cdputils.FirstChildByClass(n, trackingEventLocationClass, true); locationNode != nil {
		location, err := cleanAndParseLocation(cdputils.NodeText(locationNode))

		if err != nil {
			return []tracking.Event{event}, fmt.Errorf("[trackingEventParserChromeDp] error parsing event location: %w", err)
		}
		event.Location = location
	}

	if statusNode := cdputils.FirstChildByClass(n, trackingEventStatusClass, true); statusNode != nil {
		event.EventText = strings.TrimSpace(cdputils.NodeText(statusNode))
	}

	return []tracking.Event{event}, nil
}

func cleanAndParseDate(dateStr string) (time.Time, error) {

	dateStr = regexp.MustCompile(`\s+`).ReplaceAllString(strings.TrimSpace(dateStr), " ")

	if idx := strings.Index(dateStr, " am "); idx != -1 && len(dateStr) > idx+4 {
		dateStr = dateStr[:idx+3]
	} else if idx := strings.Index(dateStr, " pm "); idx != -1 && len(dateStr) > idx+4 {
		dateStr = dateStr[:idx+3]
	}

	result, err := time.ParseInLocation(trackingEventDateFormat, dateStr, time.Local)

	if err == nil {
		return result, err
	}

	return time.ParseInLocation(trackingEventShortDateFormat, dateStr, time.Local)
}

func cleanAndParseLocation(locationStr string) (*tracking.Location, error) {

	locationStr = strings.TrimSpace(locationStr)

	location := &tracking.Location{}

	words := strings.Fields(locationStr)

	for i := 0; i < len(words); i++ {
		words[i] = utils.CleanLocationWord(words[i])
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
		} else if stateIndex+1 < len(words) {
			location.Facility = strings.Join(words[stateIndex+1:], " ")
		}
	} else {
		location.Facility = strings.Join(words, " ")
	}

	return location, nil
}
