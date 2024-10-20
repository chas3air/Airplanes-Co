package app

import (
	"net/http"

	"github.com/chas3air/Airplanes-Co/Management_tickets/internal/router"
	"github.com/gorilla/mux"
)

func Run() {
	r := mux.NewRouter()

	r.HandleFunc("/api/tickets/getAllTickets", router.GetAllTicketsHandler).Methods(http.MethodGet)
	r.HandleFunc("/api/tickets/getTicketById/{id:[0-9]+}", router.GetTicketByIdHandler).Methods(http.MethodGet)
	r.HandleFunc("/api/tickets/insert", router.InsertTicketHandler).Methods(http.MethodPost)
	r.HandleFunc("/api/tickets/update", router.UpdateTicketHandler).Methods(http.MethodPatch)
	r.HandleFunc("/api/tickets/delete", router.DeleteTicketHandler).Methods(http.MethodDelete)

	http.ListenAndServe(":12008", r)
}
