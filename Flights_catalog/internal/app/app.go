package app

import (
	"net/http"

	"github.com/chas3air/Airplanes-Co/Flights_catalog/internal/router"
	"github.com/gorilla/mux"
)

func Run() {
	r := mux.NewRouter()

	r.HandleFunc("/catalog/flights", router.GetFlightsHandler).Methods(http.MethodGet)
	r.HandleFunc("/catalog/flights/{id}", router.GetFlightByIDHandler).Methods(http.MethodGet)
	r.HandleFunc("/catalog/flights/search", router.SearchFlightsHandler).Methods(http.MethodGet)

	r.HandleFunc("/catalog/flights/stats", router.StatsFlightsHandler).Methods(http.MethodGet)

	http.ListenAndServe(":12004", r)
}
