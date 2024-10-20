package app

import (
	"net/http"

	"github.com/chas3air/Airplanes-Co/Management_customers/internal/router"
	"github.com/gorilla/mux"
)

func Run() {
	r := mux.NewRouter()

	r.HandleFunc("/api/customers/getAllCustomers", router.GetAllCustomersHandler).Methods(http.MethodGet)
	r.HandleFunc("/api/customers/getCustomerById/{id:[0-9]+}", router.GetCustomerByIdHandler).Methods(http.MethodGet)
	r.HandleFunc("/api/customers/getCustomerByLoginAndPassword", router.GetCustomerByLoginAndPasswordHandler).Methods(http.MethodGet)
	r.HandleFunc("/api/customers/insert", router.InsertCustomerHandler).Methods(http.MethodPost)
	r.HandleFunc("/api/customers/update", router.UpdateCustomerHandler).Methods(http.MethodPatch)
	r.HandleFunc("/api/customers/delete/{id:[0-9]+}", router.DeleteCustomerHandler).Methods(http.MethodDelete)

	http.ListenAndServe(":12007", r)
}
