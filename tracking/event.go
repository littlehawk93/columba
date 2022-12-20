package tracking

import "time"

// Event defines a package tracking event. Either delivery or a package arriving at a post office, etc
type Event struct {
	ID        int
	PackageID int
	EventText string
	Details   string
	Location  *Location
	Timestamp *time.Time
	IsCurrent bool
}
