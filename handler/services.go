package handler

import (
	"net/http"

	"github.com/littlehawk93/columba/providers"
)

func getAllServices(w http.ResponseWriter, r *http.Request) {

	writeJSON(w, providers.GetAllServiceProviderNames(), http.StatusOK)
}
