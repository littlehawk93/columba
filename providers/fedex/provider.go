package fedex

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/littlehawk93/columba/providers/utils"
	"github.com/littlehawk93/columba/tracking"
)

const (
	propertiesUrl  string = "https://www.fedex.com/fedextrack/properties/WTRKProperties.json"
	urlFormat      string = "https://www.fedex.com/fedextrack/?tracknumbers=%s"
	tokenUrlFormat string = "https://api.fedex.com/auth/oauth/v2/token?grant_type=client_credentials&client_id=%s&client_secret=%s&scope=oob"
	apiUrl         string = "https://api.fedex.com/track/v2/shipments"
	id             string = "FedEx"
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

	client, err := utils.NewClientWithCookies()

	if err != nil {
		return nil, err
	}

	req, err := utils.CreateRequest(http.MethodGet, fmt.Sprintf(urlFormat, url.QueryEscape(trackingNumber)), nil)

	if err != nil {
		return nil, err
	}

	if _, err = utils.GetResponseBytes(client, req); err != nil {
		return nil, err
	}

	properties := &propertiesResponse{}

	if req, err = utils.CreateRequest(http.MethodGet, propertiesUrl, nil); err != nil {
		return nil, err
	}

	if err = utils.GetResponseJSON(client, req, properties); err != nil {
		return nil, err
	}

	if req, err = utils.CreateRequest(http.MethodPost, fmt.Sprintf(tokenUrlFormat, properties.API.ClientID, properties.API.ClientSecret), nil); err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	token := &authToken{}

	if err = utils.GetResponseJSON(client, req, token); err != nil {
		return nil, err
	}

	trackingRequest := newTrackingRequest(trackingNumber)

	if req, err = utils.CreateJSONRequest(http.MethodPost, apiUrl, &trackingRequest); err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Origin", "https://www.fedex.com")
	req.Header.Set("Referer", fmt.Sprintf(urlFormat, trackingNumber))
	req.Header.Set("x-clientid", "WTRK")
	req.Header.Set("x-version", "1.0.0")
	req.Header.Set("x-locale", "en_US")
	token.setRequestAuthorization(req)
	trackingResponse := &trackingResponse{}

	if err = utils.GetResponseJSON(client, req, trackingResponse); err != nil {
		return nil, err
	}

	events := make([]tracking.Event, 0)

	if trackingResponse != nil && len(trackingResponse.Output.Packages) > 0 {
		for _, responsePackage := range trackingResponse.Output.Packages {
			if len(responsePackage.ScanEventList) > 0 {
				for _, evt := range responsePackage.ScanEventList {

					t, err := evt.GetTimestamp()

					if err != nil {
						return nil, err
					}

					event := tracking.Event{
						EventText: evt.Status,
						Details:   evt.Details,
						Location:  evt.GetLocation(),
						Timestamp: t,
					}

					events = append(events, event)
				}
			}
		}
	}

	return events, nil
}
