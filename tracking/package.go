package tracking

import (
	"time"

	"github.com/littlehawk93/columba/service"
)

// Package defines a single package to track
type Package struct {
	ID                    int
	Status                PackageStatus
	ServiceProvider       *service.Provider
	TrackingNumber        string
	Label                 string
	TrackingURL           string
	CreatedOn             time.Time
	LastUpdatedOn         time.Time
	Origin                *Location
	Destination           *Location
	EstimatedDeliveryDate time.Time
	Events                []Event
}
