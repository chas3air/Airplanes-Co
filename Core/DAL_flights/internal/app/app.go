package app

import (
	"net/http"

	"github.com/chas3air/Airplanes-Co/Core/DAL_flights/internal/routes"
	"github.com/gorilla/mux"
)

func Run() {
	r := mux.NewRouter()

	r.HandleFunc("/postgres/flights", routes.GetFlights).Methods(http.MethodGet)
	r.HandleFunc("/postgres/flights/{id}", routes.GetFlightById).Methods(http.MethodGet)
	r.HandleFunc("/postgres/flights", routes.InsertFlight).Methods(http.MethodPost)
	r.HandleFunc("/postgres/flights", routes.UpdateFlight).Methods(http.MethodPatch)
	r.HandleFunc("/postgres/flights/{id}", routes.DeleteFlight).Methods(http.MethodDelete)

	http.ListenAndServe(":12001", r)
}
