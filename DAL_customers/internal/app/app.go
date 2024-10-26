package app

import (
	"net/http"

	"github.com/chas3air/Airplanes-Co/DAL_customers/internal/routes/customers"
	"github.com/gorilla/mux"
)

func Run() {
	router := mux.NewRouter()

	router.HandleFunc("/postgres/customers", customers.GetCustomers).Methods(http.MethodGet)
	router.HandleFunc("/postgres/customers/{id}", customers.GetCustomerById).Methods(http.MethodGet)
	router.HandleFunc("/postgres/customers/login", customers.GetCustomerByLoginAndPassword).Methods(http.MethodGet)
	router.HandleFunc("/postgres/customers", customers.InsertCustomer).Methods(http.MethodPost)
	router.HandleFunc("/postgres/customers", customers.UpdateCustomer).Methods(http.MethodPatch)
	router.HandleFunc("/postgres/customers/{id}", customers.DeleteCustomer).Methods(http.MethodDelete)

	http.ListenAndServe(":12000", router)
}
