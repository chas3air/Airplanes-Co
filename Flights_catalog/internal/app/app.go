package app

import (
	"net/http"

	"github.com/chas3air/Airplanes-Co/Flights_catalog/internal/router/catalog"
	"github.com/gorilla/mux"
)

func Run() {
	router := mux.NewRouter()

	router.HandleFunc("/api/catalog/flights", catalog.GetFlightsHandler).Methods(http.MethodGet)
	router.HandleFunc("/api/catalog/flights/{id:[0-9]+}", catalog.GetFlightByIDHandler).Methods(http.MethodGet)
	router.HandleFunc("/api/catalog/flights/search", catalog.SearchFlightsHandler).Methods(http.MethodGet)

	router.HandleFunc("/api/catalog/flights/stats", catalog.StatsFlightsHandler).Methods(http.MethodGet)

	http.ListenAndServe(":12004", router)
}
