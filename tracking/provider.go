package tracking

// Provider defines a shipping provider
type Provider interface {
	GetID() string
	GetTrackingURL(trackingNumber string) string
	GetTrackingEvents(trackingNumber string) ([]Event, error)
}
