package app

import (
	"net/http"

	"github.com/chas3air/Airplanes-Co/Core/OrderTicket/internal/routes"
	"github.com/gorilla/mux"
)

// TODO: протестировать

func Run() {
	r := mux.NewRouter()

	r.HandleFunc("/order", routes.OrderTicketHandler).Methods(http.MethodPost)

	http.ListenAndServe(":12009", r)
}
