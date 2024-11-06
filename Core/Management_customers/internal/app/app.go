package app

import (
	"net/http"

	"github.com/chas3air/Airplanes-Co/Core/Management_customers/internal/routes"
	"github.com/gorilla/mux"
)

func Run() {
	r := mux.NewRouter()

	r.HandleFunc("/management-customers/customers", routes.GetCustomersHandler).Methods(http.MethodGet)
	r.HandleFunc("/management-customers/customers/{id}/", routes.GetCustomerByIdHandler).Methods(http.MethodGet)
	r.HandleFunc("/management-customers/customers/login", routes.GetCustomerByLoginAndPasswordHandler).Methods(http.MethodGet)
	r.HandleFunc("/management-customers/customers", routes.InsertCustomerHandler).Methods(http.MethodPost)
	r.HandleFunc("/management-customers/customers", routes.UpdateCustomerHandler).Methods(http.MethodPatch)
	r.HandleFunc("/management-customers/customers/{id}", routes.DeleteCustomerHandler).Methods(http.MethodDelete)

	http.ListenAndServe(":12007", r)
}
