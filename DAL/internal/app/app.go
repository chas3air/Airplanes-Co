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
	router.HandleFunc("/postgres/customers/insert", routes.InsertCustomer).Methods(http.MethodPost)
	router.HandleFunc("/postgres/customers/update", routes.UpdateCustomer).Methods(http.MethodPost)
	router.HandleFunc("/postgres/customers/delete", routes.DeleteCustomer).Methods(http.MethodPost)

	http.ListenAndServe(":8056", router)
}
