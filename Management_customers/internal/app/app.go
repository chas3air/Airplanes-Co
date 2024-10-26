package app

import (
	"net/http"

	"github.com/chas3air/Airplanes-Co/Management_customers/internal/router"
	"github.com/gorilla/mux"
)

func Run() {
	r := mux.NewRouter()

	r.HandleFunc("/management-customers/customers", router.GetCustomersHandler).Methods(http.MethodGet)
	r.HandleFunc("/management-customers/customers/{id}", router.GetCustomerByIdHandler).Methods(http.MethodGet)
	r.HandleFunc("/management-customers/customers/login", router.GetCustomerByLoginAndPasswordHandler).Methods(http.MethodGet)
	r.HandleFunc("/management-customers/customers", router.InsertCustomerHandler).Methods(http.MethodPost)
	r.HandleFunc("/management-customers/customers", router.UpdateCustomerHandler).Methods(http.MethodPatch)
	r.HandleFunc("/management-customers/customers/{id}", router.DeleteCustomerHandler).Methods(http.MethodDelete)

	http.ListenAndServe(":12007", r)
}
