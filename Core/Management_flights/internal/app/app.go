package app

import (
	"net/http"

	"github.com/chas3air/Airplanes-Co/Core/Management_flights/internal/routes"
	"github.com/gorilla/mux"
)

func Run() {
	r := mux.NewRouter()

	r.HandleFunc("/management-flights/flights", routes.GetFlightsHandler).Methods(http.MethodGet)
	r.HandleFunc("/management-flights/flights/{id}", routes.GetFlightByIdHandler).Methods(http.MethodGet)
	r.HandleFunc("/management-flights/flights", routes.InsertFlightsHandler).Methods(http.MethodPost)
	r.HandleFunc("/management-flights/flights", routes.UpdateFlightHandler).Methods(http.MethodPatch)
	r.HandleFunc("/management-flights/flights/{id}", routes.DeleteFlightHandler).Methods(http.MethodDelete)

	http.ListenAndServe(":12006", r)
}
