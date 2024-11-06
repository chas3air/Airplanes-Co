package app

import (
	"net/http"

	"github.com/chas3air/Airplanes-Co/Core/Auth/internal/router"
	"github.com/gorilla/mux"
)

func Run() {
	r := mux.NewRouter()

	r.HandleFunc("/auth/singin", router.SignInHandler).Methods(http.MethodGet)
	r.HandleFunc("/auth/singup", router.SignUpHandler).Methods(http.MethodPost)

	http.ListenAndServe(":12005", r)
}
