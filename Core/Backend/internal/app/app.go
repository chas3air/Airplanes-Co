package app

import (
	"net/http"

	"github.com/chas3air/Airplanes-Co/Core/Backend/internal/routes/customers_routes"
	"github.com/chas3air/Airplanes-Co/Core/Backend/internal/routes/flights_routes"
	"github.com/chas3air/Airplanes-Co/Core/Backend/internal/routes/tickets_routes"
	"github.com/gorilla/mux"
)

func Run() {
	r := mux.NewRouter()

	r.HandleFunc("/flights/get", flights_routes.GetFlightsHandler).Methods(http.MethodGet)
	r.HandleFunc("/flights/get/{id}", flights_routes.GetFlightByIdHandler).Methods(http.MethodGet)
	r.HandleFunc("/flights/insert", flights_routes.InsertFlightHandler).Methods(http.MethodPost)
	r.HandleFunc("/flights/update", flights_routes.UpdateFlightHandler).Methods(http.MethodPatch)
	r.HandleFunc("/flights/delete/{id}", flights_routes.DeleteFlightHandler).Methods(http.MethodDelete)

	r.HandleFunc("/catalog/flights", flights_routes.GetFlightsHandler).Methods(http.MethodGet)

	r.HandleFunc("/customers/get", customers_routes.GetCustomersHandler).Methods(http.MethodGet)
	r.HandleFunc("/customers/get/{id}", customers_routes.GetCustomerByIdHandler).Methods(http.MethodGet)
	r.HandleFunc("/customers/insert", customers_routes.InsertCustomerHandler).Methods(http.MethodPost)
	r.HandleFunc("/customers/update", customers_routes.UpdateCustomerHandler).Methods(http.MethodPatch)
	r.HandleFunc("/customers/delete/{id}", customers_routes.DeleteCustomerHandler).Methods(http.MethodDelete)

	r.HandleFunc("/ticket/get", tickets_routes.GetTicketHandler).Methods(http.MethodGet)
	r.HandleFunc("/ticket/get/{id}", tickets_routes.GetTicketByIdHandler).Methods(http.MethodGet)
	r.HandleFunc("/ticket/insert", tickets_routes.InsertTicketHandler).Methods(http.MethodPost)
	r.HandleFunc("/ticket/update", tickets_routes.UpdateTicketHandler).Methods(http.MethodPatch)
	r.HandleFunc("/ticket/delete/{id}", tickets_routes.DeleteTicketHandler).Methods(http.MethodDelete)

	r.HandleFunc("/sign-up", customers_routes.SignUpHandler).Methods(http.MethodPost)
	r.HandleFunc("/sign-in", customers_routes.SignInHandler).Methods(http.MethodGet)

	http.ListenAndServe(":12013", r)
}
