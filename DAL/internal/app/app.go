package app

import (
	"net/http"

	"github.com/chas3air/Airplanes-Co/DAL/internal/routes"
	"github.com/gorilla/mux"
)

func Run() {
	router := mux.NewRouter()

	router.HandleFunc("/postgres/customers/get", routes.GetCustomers).Methods(http.MethodGet)
	router.HandleFunc("/postgres/customers/get/{id:[0-9]+}", routes.GetCustomerById).Methods(http.MethodGet)
	router.HandleFunc("/postgres/customers/insert", routes.InsertCustomer).Methods(http.MethodPut)
	router.HandleFunc("/postgres/customers/update", routes.UpdateCustomer).Methods(http.MethodPatch)
	router.HandleFunc("/postgres/customers/delete", routes.DeleteCustomer).Methods(http.MethodDelete)

	router.HandleFunc("/postgres/flights/get", routes.GetFlights).Methods(http.MethodGet)
	router.HandleFunc("/postgres/flights/get/{id:[0-9]+}", routes.GetFlightById).Methods(http.MethodGet)
	router.HandleFunc("/postgres/flights/insert", routes.InsertFlight).Methods(http.MethodPut)
	router.HandleFunc("/postgres/flights/update", routes.UpdateFlight).Methods(http.MethodPatch)
	router.HandleFunc("/postgres/flights/delete", routes.DeleteFlight).Methods(http.MethodDelete)

	router.HandleFunc("/postgres/tickets/get", routes.GetTickets).Methods(http.MethodGet)
	router.HandleFunc("/postgres/tickets/get/{id:[0-9]+}", routes.GetTicketById).Methods(http.MethodGet)
	router.HandleFunc("/postgres/ticket/insert", routes.InsertTicket).Methods(http.MethodPut)
	router.HandleFunc("/postgres/ticket/update", routes.UpdateTicket).Methods(http.MethodPatch)
	router.HandleFunc("/postgres/ticket/delete", routes.DeleteTicket).Methods(http.MethodDelete)

	http.ListenAndServe(":8056", router)
}
