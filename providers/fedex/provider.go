package fedex

import (
	"fmt"
	"net/url"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/littlehawk93/columba/providers/utils"
	"github.com/littlehawk93/columba/tracking"
)

const (
	historyButtonSelector         string = "div.travel-history-link button"
	trackingEventTableRowSelector string = "tr.travel-history-table__row"
	trackingEventTableDateClass   string = "travel-history-table__scan-event-date"
	urlFormat                     string = "https://www.fedex.com/fedextrack/?trknbr=%s"
	id                            string = "FedEx"
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
func (me *Provider) GetTrackingEvents(trackingNumber string, options tracking.Options) ([]tracking.Event, error) {

	actions := []chromedp.Action{
		chromedp.WaitVisible(historyButtonSelector),
		chromedp.Click(historyButtonSelector),
	}

	return utils.ParseTrackingEventsHeadlessChrome(me.GetTrackingURL(trackingNumber), trackingEventTableRowSelector, options, actions, trackingEventParserChromeDp)
}

func trackingEventParserChromeDp(n *cdp.Node) ([]tracking.Event, error) {
	return nil, nil
}
