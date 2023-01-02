package ups

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/littlehawk93/columba/providers/utils"
	"github.com/littlehawk93/columba/tracking"
)

const (
	trackingURL string = "https://www.ups.com/track?loc=en_US"
	id          string = "UPS"
	apiUrl      string = "https://www.ups.com/track/api/Track/GetStatus?loc=en_US"
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
func (me *Provider) GetTrackingEvents(trackingNumber string) ([]tracking.Event, error) {

	requestBytes, err := json.Marshal(&trackingRequest{
		Locale:          "en_US",
		TrackingNumbers: []string{trackingNumber},
	})

	if err != nil {
		return nil, err
	}

	requestBody := bytes.NewReader(requestBytes)
	req, err := http.NewRequest(http.MethodPost, apiUrl, requestBody)

	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", utils.GetUserAgent())
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	resBytes, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	response := trackingResponse{}

	if err = json.Unmarshal(resBytes, &response); err != nil {
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
						Timestamp: &timestamp,
						Details:   "",
						IsCurrent: false,
					})
				}
			}
		}
	}

	return events, nil
}
