package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/littlehawk93/columba/providers"
	"github.com/littlehawk93/columba/tracking"
)

func getPackageEvents(w http.ResponseWriter, r *http.Request) {

	packageIDStr, ok := mux.Vars(r)["id"]

	if !ok {
		writeError(w, errors.New("must provide package ID"), http.StatusBadRequest)
		return
	}

	packageID, err := strconv.ParseInt(packageIDStr, 10, 32)

	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	tx, db, err := openTx()

	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	defer func() {
		tx.Commit()
		db.Close()
	}()

	pkg, err := tracking.GetPackage(int(packageID), tx)

	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	if pkg == nil {
		writeError(w, errors.New("package not found"), http.StatusNotFound)
		return
	}

	dbEvents, err := tracking.GetPackageEvents(tx, pkg)

	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	provider := providers.GetServiceProvider(pkg.ServiceID)

	if provider == nil {
		writeJSON(w, dbEvents, http.StatusOK)
		return
	}

	newEvents, err := provider.GetTrackingEvents(pkg.TrackingNumber)

	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	newEvents = tracking.GetNewEvents(newEvents, dbEvents)

	if len(newEvents) > 0 {
		if err = tracking.InsertEvents(newEvents, tx); err != nil {
			tx.Rollback()
			writeError(w, err, http.StatusInternalServerError)
			return
		}

		if err = pkg.Update(tx); err != nil {
			tx.Rollback()
			writeError(w, err, http.StatusInternalServerError)
			return
		}
	}

	totalEvents := append(newEvents, dbEvents...)
	writeJSON(w, totalEvents, http.StatusOK)
}
