package tracking

import (
	"database/sql"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

const (
	packageTimestampDateFormat         string = "2006-01-02 15:04:05.999"
	packageTimestampDeliveryDateFormat string = "2006-01-02"
)

// Package defines a single package to track
type Package struct {
	id                    int
	createdOn             time.Time
	lastUpdatedOn         time.Time
	Status                PackageStatus
	TrackingNumber        string
	ServiceID             string
	Provider              Provider
	Label                 string
	Origin                *Location
	Destination           *Location
	EstimatedDeliveryDate *time.Time
	Events                []Event
}

// GetID get this package's ID
func (me Package) GetID() int {
	return me.id
}

// GetCreatedOn get this package's creation date and timestamp
func (me Package) GetCreatedOn() time.Time {
	return me.createdOn
}

// GetLastUpdatedOn get this package's last updated date and timestamp
func (me Package) GetLastUpdatedOn() time.Time {
	return me.lastUpdatedOn
}

// Insert inserts this record into the database, returns any errors
func (me *Package) Insert(db *sql.Tx) error {

	me.createdOn = time.Now()
	me.lastUpdatedOn = time.Now()

	stmt, err := db.Prepare("INSERT INTO packages (status, tracking_number, service, label, created_on, last_updated_on, origin, destination, estimated_delivery_date) VALUES (?,?,?,?,?,?,?,?,?) RETURNING id")

	if err != nil {
		return err
	}

	deliveryDateStr := ""

	if me.EstimatedDeliveryDate != nil {
		deliveryDateStr = me.EstimatedDeliveryDate.Format(packageTimestampDeliveryDateFormat)
	}

	originStr := ""

	if me.Origin != nil {
		originStr = me.Origin.String()
	}

	destinationStr := ""

	if me.Destination != nil {
		destinationStr = me.Destination.String()
	}

	row := stmt.QueryRow(int(me.Status), me.TrackingNumber, me.ServiceID, me.Label, me.createdOn.Format(packageTimestampDateFormat), me.lastUpdatedOn.Format(packageTimestampDateFormat), originStr, destinationStr, deliveryDateStr)

	if row.Err() != nil {
		return row.Err()
	}

	return row.Scan(&(me.id))
}

func (me *Package) UpdateStatus(status PackageStatus, db *sql.Tx) error {

	me.lastUpdatedOn = time.Now()
	me.Status = status

	stmt, err := db.Prepare("UPDATE packages SET status = ?, last_updated_on = ? WHERE id = ?")

	if err != nil {
		return err
	}

	_, err = stmt.Exec(int(me.Status), me.lastUpdatedOn.Format(packageTimestampDateFormat), me.id)
	return err
}

// Update updates information about this record in the database, returns any errors
func (me *Package) Update(db *sql.Tx) error {

	me.lastUpdatedOn = time.Now()

	stmt, err := db.Prepare("UPDATE packages SET label = ?, last_updated_on = ?, origin = ?, destination = ?, estimated_delivery_date = ? WHERE id = ?")

	if err != nil {
		return err
	}

	deliveryDateStr := ""

	if me.EstimatedDeliveryDate != nil {
		deliveryDateStr = me.EstimatedDeliveryDate.Format(packageTimestampDeliveryDateFormat)
	}

	originStr := ""

	if me.Origin != nil {
		originStr = me.Origin.String()
	}

	destinationStr := ""

	if me.Destination != nil {
		destinationStr = me.Destination.String()
	}

	_, err = stmt.Exec(me.Label, me.lastUpdatedOn.Format(packageTimestampDateFormat), originStr, destinationStr, deliveryDateStr, me.id)
	return err
}

// MarshalJSON custom json marshalling interface method
func (me Package) MarshalJSON() ([]byte, error) {

	var tmp struct {
		ID                    int        `json:"id"`
		CreatedOn             time.Time  `json:"created_on"`
		LastUpdatedOn         time.Time  `json:"last_updated_on"`
		Status                string     `json:"status"`
		TrackingNumber        string     `json:"tracking_number"`
		TrackingURL           string     `json:"tracking_url"`
		Service               string     `json:"service"`
		Label                 string     `json:"label"`
		Origin                *Location  `json:"origin"`
		Destination           *Location  `json:"destination"`
		EstimatedDeliveryDate *time.Time `json:"estimated_delivery_date"`
		Events                []Event    `json:"events"`
	}

	tmp.ID = me.id
	tmp.CreatedOn = me.createdOn
	tmp.LastUpdatedOn = me.lastUpdatedOn
	tmp.Status = me.Status.String()
	tmp.TrackingNumber = me.TrackingNumber
	tmp.Service = me.ServiceID
	tmp.Label = me.Label
	tmp.Origin = me.Origin
	tmp.Destination = me.Destination
	tmp.EstimatedDeliveryDate = me.EstimatedDeliveryDate
	tmp.Events = me.Events

	if me.Provider != nil {
		tmp.TrackingURL = me.Provider.GetTrackingURL(me.TrackingNumber)
	}

	return json.Marshal(&tmp)
}

// UnmarshalJSON custom JSON unmarshalling interface method
func (me *Package) UnmarshalJSON(data []byte) error {

	var tmp struct {
		ID             int    `json:"id"`
		Status         string `json:"string"`
		TrackingNumber string `json:"tracking_number"`
		Service        string `json:"service"`
		Label          string `json:"label"`
	}

	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	me.id = tmp.ID
	me.Status = ParsePackageStatus(tmp.Status)
	me.TrackingNumber = tmp.TrackingNumber
	me.ServiceID = tmp.Service
	me.Label = tmp.Label

	return nil
}

// GetPackage retrieve a single package from the database using it's ID. Returns nil if no package was found with the id
func GetPackage(id int, db *sql.Tx) (*Package, error) {

	stmt, err := db.Prepare("SELECT id, status, tracking_number, service, label, created_on, last_updated_on, origin, destination, estimated_delivery_date FROM packages WHERE id = ?")

	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	defer rows.Close()

	var pkg Package

	if rows.Next() {

		if pkg, err = newPackageFromSQL(rows); err != nil {
			return nil, err
		}
	}
	return &pkg, nil
}

// GetAllPackagesWithEvents retrieve all packages that have a particular status along with the events for each returned package
func GetAllPackagesWithEvents(status PackageStatus, db *sql.Tx) ([]Package, error) {

	packages, err := GetAllPackages(status, db)

	if err != nil {
		return nil, err
	}

	for i := range packages {

		events, err := GetPackageEvents(db, &packages[i])

		if err != nil {
			return nil, err
		}

		packages[i].Events = events
	}

	return packages, nil
}

// GetAllPackages retrieve all packages that have a particular status
func GetAllPackages(status PackageStatus, db *sql.Tx) ([]Package, error) {

	stmt, err := db.Prepare("SELECT id, status, tracking_number, service, label, created_on, last_updated_on, origin, destination, estimated_delivery_date FROM packages WHERE status = ?")

	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(int(status))

	results := make([]Package, 0)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return results, nil
		}
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var pkg Package

		if pkg, err = newPackageFromSQL(rows); err != nil {
			return results, err
		}

		results = append(results, pkg)
	}

	return results, nil
}

func newPackageFromSQL(rows *sql.Rows) (Package, error) {

	pkg := Package{}
	var err error

	var id, status int
	var trackingNumber, service, label, createdOnStr, lastUpdatedOnStr, origin, destination, estimatedDeliveryDateStr string

	if err = rows.Scan(&id, &status, &trackingNumber, &service, &label, &createdOnStr, &lastUpdatedOnStr, &origin, &destination, &estimatedDeliveryDateStr); err != nil {
		return pkg, err
	}

	pkg.id = id
	pkg.Status = PackageStatus(status)
	pkg.TrackingNumber = trackingNumber
	pkg.ServiceID = service
	pkg.Label = label

	if pkg.createdOn, err = time.Parse(packageTimestampDateFormat, createdOnStr); err != nil {
		return pkg, err
	}

	if pkg.lastUpdatedOn, err = time.Parse(packageTimestampDateFormat, lastUpdatedOnStr); err != nil {
		return pkg, err
	}

	pkg.Origin = NewLocationFromString(origin)
	pkg.Destination = NewLocationFromString(destination)

	if strings.TrimSpace(estimatedDeliveryDateStr) == "" {
		pkg.EstimatedDeliveryDate = nil
	} else {
		var tmp time.Time

		if tmp, err = time.Parse(packageTimestampDeliveryDateFormat, estimatedDeliveryDateStr); err != nil {
			return pkg, err
		}
		pkg.EstimatedDeliveryDate = &tmp
	}

	return pkg, nil
}
