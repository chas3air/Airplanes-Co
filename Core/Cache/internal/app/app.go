package app

import (
	"net/http"

	"github.com/chas3air/Airplanes-Co/Core/Cache/internal/routes"
	"github.com/gorilla/mux"
)

func Run() {
	r := mux.NewRouter()
	r.HandleFunc("/cache/{key}", routes.GetItemHandler).Methods(http.MethodGet)
	r.HandleFunc("/cache", routes.SetItemHandler).Methods(http.MethodPost)
	r.HandleFunc("/cache{key}", routes.DeleteItemHandler).Methods(http.MethodDelete)

	http.ListenAndServe(":12011", r)
}
