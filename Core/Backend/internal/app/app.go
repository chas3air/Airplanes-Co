package app

import (
	"net/http"

	"github.com/gorilla/mux"
)

func Run() {
	r := mux.NewRouter()

	http.ListenAndServe(":84", r)
}
