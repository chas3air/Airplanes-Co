package app

import (
	"net/http"

	"github.com/chas3air/Airplanes-Co/Cart/internal/router/cart"
	"github.com/gorilla/mux"
)

func Run() {
	router := mux.NewRouter()

	router.HandleFunc("/api/get", cart.GetAllTicketsHandler).Methods(http.MethodGet)
	router.HandleFunc("/api/insert", cart.InsertTicketHandler).Methods(http.MethodPost)
	router.HandleFunc("/api/update", cart.UpdateTicketHandler).Methods(http.MethodPatch)
	router.HandleFunc("/api/delete/{id:[0-9]+}", cart.DeleteTicketHandler).Methods(http.MethodDelete)
	router.HandleFunc("/api/clear", cart.ClearHandler).Methods(http.MethodDelete)

	http.ListenAndServe(":12003", router)
}
