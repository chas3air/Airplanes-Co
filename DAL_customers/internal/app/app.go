package app

import (
	"net/http"

	"github.com/chas3air/Airplanes-Co/DAL_customers/internal/routes/customers"
	"github.com/gorilla/mux"
)

func Run() {
	router := mux.NewRouter()

	router.HandleFunc("/dal/flights/postgres/get", customers.GetCustomers).Methods(http.MethodGet)
	router.HandleFunc("/dal/flights/postgres/get/{id:[0-9]+}", customers.GetCustomerById).Methods(http.MethodGet)
	router.HandleFunc("/dal/flights/postgres/insert", customers.InsertCustomer).Methods(http.MethodPost)
	router.HandleFunc("/dal/flights/postgres/update", customers.UpdateCustomer).Methods(http.MethodPatch)
	router.HandleFunc("/dal/flights/postgres/delete/{id:[0-9]+}", customers.DeleteCustomer).Methods(http.MethodDelete)

	http.ListenAndServe(":12000", router)
}
