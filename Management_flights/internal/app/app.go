package app

import (
	"net/http"

	"github.com/chas3air/Airplanes-Co/Management_flights/internal/router"
	"github.com/gorilla/mux"
)

func Run() {
	r := mux.NewRouter()

	r.HandleFunc("/management-flights/flights", router.GetFlightsHandler).Methods(http.MethodGet)
	r.HandleFunc("/management-flights/flights/{id}", router.GetFlightByIdHandler).Methods(http.MethodGet)
	r.HandleFunc("/management-flights/flights", router.InsertFlightsHandler).Methods(http.MethodPost)
	r.HandleFunc("/management-flights/flights", router.UpdateFlightHandler).Methods(http.MethodPatch)
	r.HandleFunc("/management-flights/flights/{id}", router.DeleteFlightHandler).Methods(http.MethodDelete)

	http.ListenAndServe(":12006", r)
}
