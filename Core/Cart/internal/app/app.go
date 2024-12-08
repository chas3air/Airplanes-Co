package app

import (
	"net/http"

	"github.com/chas3air/Airplanes-Co/Core/Cart/internal/router"
	"github.com/gorilla/mux"
)

func Run() {
	r := mux.NewRouter()

	r.HandleFunc("/cart/{id}", router.GetTicketsHandler).Methods(http.MethodGet)
	r.HandleFunc("/cart", router.InsertTicketHandler).Methods(http.MethodPost)
	r.HandleFunc("/cart", router.UpdateTicketHandler).Methods(http.MethodPatch)
	r.HandleFunc("/cart/{id}", router.DeleteTicketHandler).Methods(http.MethodDelete)
	r.HandleFunc("/cart/clear/{id}", router.ClearHandler).Methods(http.MethodDelete)

	http.ListenAndServe(":12003", r)
}
