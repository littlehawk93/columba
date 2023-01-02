package utils

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/littlehawk93/columba/tracking"
)

const (
	// userAgentWin10Edge Win10 MS Edge User Agent
	userAgentWin10Edge       string = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.135 Safari/537.36 Edge/12.246"
	cleanLocationRegexCutset string = `(^[^A-Za-z0-9]+)|([^A-Za-z0-9]+$)`
)

// TrackingEventParser definition for a function that returns a tracking event from an HTML selection
type TrackingEventParser func(doc *goquery.Selection) (tracking.Event, error)

// ParseTrackingEvents helper function for parsing HTML from a given URL and parsing a set of events from the returned response
func ParseTrackingEvents(url, selector string, parser TrackingEventParser) ([]tracking.Event, error) {

	results := make([]tracking.Event, 0)

	res, err := http.Get(url)

	if err != nil {
		return results, err
	}

	defer res.Body.Close()

	document, err := goquery.NewDocumentFromReader(res.Body)

	if err != nil {
		return results, err
	}

	document.Find(selector).Each(func(i int, s *goquery.Selection) {

		if err != nil {
			return
		}

		evt, err2 := parser(s)

		if err2 != nil {
			err = err2
		} else if evt.EventText != "" || evt.Location != nil || evt.Timestamp != nil {
			results = append(results, evt)
		}
	})

	return results, err
}

// ElemExistsWithData returns the first matching element or nil if no element found and returns true if the elem has any text in it
func ElemExistsWithData(s *goquery.Selection, selector string) (*goquery.Selection, bool) {

	elem := s.Find(selector).First()

	if elem == nil {
		return nil, false
	}

	return elem, len(strings.TrimSpace(elem.Text())) > 0
}

// GetUserAgent retreive a user agent to make Columba requests appear as desktop client apps to tracking websites
func GetUserAgent() string {
	return userAgentWin10Edge
}

// CleanLocationWord cleans a particular city or state name for a location to standardize naming format across all shipping providers
func CleanLocationWord(word string) string {
	return strings.ToUpper(regexp.MustCompile(cleanLocationRegexCutset).ReplaceAllString(word, ""))
}
