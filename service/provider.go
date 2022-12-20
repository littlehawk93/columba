package service

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/PuerkitoBio/goquery"
	"github.com/littlehawk93/columba/tracking"
)

// TrackingEventParser definition for a function that returns a tracking event from an HTML selection
type TrackingEventParser func(doc *goquery.Selection) (tracking.Event, error)

// Provider defines a shipping provider
type Provider interface {
	GetID() string
	GetTrackingURL(trackingNumber string) string
	GetTrackingEvents(trackingNumber string) ([]tracking.Event, error)
}

// ParseTrackingEvents helper function for parsing HTML from a given URL and parsing a set of events from the returned response
func ParseTrackingEvents(url, selector string, parser TrackingEventParser) ([]tracking.Event, error) {

	log.Println(url)
	log.Println(selector)

	results := make([]tracking.Event, 0)

	res, err := http.Get(url)

	if err != nil {
		return results, err
	}

	defer res.Body.Close()

	f, err := os.Create("tmp.html")

	if err != nil {
		return results, err
	}

	defer f.Close()

	io.Copy(f, res.Body)

	return results, nil

	/*
		document, err := goquery.NewDocumentFromReader(f)

		if err != nil {
			return results, err
		}

		document.Find("div").Each(func(i int, s *goquery.Selection) {

			if err != nil {
				return
			}

			evt, err2 := parser(s)

			if err2 != nil {
				log.Fatal(err2)
			} else if evt.EventText != "" || evt.Location != nil || evt.Timestamp != nil {
				results = append(results, evt)
			}
		})

		return results, err
	*/
}
