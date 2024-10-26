package app

import (
	"net/http"

	"github.com/chas3air/Airplanes-Co/Management_tickets/internal/router"
	"github.com/gorilla/mux"
)

func Run() {
	r := mux.NewRouter()

	r.HandleFunc("/management-tickets/tickets", router.GetTicketsHandler).Methods(http.MethodGet)
	r.HandleFunc("/management-tickets/tickets/{id}", router.GetTicketByIdHandler).Methods(http.MethodGet)
	r.HandleFunc("/management-tickets/tickets", router.InsertTicketHandler).Methods(http.MethodPost)
	r.HandleFunc("/management-tickets/tickets", router.UpdateTicketHandler).Methods(http.MethodPatch)
	r.HandleFunc("/management-tickets/tickets/{id}", router.DeleteTicketHandler).Methods(http.MethodDelete)

	http.ListenAndServe(":12008", r)
}
