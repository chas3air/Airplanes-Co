package app

import (
	"net/http"

	"github.com/chas3air/Airplanes-Co/Core/Management_airplanes/internal/routes"
	"github.com/gorilla/mux"
)

func Run() {
	r := mux.NewRouter()

	r.HandleFunc("/airplanes", routes.GetAirplanesHandler).Methods(http.MethodGet)
	r.HandleFunc("/airplanes/{id}", routes.GetAirplaneByIdHandler).Methods(http.MethodGet)
	r.HandleFunc("/airplanes", routes.InsertAirplaneHandler).Methods(http.MethodPost)
	r.HandleFunc("/airplanes", routes.UpdateAirplaneHandler).Methods(http.MethodPatch)
	r.HandleFunc("/airplanes/{id}", routes.DeleteAirplaneHandler).Methods(http.MethodDelete)

	http.ListenAndServe(":12008", r)
}
