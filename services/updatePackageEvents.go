package services

import (
	"fmt"

	"github.com/littlehawk93/columba/config"
	"github.com/littlehawk93/columba/providers"
	"github.com/littlehawk93/columba/tracking"
)

// UpdatePackageEvents get all active packages and update their events
func UpdatePackageEvents(cfg config.ApplicationConfiguration) error {

	db, err := cfg.Database.Open()

	if err != nil {
		return err
	}

	defer db.Close()

	tx, err := db.Begin()

	if err != nil {
		return err
	}

	packages, err := tracking.GetAllPackagesWithEvents(tracking.PackageStatusActive, tx)

	if err != nil {
		tx.Rollback()
		return err
	}

	for _, pkg := range packages {

		prov := providers.GetServiceProvider(pkg.ServiceID)

		if prov == nil {
			tx.Rollback()
			return fmt.Errorf("unable to get service provider for id '%s'", pkg.ServiceID)
		}

		events, err := prov.GetTrackingEvents(pkg.TrackingNumber)

		if err != nil {
			tx.Rollback()
			return err
		}

		newEvents := tracking.GetNewEvents(events, pkg.Events, &pkg)

		if len(newEvents) > 0 {
			if err = tracking.InsertEvents(newEvents, &pkg, tx); err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
	}
	return err
}
