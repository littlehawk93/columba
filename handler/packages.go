package handler

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/littlehawk93/columba/providers"
	"github.com/littlehawk93/columba/tracking"
)

func getAllPackages(w http.ResponseWriter, r *http.Request) {

	status := tracking.PackageStatusActive

	vals := r.URL.Query()

	if s := vals.Get("status"); s != "" {
		if sVal := tracking.ParsePackageStatus(s); sVal != tracking.PackageStatusError {
			status = sVal
		}
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

	packages, err := tracking.GetAllPackagesWithEvents(status, tx)

	for i := range packages {
		packages[i].Provider = providers.GetServiceProvider(packages[i].ServiceID)
	}

	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	writeJSON(w, packages, http.StatusOK)
}

func createPackage(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	bytes, err := ioutil.ReadAll(r.Body)

	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	var pkg tracking.Package

	if err := json.Unmarshal(bytes, &pkg); err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	tx, db, err := openTx()

	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	defer func() {
		db.Close()
	}()

	pkg.Status = tracking.PackageStatusActive

	if err = pkg.Insert(tx); err != nil {
		tx.Rollback()
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	pkg.Provider = providers.GetServiceProvider(pkg.ServiceID)
	writeJSON(w, pkg, http.StatusCreated)
}

func deletePackage(w http.ResponseWriter, r *http.Request) {

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
		db.Close()
	}()

	pkg, err := tracking.GetPackage(int(packageID), tx)

	if err != nil {
		tx.Rollback()
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	if pkg == nil {
		tx.Rollback()
		writeError(w, errors.New("package not found"), http.StatusNotFound)
		return
	}

	if err = pkg.UpdateStatus(tracking.PackageStatusArchived, tx); err != nil {
		tx.Rollback()
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte{})
}
