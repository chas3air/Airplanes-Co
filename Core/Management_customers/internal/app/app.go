package app

import (
	"net/http"

	"github.com/chas3air/Airplanes-Co/Core/Management_customers/internal/routes"
	"github.com/gorilla/mux"
)

func Run() {
	r := mux.NewRouter()

	r.HandleFunc("/customers", routes.GetCustomersHandler).Methods(http.MethodGet)
	r.HandleFunc("/customers/{id}/", routes.GetCustomerByIdHandler).Methods(http.MethodGet)
	r.HandleFunc("/customers/login", routes.GetCustomerByLoginAndPasswordHandler).Methods(http.MethodGet)
	r.HandleFunc("/customers", routes.InsertCustomerHandler).Methods(http.MethodPost)
	r.HandleFunc("/customers", routes.UpdateCustomerHandler).Methods(http.MethodPatch)
	r.HandleFunc("/customers/{id}", routes.DeleteCustomerHandler).Methods(http.MethodDelete)

	http.ListenAndServe(":12007", r)
}
