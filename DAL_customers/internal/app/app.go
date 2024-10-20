package app

import (
	"net/http"

	"github.com/chas3air/Airplanes-Co/DAL_customers/internal/routes/customers"
	"github.com/gorilla/mux"
)

func Run() {
	router := mux.NewRouter()

	router.HandleFunc("/dal/customers/postgres/get", customers.GetCustomers).Methods(http.MethodGet)
	router.HandleFunc("/dal/customers/postgres/get/{id:[0-9]+}", customers.GetCustomerById).Methods(http.MethodGet)
	router.HandleFunc("/dal/customers/postgres/get/lp", customers.GetCustomerByLoginAndPassword).Methods(http.MethodGet)
	router.HandleFunc("/dal/customers/postgres/insert", customers.InsertCustomer).Methods(http.MethodPost)
	router.HandleFunc("/dal/customers/postgres/update", customers.UpdateCustomer).Methods(http.MethodPatch)
	router.HandleFunc("/dal/customers/postgres/delete/{id:[0-9]+}", customers.DeleteCustomer).Methods(http.MethodDelete)

	http.ListenAndServe(":12000", router)
}
