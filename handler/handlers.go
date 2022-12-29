package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"regexp"
	"strings"

	"github.com/gorilla/mux"
	"github.com/littlehawk93/columba/config"
	"github.com/littlehawk93/columba/service"
	"github.com/littlehawk93/columba/service/usps"
)

var configuration config.ApplicationConfiguration

var serviceProviders = map[string]service.Provider{
	"usps": usps.Provider{},
}

// SetConfiguration sets the configuration for all handlers
func SetConfiguration(conf config.ApplicationConfiguration) {
	configuration = conf
}

// AddHandlers binds all the API handlers for the web app to their endpoints
func AddHandlers(r *mux.Router) {

	api := r.PathPrefix("/api").Subrouter()

	services := api.PathPrefix("/service").Subrouter()

	services.HandleFunc("/", getAllServices).Methods("GET")

	packages := api.PathPrefix("/package").Subrouter()

	packages.HandleFunc("/", getAllPackages).Methods("GET")
	packages.HandleFunc("/", createPackage).Methods("POST").Headers("Content-Type", "application/json")
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

func getAllServiceProviderNames() []string {

	names := make([]string, len(serviceProviders))

	for i := 0; i < len(names); i++ {
		names[i] = ""
	}

	for _, v := range serviceProviders {

		name := v.GetID()

		for i := 0; i < len(names); i++ {
			if names[i] == "" {
				names[i] = name
				break
			} else if names[i] > name {
				tmp := names[i]
				names[i] = name
				name = tmp
			}
		}
	}

	return names
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

func getServiceProvider(serviceID string) service.Provider {

	serviceID = strings.ToLower(regexp.MustCompile(`[^a-zA-z]`).ReplaceAllString(serviceID, ""))

	if provider, ok := serviceProviders[serviceID]; !ok {
		panic("Invalid service ID provided")
	} else {
		return provider
	}
}
