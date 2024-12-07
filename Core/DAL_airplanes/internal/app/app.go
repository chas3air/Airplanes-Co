package app

import (
	"net/http"

	ded "github.com/chas3air/Airplanes-Co/Core/DAL_airplanes/internal/router"
	"github.com/gorilla/mux"
)

func Run() {
	r := mux.NewRouter()

	r.HandleFunc("/airplanes", ded.GetAirplanesHandler).Methods(http.MethodGet)
	r.HandleFunc("/airplanes", ded.InsertAirplaneHandler).Methods(http.MethodPost)
	r.HandleFunc("/airplanes", ded.UpdateAirplaneHandler).Methods(http.MethodPatch)
	r.HandleFunc("/airplanes/{id}", ded.DeleteAirplaneHandler).Methods(http.MethodDelete)

	http.ListenAndServe(":12010", r)
}
