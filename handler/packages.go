package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

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

	packages, err := tracking.GetAllPackages(status, tx)

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
