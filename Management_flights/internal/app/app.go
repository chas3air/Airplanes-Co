package app

import (
	"net/http"

	"github.com/chas3air/Airplanes-Co/Management_flights/internal/router"
	"github.com/gorilla/mux"
)

func Run() {
	r := mux.NewRouter()

	r.HandleFunc("/api/flight/getAllFlights", router.GetAllFlightsHandler).Methods(http.MethodGet)
	r.HandleFunc("/api/flight/getFlightById/{id:[0-9]+}", router.GetFlightByIdHandler).Methods(http.MethodGet)
	r.HandleFunc("/api/flight/insert", router.InsertFlightsHandler).Methods(http.MethodPost)
	r.HandleFunc("/api/flight/update", router.UpdateFlightHandler).Methods(http.MethodPatch)
	r.HandleFunc("/api/flight/delete/{id:[0-9]+}", router.DeleteFlightHandler).Methods(http.MethodDelete)

	http.ListenAndServe(":12006", r)
}
