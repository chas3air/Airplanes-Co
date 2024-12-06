package app

import (
	"net/http"

	"github.com/chas3air/Airplanes-Co/Core/Payment/internal/router"
	"github.com/gorilla/mux"
)

func Run() {
	r := mux.NewRouter()

	r.HandleFunc("/pay", router.PayForTicket).Methods(http.MethodPost)

	http.ListenAndServe(":12010", r)
}
