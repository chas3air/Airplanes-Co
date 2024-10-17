package app

import (
	"net/http"

	"github.com/chas3air/Airplanes-Co/Cart/internal/router/cart"
	"github.com/gorilla/mux"
)

func Run() {
	router := mux.NewRouter()

	router.HandleFunc("/get", cart.GetAllTicketsHandler).Methods(http.MethodGet)
	router.HandleFunc("/insert", cart.InsertTicketHandler).Methods(http.MethodPost)
	router.HandleFunc("/update", cart.UpdateTicketHandler).Methods(http.MethodPatch)
	router.HandleFunc("/delete/{id:[0-9]+}", cart.DeleteTicketHandler).Methods(http.MethodDelete)

	http.ListenAndServe(":12003", router)
}
