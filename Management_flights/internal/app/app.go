package app

import (
	"net/http"

	"github.com/chas3air/Airplanes-Co/Management_flights/internal/router/management"
	"github.com/gorilla/mux"
)

func Run() {
	router := mux.NewRouter()

	router.HandleFunc("/api/flight/getAllFlights", management.GetAllFlightsHandler).Methods(http.MethodGet)
	router.HandleFunc("/api/flight/getFlightById/{id:[0-9]+}", management.GetFlightByIdHandler).Methods(http.MethodGet)
	router.HandleFunc("/api/flight/insert", management.InsertFlightsHandler).Methods(http.MethodPost)
	router.HandleFunc("/api/flight/update", management.UpdateFlightHandler).Methods(http.MethodPatch)
	router.HandleFunc("/api/flight/delete/{id:[0-9]+}", management.DeleteFlightHandler).Methods(http.MethodDelete)

	http.ListenAndServe(":12006", router)
}
