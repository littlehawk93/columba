package ups

type trackingRequest struct {
	Locale          string   `json:"Locale"`
	Requester       string   `json:"Requester"`
	TrackingNumbers []string `json:"TrackingNumber"`
	ReturnToValue   string   `json:"returnToValue"`
}
