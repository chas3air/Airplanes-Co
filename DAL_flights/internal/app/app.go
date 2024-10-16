package app

import (
	"net/http"

	"github.com/chas3air/Airplanes-Co/DAL_flights/internal/routes/flights"
	"github.com/gorilla/mux"
)

func Run() {
	router := mux.NewRouter()

	router.HandleFunc("/postgres/flights/get", flights.GetFlights).Methods(http.MethodGet)
	router.HandleFunc("/postgres/flights/get/{id:[0-9]+}", flights.GetFlightById).Methods(http.MethodGet)
	router.HandleFunc("/postgres/flights/insert", flights.InsertFlight).Methods(http.MethodPut)
	router.HandleFunc("/postgres/flights/update", flights.UpdateFlight).Methods(http.MethodPatch)
	router.HandleFunc("/postgres/flights/delete/{id:[0-9]+}", flights.DeleteFlight).Methods(http.MethodDelete)

	http.ListenAndServe(":8058", router)
}
