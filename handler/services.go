package handler

import (
	"net/http"
)

func getAllServices(w http.ResponseWriter, r *http.Request) {

	writeJSON(w, getAllServiceProviderNames(), http.StatusOK)
}
