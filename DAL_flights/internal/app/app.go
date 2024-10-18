package app

import (
	"net/http"

	"github.com/chas3air/Airplanes-Co/DAL_flights/internal/routes/flights"
	"github.com/gorilla/mux"
)

func Run() {
	router := mux.NewRouter()

	router.HandleFunc("/dal/flights/postgres/get", flights.GetFlights).Methods(http.MethodGet)
	router.HandleFunc("/dal/flights/postgres/get/{id:[0-9]+}", flights.GetFlightById).Methods(http.MethodGet)
	router.HandleFunc("/dal/flights/postgres/insert", flights.InsertFlight).Methods(http.MethodPost)
	router.HandleFunc("/dal/flights/postgres/update", flights.UpdateFlight).Methods(http.MethodPatch)
	router.HandleFunc("/dal/flights/postgres/delete/{id:[0-9]+}", flights.DeleteFlight).Methods(http.MethodDelete)

	http.ListenAndServe(":12001", router)
}
