package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/littlehawk93/columba/tracking"
)

func getAllPackages(w http.ResponseWriter, r *http.Request) {

	tx, db, err := openTx()

	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	defer func() {
		tx.Commit()
		db.Close()
	}()

	packages, err := tracking.GetAllPackages(tracking.PackageStatusActive, tx)

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

	writeJSON(w, pkg, http.StatusCreated)
}
