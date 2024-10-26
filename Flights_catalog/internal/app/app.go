package app

import (
	"net/http"

	"github.com/chas3air/Airplanes-Co/Flights_catalog/internal/router/catalog"
	"github.com/gorilla/mux"
)

func Run() {
	router := mux.NewRouter()

	router.HandleFunc("/catalog/flights", catalog.GetFlightsHandler).Methods(http.MethodGet)
	router.HandleFunc("/catalog/flights/{id}", catalog.GetFlightByIDHandler).Methods(http.MethodGet)
	router.HandleFunc("/catalog/flights/search", catalog.SearchFlightsHandler).Methods(http.MethodGet)

	router.HandleFunc("/catalog/flights/stats", catalog.StatsFlightsHandler).Methods(http.MethodGet)

	http.ListenAndServe(":12004", router)
}
