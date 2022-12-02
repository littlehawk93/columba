package tracking

import "time"

// Event defines a package tracking event. Either delivery or a package arriving at a post office, etc
type Event struct {
	ID        int
	Text      string
	Location  *Location
	EventCode string
	Timestamp time.Time
}
