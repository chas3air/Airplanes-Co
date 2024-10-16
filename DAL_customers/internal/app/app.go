package app

import (
	"net/http"

	"github.com/chas3air/Airplanes-Co/DAL_customers/internal/routes/customers"
	"github.com/gorilla/mux"
)

func Run() {
	router := mux.NewRouter()

	router.HandleFunc("/postgres/customers/get", customers.GetCustomers).Methods(http.MethodGet)
	router.HandleFunc("/postgres/customers/get/{id:[0-9]+}", customers.GetCustomerById).Methods(http.MethodGet)
	router.HandleFunc("/postgres/customers/insert", customers.InsertCustomer).Methods(http.MethodPut)
	router.HandleFunc("/postgres/customers/update", customers.UpdateCustomer).Methods(http.MethodPatch)
	router.HandleFunc("/postgres/customers/delete/{id:[0-9]+}", customers.DeleteCustomer).Methods(http.MethodDelete)

	http.ListenAndServe(":8057", router)
}
