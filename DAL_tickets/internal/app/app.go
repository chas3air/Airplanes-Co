package app

import (
	"net/http"

	"github.com/chas3air/Airplanes-Co/DAL_tickets/internal/routes/tickets"
	"github.com/gorilla/mux"
)

func Run() {
	router := mux.NewRouter()

	router.HandleFunc("/dal/tickets/postgres/get", tickets.GetTickets).Methods(http.MethodGet)
	router.HandleFunc("/dal/tickets/postgres/get/{id:[0-9]+}", tickets.GetTicketById).Methods(http.MethodGet)
	router.HandleFunc("/dal/tickets/postgres/insert", tickets.InsertTicket).Methods(http.MethodPost)
	router.HandleFunc("/dal/tickets/postgres/update", tickets.UpdateTicket).Methods(http.MethodPatch)
	router.HandleFunc("/dal/tickets/postgres/delete/{id:[0-9]+}", tickets.DeleteTicket).Methods(http.MethodDelete)

	http.ListenAndServe(":12002", router)
}
