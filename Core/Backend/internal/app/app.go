package app

import (
	"net/http"

	"github.com/chas3air/Airplanes-Co/Core/Backend/internal/routes/flights_routes"
	"github.com/gorilla/mux"
)

func Run() {
	r := mux.NewRouter()

	r.HandleFunc("/flights/get", flights_routes.GetFlightsHandler).Methods(http.MethodGet)
	r.HandleFunc("/flights/get/{id}", flights_routes.GetFlightByIdHandler).Methods(http.MethodGet)
	r.HandleFunc("/flights/insert", flights_routes.InsertFlightHandler).Methods(http.MethodPost)
	r.HandleFunc("/flights/update", flights_routes.UpdateFlightHandler).Methods(http.MethodPatch)
	r.HandleFunc("/flights/delete/{id}", flights_routes.DeleteFlightHandler).Methods(http.MethodDelete)

	r.HandleFunc("/catalog/flights", flights_routes.GetFlightsHandler).Methods(http.MethodGet)

	http.ListenAndServe(":12013", r)
}
