package app

import (
	"net/http"

	"github.com/chas3air/Airplanes-Co/Core/Management_tickets/internal/routes"
	"github.com/gorilla/mux"
)

func Run() {
	r := mux.NewRouter()

	r.HandleFunc("/management-tickets/tickets", routes.GetTicketsHandler).Methods(http.MethodGet)
	r.HandleFunc("/management-tickets/tickets/{id}", routes.GetTicketByIdHandler).Methods(http.MethodGet)
	r.HandleFunc("/management-tickets/tickets", routes.InsertTicketHandler).Methods(http.MethodPost)
	r.HandleFunc("/management-tickets/tickets", routes.UpdateTicketHandler).Methods(http.MethodPatch)
	r.HandleFunc("/management-tickets/tickets/{id}", routes.DeleteTicketHandler).Methods(http.MethodDelete)

	http.ListenAndServe(":12008", r)
}
