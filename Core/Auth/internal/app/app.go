package app

import (
	"net/http"

	"github.com/chas3air/Airplanes-Co/Core/Auth/internal/router"
	"github.com/gorilla/mux"
)

func Run() {
	r := mux.NewRouter()

	r.HandleFunc("/auth/signin", router.SignInHandler).Methods(http.MethodGet)
	r.HandleFunc("/auth/signup", router.SignUpHandler).Methods(http.MethodPost)

	http.ListenAndServe(":12005", r)
}
