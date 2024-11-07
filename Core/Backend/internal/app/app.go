package app

import (
	"net/http"

	"github.com/chas3air/Airplanes-Co/Core/Backend/internal/routes"
	"github.com/gorilla/mux"
)

func Run() {
	r := mux.NewRouter()

	r.HandleFunc("/flights/get", routes.GetFlightsHandler).Methods(http.MethodGet)
	r.HandleFunc("/flights/insert", routes.InsertFlightHandler).Methods(http.MethodPost)
	r.HandleFunc("/flights/delete/{id}", routes.DeleteFlightHandler).Methods(http.MethodDelete)

	http.ListenAndServe(":84", r)
}
