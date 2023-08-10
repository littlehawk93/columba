package ups

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/littlehawk93/columba/providers/utils"
	"github.com/littlehawk93/columba/tracking"
)

const (
	trackingURL     string = "https://www.ups.com/track?loc=en_US"
	id              string = "UPS"
	apiUrl          string = "https://www.ups.com/track/api/Track/GetStatus?loc=en_US"
	tokenCookieName string = "X-XSRF-TOKEN-ST"
	tokenHeaderName string = "X-XSRF-TOKEN"
)

// Provider extends the service.Provider interface for UPS
type Provider struct {
}

// GetID returns the USPS provider ID
func (me *Provider) GetID() string {
	return id
}

// GetTrackingURL returns the UPS tracking URL for a given tracking number
func (me *Provider) GetTrackingURL(trackingNumber string) string {
	return trackingURL
}

// GetTrackingEvents get all tracking events for a given tracking number
func (me *Provider) GetTrackingEvents(trackingNumber string, options tracking.Options) ([]tracking.Event, error) {

	client, err := utils.NewClientWithCookies()

	if err != nil {
		return nil, err
	}

	req, err := utils.CreateRequest(http.MethodGet, trackingURL, nil)

	if err != nil {
		return nil, err
	}

	if _, err = utils.GetResponseBytes(client, req); err != nil {
		return nil, err
	}

	if req, err = utils.CreateJSONRequest(http.MethodPost, apiUrl, &trackingRequest{
		Locale:          "en_US",
		TrackingNumbers: []string{trackingNumber},
	}); err != nil {
		return nil, err
	}

	token := getTokenFromCookies(client)

	req.Header.Set(tokenHeaderName, token)

	response := &trackingResponse{}

	if err = utils.GetResponseJSON(client, req, response); err != nil {
		return nil, err
	}

	events := make([]tracking.Event, 0)

	if len(response.Details) > 0 {
		for _, detail := range response.Details {
			if len(detail.ProgressDetails) > 0 {
				for _, progressDetail := range detail.ProgressDetails {

					timestamp, err := progressDetail.GetTimestamp()

					if err != nil {
						return nil, err
					}

					events = append(events, tracking.Event{
						EventText: progressDetail.Activity,
						Location:  progressDetail.GetLocation(),
						Timestamp: timestamp,
						Details:   "",
					})
				}
			}
		}
	}

	return events, nil
}

func getTokenFromCookies(client *http.Client) string {

	if u, err := url.Parse(trackingURL); err == nil {
		for _, cookie := range client.Jar.Cookies(u) {
			if strings.TrimSpace(strings.ToUpper(cookie.Name)) == tokenCookieName {
				return cookie.Value
			}
		}
	}
	return ""
}
