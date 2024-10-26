package app

import (
	"net/http"

	"github.com/chas3air/Airplanes-Co/DAL_tickets/internal/routes/tickets"
	"github.com/gorilla/mux"
)

func Run() {
	router := mux.NewRouter()

	router.HandleFunc("/postgres/tickets", tickets.GetTickets).Methods(http.MethodGet)
	router.HandleFunc("/postgres/tickets/{id}", tickets.GetTicketById).Methods(http.MethodGet)
	router.HandleFunc("/postgres/tickets", tickets.InsertTicket).Methods(http.MethodPost)
	router.HandleFunc("/postgres/tickets", tickets.UpdateTicket).Methods(http.MethodPatch)
	router.HandleFunc("/postgres/tickets/{id}", tickets.DeleteTicket).Methods(http.MethodDelete)

	http.ListenAndServe(":12002", router)
}
