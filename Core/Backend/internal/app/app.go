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

	r.HandleFunc("/flights", flights_routes.GetFlightsHandler).Methods(http.MethodGet)
	r.HandleFunc("/flights/{id}", flights_routes.GetFlightByIdHandler).Methods(http.MethodGet)
	r.HandleFunc("/flights", flights_routes.InsertFlightHandler).Methods(http.MethodPost)
	r.HandleFunc("/flights", flights_routes.UpdateFlightHandler).Methods(http.MethodPatch)
	r.HandleFunc("/flights/{id}", flights_routes.DeleteFlightHandler).Methods(http.MethodDelete)

	r.HandleFunc("/catalog/flights", flights_routes.GetFlightsHandler).Methods(http.MethodGet)

	r.HandleFunc("/customers", customers_routes.GetCustomersHandler).Methods(http.MethodGet)
	r.HandleFunc("/customers/{id}", customers_routes.GetCustomerByIdHandler).Methods(http.MethodGet)
	r.HandleFunc("/customers", customers_routes.InsertCustomerHandler).Methods(http.MethodPost)
	r.HandleFunc("/customers", customers_routes.UpdateCustomerHandler).Methods(http.MethodPatch)
	r.HandleFunc("/customers/{id}", customers_routes.DeleteCustomerHandler).Methods(http.MethodDelete)

	r.HandleFunc("/tickets", tickets_routes.GetTicketsHandler).Methods(http.MethodGet)
	r.HandleFunc("/tickets/{id}", tickets_routes.GetTicketByIdHandler).Methods(http.MethodGet)
	r.HandleFunc("/tickets", tickets_routes.InsertTicketHandler).Methods(http.MethodPost)
	r.HandleFunc("/tickets", tickets_routes.UpdateTicketHandler).Methods(http.MethodPatch)
	r.HandleFunc("/tickets/{id}", tickets_routes.DeleteTicketHandler).Methods(http.MethodDelete)

	r.HandleFunc("/cart", tickets_routes.GetTicketsFromCartHandler).Methods(http.MethodGet)
	r.HandleFunc("/cart", tickets_routes.InsertTicketToCartHandler).Methods(http.MethodPost)
	r.HandleFunc("/cart/clear/{id}", tickets_routes.ClearTicketCartHandler).Methods(http.MethodDelete)

	r.HandleFunc("/purchasedTickets", tickets_routes.GetPurchasedTicketHandler).Methods(http.MethodGet)

	r.HandleFunc("/payment", tickets_routes.PayForTickets).Methods(http.MethodPost)

	r.HandleFunc("/sign-up", customers_routes.SignUpHandler).Methods(http.MethodPost)
	r.HandleFunc("/sign-in", customers_routes.SignInHandler).Methods(http.MethodGet)

	http.ListenAndServe(":12013", r)
}
