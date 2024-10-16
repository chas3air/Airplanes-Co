package app

import (
	"net/http"

	"github.com/chas3air/Airplanes-Co/DAL/internal/routes/flights"
	"github.com/chas3air/Airplanes-Co/DAL/internal/routes/tickets"
	"github.com/gorilla/mux"
)

func Run() {
	router := mux.NewRouter()

	router.HandleFunc("/postgres/flights/get", flights.GetFlights).Methods(http.MethodGet)
	router.HandleFunc("/postgres/flights/get/{id:[0-9]+}", flights.GetFlightById).Methods(http.MethodGet)
	router.HandleFunc("/postgres/flights/insert", flights.InsertFlight).Methods(http.MethodPut)
	router.HandleFunc("/postgres/flights/update", flights.UpdateFlight).Methods(http.MethodPatch)
	router.HandleFunc("/postgres/flights/delete/{id:[0-9]+}", flights.DeleteFlight).Methods(http.MethodDelete)

	router.HandleFunc("/postgres/tickets/get", tickets.GetTickets).Methods(http.MethodGet)
	router.HandleFunc("/postgres/tickets/get/{id:[0-9]+}", tickets.GetTicketById).Methods(http.MethodGet)
	router.HandleFunc("/postgres/ticket/insert", tickets.InsertTicket).Methods(http.MethodPut)
	router.HandleFunc("/postgres/ticket/update", tickets.UpdateTicket).Methods(http.MethodPatch)
	router.HandleFunc("/postgres/ticket/delete", tickets.DeleteTicket).Methods(http.MethodDelete)

	http.ListenAndServe(":8056", router)
}
