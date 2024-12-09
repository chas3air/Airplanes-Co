package app

import (
	"net/http"

	"github.com/chas3air/Airplanes-Co/Core/Checkbook/internal/routes"
	"github.com/gorilla/mux"
)

// TODO: протестировать

func Run() {
	r := mux.NewRouter()

	r.HandleFunc("/checkbook/{id}", routes.GetTicketsFromCheckbookHandler).Methods(http.MethodGet)
	r.HandleFunc("/checkbook/{id}", routes.InsertTicketsToCheckbookHandler).Methods(http.MethodPost)

	http.ListenAndServe(":12009", r)
}
