package routes

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func GetFlightsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Get flights via management-flights")
}
func InsertFlightHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Insert flight via management-flights")
}

func DeleteFlightHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Delete flight via management-flights")

	id, ok := mux.Vars(r)["id"]
	if !ok {
		log.Println("Bad request: cannot read var id")
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	_ = id
}
