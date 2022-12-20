package tracking

import (
	"time"
)

// Package defines a single package to track
type Package struct {
	ID                    int
	Status                PackageStatus
	TrackingNumber        string
	ServiceID             string
	Label                 string
	TrackingURL           string
	CreatedOn             time.Time
	LastUpdatedOn         time.Time
	Origin                *Location
	Destination           *Location
	EstimatedDeliveryDate *time.Time
	Events                []Event
}
