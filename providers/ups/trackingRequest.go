package ups

type trackingRequest struct {
	Locale          string   `json:"Locale"`
	TrackingNumbers []string `json:"TrackingNumber"`
}
