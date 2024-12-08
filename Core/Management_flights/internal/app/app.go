package app

import (
	"net/http"

	"github.com/chas3air/Airplanes-Co/Core/Management_flights/internal/routes"
	"github.com/gorilla/mux"
)

func Run() {
	r := mux.NewRouter()

	r.HandleFunc("/flights", routes.GetFlightsHandler).Methods(http.MethodGet)
	r.HandleFunc("/flights/{id}", routes.GetFlightByIdHandler).Methods(http.MethodGet)
	r.HandleFunc("/flights", routes.InsertFlightsHandler).Methods(http.MethodPost)
	r.HandleFunc("/flights", routes.UpdateFlightHandler).Methods(http.MethodPatch)
	r.HandleFunc("/flights/{id}", routes.DeleteFlightHandler).Methods(http.MethodDelete)

	http.ListenAndServe(":12006", r)
}
