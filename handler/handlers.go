package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/littlehawk93/columba/config"
)

var configuration config.ApplicationConfiguration

// SetConfiguration sets the configuration for all handlers
func SetConfiguration(conf config.ApplicationConfiguration) {
	configuration = conf
}

// AddAPIHandlers binds all the API handlers for the web app to their endpoints
func AddAPIHandlers(r *mux.Router) {

	api := r.PathPrefix("/api").Subrouter()

	services := api.PathPrefix("/service").Subrouter()

	services.HandleFunc("", getAllServices).Methods("GET")

	packages := api.PathPrefix("/package").Subrouter()

	packages.HandleFunc("", getAllPackages).Methods("GET")
	packages.HandleFunc("", createPackage).Methods("POST").Headers("Content-Type", "application/json")
	packages.HandleFunc("/{id:[0-9]+}", deletePackage).Methods("DELETE")

	events := api.PathPrefix("/event").Subrouter()

	events.HandleFunc("/{id:[0-9]+}", getPackageEvents).Methods("GET")
}

func openTx() (*sql.Tx, *sql.DB, error) {

	db, err := configuration.Database.Open()

	if err != nil {
		return nil, nil, err
	}
	tx, err := db.Begin()

	if err != nil {
		db.Close()
		return nil, nil, err
	}
	return tx, db, nil
}

func writeJSON(w http.ResponseWriter, data interface{}, status int) {

	d, err := json.Marshal(&data)

	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(d)
}

func writeError(w http.ResponseWriter, err error, status int) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(status)
	w.Write([]byte(err.Error()))
}
