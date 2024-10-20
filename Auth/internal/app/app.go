package app

import (
	"net/http"

	"github.com/chas3air/Airplanes-Co/Auth/internal/routes/auth"
	"github.com/gorilla/mux"
)

func Run() {
	router := mux.NewRouter()

	router.HandleFunc("/api/auth/singin", auth.SignInHandler).Methods(http.MethodGet)
	router.HandleFunc("/api/auth/singup", auth.SignUpHandler).Methods(http.MethodPost)

	http.ListenAndServe(":12005", router)
}
