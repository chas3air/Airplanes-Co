package flights

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/chas3air/Airplanes-Co/DAL/internal/models"
	"github.com/chas3air/Airplanes-Co/DAL/internal/storage"
	"github.com/gorilla/mux"
)

var FlightsDB = storage.MustGetInstanceOfFlightsStorage("psql")

func GetFlights(w http.ResponseWriter, r *http.Request) {
	log.Println("Fetching all flights")

	entities, err := FlightsDB.GetAll(context.Background())
	if err != nil {
		log.Println("Error retriecing flights")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	flights, ok := entities.([]models.Flight)
	if !ok {
		log.Println("Invalid data type")
		http.Error(w, "Invalid data type", http.StatusInternalServerError)
		return
	}

	bs, err := json.Marshal(flights)
	if err != nil {
		log.Println("Cannot marshal object")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bs)
	log.Println("Succesfully fetched flights.")
}

func GetFlightById(w http.ResponseWriter, r *http.Request) {
	log.Println("Fecthing flight by id")

	id_s := mux.Vars(r)["id"]
	id, err := strconv.Atoi(id_s)
	if err != nil {
		log.Printf("Bad request: invalid ID: %s\n", id_s)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	entity, err := FlightsDB.GetById(context.Background(), id)
	if err != nil {
		log.Printf("Error retrieving flight by id: %d\n", id)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	flight, ok := entity.(models.Flight)
	if !ok {
		log.Println("Invalid data type")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bs, err := json.Marshal(flight)
	if err != nil {
		log.Printf("Cannot marshal object: %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bs)
	log.Println("Succesfully fatched flight by id.")
}

func InsertFlight(w http.ResponseWriter, r *http.Request) {
	log.Println("Flight insertion")

	bs, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Bad request: cannot read request body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var flight models.Flight
	err = json.Unmarshal(bs, &flight)
	if err != nil {
		log.Printf("Cannot unmarshal request body to flight")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	obj, err := FlightsDB.Insert(context.Background(), flight)
	if err != nil {
		log.Printf("Error inserting flight with id: %d\n", flight.Id)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	out_bs, err := json.Marshal(obj)
	if err != nil {
		log.Println("Cannot marshal inserted flight")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out_bs)
	log.Printf("Succesfully inserted flight with id: %d\n", flight.Id)
}

func UpdateFlight(w http.ResponseWriter, r *http.Request) {
	log.Println("Flight updating")
	bs, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Bad request: cannot read request body")
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	defer r.Body.Close()

	var flight models.Flight
	err = json.Unmarshal(bs, &flight)
	if err != nil {
		log.Println("Cannot unmarshal request body to flight")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	obj, err := FlightsDB.Update(context.Background(), flight)
	if err != nil {
		log.Printf("Error updating flight with id: %d\n", flight.Id)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	out_bs, err := json.Marshal(obj)
	if err != nil {
		log.Println("Canot mardhsl updated flight")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out_bs)
	log.Printf("Successfully updated customer with ID: %d", flight.Id)
}

func DeleteFlight(w http.ResponseWriter, r *http.Request) {
	log.Println("Flight deleting")

	id_s := mux.Vars(r)["id"]
	id, err := strconv.Atoi(id_s)
	if err != nil {
		log.Println("Bad request: wrong flight id")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	obj, err := FlightsDB.Delete(context.Background(), id)
	if err != nil {
		log.Println("Error deliting flight with id:", id)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	out_bs, err := json.Marshal(obj)
	if err != nil {
		log.Println("Cannot marshal deleted flight")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out_bs)
	log.Printf("Successfully deleted flight with ID: %d", id)
}
