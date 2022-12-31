package tracking

import (
	"database/sql"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"hash/crc32"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	eventTimestampDateFormat string = "2006-01-02 15:04:05"
)

// Event defines a package tracking event. Either delivery or a package arriving at a post office, etc
type Event struct {
	id        uuid.UUID
	PackageID int
	EventText string
	Details   string
	Location  *Location
	Timestamp *time.Time
	IsCurrent bool
}

// GetID get the ID of this event
func (me Event) GetID() uuid.UUID {

	return me.id
}

// SetID generates and sets the ID property of this event. ID is generated using a deterministic algorithm in order to make comparing equivalent events easy
func (me *Event) SetID() {

	idBytes := make([]byte, 16)

	if me.Timestamp == nil {
		panic("Event timestamp is nil")
	}

	minutes := uint32(me.Timestamp.Sub(time.Date(1970, 1, 1, 0, 0, 0, 0, time.Local)).Minutes())

	tmpBytes := make([]byte, 4)

	binary.BigEndian.PutUint32(tmpBytes, minutes)

	copy(idBytes, tmpBytes)

	binary.BigEndian.PutUint32(tmpBytes, uint32(me.PackageID))

	for i, b := range tmpBytes {
		idBytes[i+4] = b
	}

	checksum := crc32.ChecksumIEEE([]byte(strings.TrimSpace(strings.ToLower(fmt.Sprintf("%s|%s", me.EventText, me.Details)))))
	binary.BigEndian.PutUint32(tmpBytes, checksum)

	for i, b := range tmpBytes {
		idBytes[i+8] = b
	}

	if me.Location != nil {
		checksum = crc32.ChecksumIEEE([]byte(strings.TrimSpace(strings.ToLower(me.Location.String()))))
		binary.BigEndian.PutUint32(tmpBytes, checksum)

		for i, b := range tmpBytes {
			idBytes[i+12] = b
		}
	}

	me.id, _ = uuid.FromBytes(idBytes)
}

func (me Event) Insert(db *sql.Tx) error {

	stmt, err := db.Prepare("INSERT INTO events (id, package_id, event_text, details, location, event_timestamp, is_current) VALUES (?,?,?,?,?,?,?)")

	if err != nil {
		return err
	}

	_, err = stmt.Exec(me.GetID().String(), me.PackageID, me.EventText, me.Details, me.getLocationString(), me.getTimestampString(), me.getIsCurrentFlag())
	return err
}

func (me Event) Update(db *sql.Tx) error {

	stmt, err := db.Prepare("UPDATE events SET is_current = ? WHERE id = ?")

	if err != nil {
		return err
	}

	_, err = stmt.Exec(me.getIsCurrentFlag(), me.GetID().String())
	return err
}

func (me Event) getIsCurrentFlag() int {

	if me.IsCurrent {
		return 1
	}
	return 0
}

func (me Event) getLocationString() string {

	if me.Location == nil {
		return ""
	}
	return me.Location.String()
}

func (me Event) getTimestampString() string {

	if me.Timestamp == nil {
		return ""
	}
	return me.Timestamp.Format(eventTimestampDateFormat)
}

// MarshalJSON custom json marshalling interface method
func (me Event) MarshalJSON() ([]byte, error) {

	var tmp struct {
		ID        uuid.UUID  `json:"id"`
		EventText string     `json:"event_text"`
		Location  *Location  `json:"location"`
		Timestamp *time.Time `json:"timestamp"`
	}

	tmp.ID = me.id
	tmp.EventText = me.EventText
	tmp.Location = me.Location
	tmp.Timestamp = me.Timestamp

	return json.Marshal(&tmp)
}

// GetNewEvents compares a new list of events with an old list and returns only events not found in the old list
func GetNewEvents(newList, oldList []Event, pkg *Package) []Event {

	results := make([]Event, 0)

	for _, newEvent := range newList {

		newEvent.PackageID = pkg.id
		newEvent.SetID()

		exists := false

		for _, oldEvent := range oldList {
			if newEvent.GetID().String() == oldEvent.GetID().String() {
				exists = true
				break
			}
		}

		if !exists {
			results = append(results, newEvent)
		}
	}
	return results
}

// InsertEvents inserts multiple events and returns an error on the first error encountered
func InsertEvents(events []Event, pkg *Package, db *sql.Tx) error {

	for _, e := range events {

		e.PackageID = pkg.id

		if e.id == uuid.Nil {
			e.SetID()
		}

		if err := e.Insert(db); err != nil {
			return err
		}
	}
	return nil
}

// GetPackageEvents returns all events for a particular package
func GetPackageEvents(db *sql.Tx, p *Package) ([]Event, error) {

	results := make([]Event, 0)

	stmt, err := db.Prepare("SELECT id, package_id, event_text, details, location, event_timestamp, is_current FROM events WHERE package_id = ? ORDER BY event_timestamp DESC")

	if err != nil {
		return results, err
	}

	rows, err := stmt.Query(p.GetID())

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return results, err
		}
		return nil, err
	}

	for rows.Next() {
		var evt Event

		evt, err = newEventFromSQL(rows)

		if err != nil {
			return results, err
		}

		results = append(results, evt)
	}

	return results, nil
}

func newEventFromSQL(rows *sql.Rows) (Event, error) {

	event := Event{}

	var err error
	var idStr, eventText, details, locationStr, timestampStr string
	var packageId, isCurrentFlag int

	if err = rows.Scan(&idStr, &packageId, &eventText, &details, &locationStr, &timestampStr, &isCurrentFlag); err != nil {
		return event, err
	}

	if event.id, err = uuid.Parse(idStr); err != nil {
		return event, err
	}

	event.PackageID = packageId
	event.EventText = eventText
	event.Details = details

	event.Location = NewLocationFromString(locationStr)

	if event.Timestamp, err = parseEventTimestamp(timestampStr); err != nil {
		return event, err
	}

	event.IsCurrent = parseEventIsCurrentFlag(isCurrentFlag)

	return event, nil
}

func parseEventTimestamp(timestampStr string) (*time.Time, error) {

	timestampStr = strings.TrimSpace(timestampStr)

	if timestampStr == "" {
		return nil, nil
	}
	result, err := time.Parse(eventTimestampDateFormat, timestampStr)
	return &result, err
}

func parseEventIsCurrentFlag(isCurrentFlag int) bool {
	return isCurrentFlag == 1
}
