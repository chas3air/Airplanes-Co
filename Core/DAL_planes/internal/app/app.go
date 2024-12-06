package app

import (
	"net/http"

	"github.com/chas3air/Airplanes-Co/Core/DAL_planes/internal/routes"
	"github.com/gorilla/mux"
)

func Run() {
	r := mux.NewRouter()

	r.HandleFunc("/planes", routes.GetPlanesHandler).Methods(http.MethodGet)
	r.HandleFunc("/planes/{id}", routes.GetPlanesByIdHandler).Methods(http.MethodGet)
	r.HandleFunc("/planes", routes.InsertPlaneHandler).Methods(http.MethodPost)
	r.HandleFunc("/planes", routes.UpdatePlaneHandler).Methods(http.MethodPatch)
	r.HandleFunc("/planes/{id}", routes.DeletePlaneHandler).Methods(http.MethodDelete)

	http.ListenAndServe(":12010", r)
}
