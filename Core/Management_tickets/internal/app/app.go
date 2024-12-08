package app

import (
	"net/http"

	"github.com/chas3air/Airplanes-Co/Core/Management_tickets/internal/routes"
	"github.com/gorilla/mux"
)

func Run() {
	r := mux.NewRouter()

	r.HandleFunc("/tickets", routes.GetTicketsHandler).Methods(http.MethodGet)
	r.HandleFunc("/tickets/{id}", routes.GetTicketByIdHandler).Methods(http.MethodGet)
	r.HandleFunc("/tickets", routes.InsertTicketHandler).Methods(http.MethodPost)
	r.HandleFunc("/tickets", routes.UpdateTicketHandler).Methods(http.MethodPatch)
	r.HandleFunc("/tickets/{id}", routes.DeleteTicketHandler).Methods(http.MethodDelete)

	http.ListenAndServe(":12008", r)
}
