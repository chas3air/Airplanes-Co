package app

import (
	"net/http"

	"github.com/chas3air/Airplanes-Co/OrderTicket/internal/router"
	"github.com/gorilla/mux"
)

func Run() {
	r := mux.NewRouter()

	r.HandleFunc("/order", router.OrderTicketHandler).Methods(http.MethodPost)

	http.ListenAndServe(":12009", r)
}
