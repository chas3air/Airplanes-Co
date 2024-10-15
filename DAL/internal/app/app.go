package app

import (
	"net/http"

	"github.com/chas3air/Airplanes-Co/DAL/internal/routes/customers"
	"github.com/chas3air/Airplanes-Co/DAL/internal/routes/flights"
	"github.com/gorilla/mux"
)

func Run() {
	router := mux.NewRouter()

	router.HandleFunc("/postgres/customers/get", customers.GetCustomers).Methods(http.MethodGet)
	router.HandleFunc("/postgres/customers/get/{id:[0-9]+}", customers.GetCustomerById).Methods(http.MethodGet)
	router.HandleFunc("/postgres/customers/insert", customers.InsertCustomer).Methods(http.MethodPut)
	router.HandleFunc("/postgres/customers/update", customers.UpdateCustomer).Methods(http.MethodPatch)
	router.HandleFunc("/postgres/customers/delete/{id:[0-9]+}", customers.DeleteCustomer).Methods(http.MethodDelete)

	router.HandleFunc("/postgres/flights/get", flights.GetFlights).Methods(http.MethodGet)
	router.HandleFunc("/postgres/flights/get/{id:[0-9]+}", flights.GetFlightById).Methods(http.MethodGet)
	router.HandleFunc("/postgres/flights/insert", flights.InsertFlight).Methods(http.MethodPut)
	router.HandleFunc("/postgres/flights/update", flights.UpdateFlight).Methods(http.MethodPatch)
	router.HandleFunc("/postgres/flights/delete/{id:[0-9]+}", flights.DeleteFlight).Methods(http.MethodDelete)

	// router.HandleFunc("/postgres/tickets/get", routes.GetTickets).Methods(http.MethodGet)
	// router.HandleFunc("/postgres/tickets/get/{id:[0-9]+}", routes.GetTicketById).Methods(http.MethodGet)
	// router.HandleFunc("/postgres/ticket/insert", routes.InsertTicket).Methods(http.MethodPut)
	// router.HandleFunc("/postgres/ticket/update", routes.UpdateTicket).Methods(http.MethodPatch)
	// router.HandleFunc("/postgres/ticket/delete", routes.DeleteTicket).Methods(http.MethodDelete)

	http.ListenAndServe(":8056", router)
}
