package app

import (
	"net/http"

	"github.com/chas3air/Airplanes-Co/Core/DAL_tickets/internal/routes"
	"github.com/gorilla/mux"
)

func Run() {
	r := mux.NewRouter()

	r.HandleFunc("/postgres/tickets", routes.GetTickets).Methods(http.MethodGet)
	r.HandleFunc("/postgres/tickets/{id}", routes.GetTicketById).Methods(http.MethodGet)
	r.HandleFunc("/postgres/tickets", routes.InsertTicket).Methods(http.MethodPost)
	r.HandleFunc("/postgres/tickets", routes.UpdateTicket).Methods(http.MethodPatch)
	r.HandleFunc("/postgres/tickets/{id}", routes.DeleteTicket).Methods(http.MethodDelete)

	http.ListenAndServe(":12002", r)
}
