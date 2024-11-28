package app

import (
	"net/http"

	"github.com/chas3air/Airplanes-Co/Core/DAL_customers/internal/routes"
	"github.com/gorilla/mux"
)

//TODO: изменить энд поинты

func Run() {
	r := mux.NewRouter()

	r.HandleFunc("/postgres/customers", routes.GetCustomers).Methods(http.MethodGet)
	r.HandleFunc("/postgres/customers/{id}/", routes.GetCustomerById).Methods(http.MethodGet)
	r.HandleFunc("/postgres/customers/login", routes.GetCustomerByLoginAndPassword).Methods(http.MethodGet)
	r.HandleFunc("/postgres/customers", routes.InsertCustomer).Methods(http.MethodPost)
	r.HandleFunc("/postgres/customers", routes.UpdateCustomer).Methods(http.MethodPatch)
	r.HandleFunc("/postgres/customers/k", routes.DeleteCustomer).Methods(http.MethodDelete)

	http.ListenAndServe(":12000", r)
}
