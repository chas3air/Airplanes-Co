package app

import (
	"net/http"

	"github.com/chas3air/Airplanes-Co/DAL_flights/internal/router/flights"
	"github.com/gorilla/mux"
)

func Run() {
	router := mux.NewRouter()

	router.HandleFunc("/postgres/flights", flights.GetFlights).Methods(http.MethodGet)
	router.HandleFunc("/postgres/flights/{id}", flights.GetFlightById).Methods(http.MethodGet)
	router.HandleFunc("/postgres/flights", flights.InsertFlight).Methods(http.MethodPost)
	router.HandleFunc("/postgres/flights", flights.UpdateFlight).Methods(http.MethodPatch)
	router.HandleFunc("/postgres/flights/{id}", flights.DeleteFlight).Methods(http.MethodDelete)

	http.ListenAndServe(":12001", router)
}
