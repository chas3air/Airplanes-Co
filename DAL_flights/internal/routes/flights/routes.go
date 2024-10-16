package flights

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/chas3air/Airplanes-Co/DAL_flights/internal/config"
	"github.com/chas3air/Airplanes-Co/DAL_flights/internal/models"
	"github.com/chas3air/Airplanes-Co/DAL_flights/internal/storage"
	"github.com/gorilla/mux"
)

var FlightsDB = storage.MustGetInstanceOfFlightsStorage("psql")

func GetFlights(w http.ResponseWriter, r *http.Request) {
	log.Println("Fetching all flights")

	ctx, cancel := context.WithTimeout(context.Background(), config.PSQL_LIMIT_RESPONSE_TIME*time.Second)
	defer cancel()

	entities, err := FlightsDB.GetAll(ctx)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("Request timed out")
			http.Error(w, "Request timed out", http.StatusGatewayTimeout)
			return
		} else {
			log.Println("Error retrieving flights")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
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
	log.Println("Successfully fetched flights.")
}

func GetFlightById(w http.ResponseWriter, r *http.Request) {
	log.Println("Fetching flight by ID")

	ctx, cancel := context.WithTimeout(context.Background(), config.PSQL_LIMIT_RESPONSE_TIME*time.Second)
	defer cancel()

	id_s := mux.Vars(r)["id"]
	id, err := strconv.Atoi(id_s)
	if err != nil {
		log.Printf("Bad request: invalid ID: %s\n", id_s)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	entity, err := FlightsDB.GetById(ctx, id)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("Request timed out")
			http.Error(w, "Request timed out", http.StatusGatewayTimeout)
			return
		} else {
			log.Println("Error retrieving flight by ID:", id)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	flight, ok := entity.(models.Flight)
	if !ok {
		log.Println("Invalid data type")
		http.Error(w, "Invalid data type", http.StatusInternalServerError)
		return
	}

	bs, err := json.Marshal(flight)
	if err != nil {
		log.Println("Cannot marshal object")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bs)
	log.Println("Successfully fetched flight by ID.")
}

func InsertFlight(w http.ResponseWriter, r *http.Request) {
	log.Println("Inserting flight")

	ctx, cancel := context.WithTimeout(context.Background(), config.PSQL_LIMIT_RESPONSE_TIME*time.Second)
	defer cancel()

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
		log.Println("Cannot unmarshal request body to flight")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	obj, err := FlightsDB.Insert(ctx, flight)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("Request timed out")
			http.Error(w, "Request timed out", http.StatusGatewayTimeout)
			return
		} else {
			log.Printf("Error inserting flight with ID: %d\n", flight.Id)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
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
	log.Printf("Successfully inserted flight with ID: %d\n", flight.Id)
}

func UpdateFlight(w http.ResponseWriter, r *http.Request) {
	log.Println("Updating flight")

	ctx, cancel := context.WithTimeout(context.Background(), config.PSQL_LIMIT_RESPONSE_TIME*time.Second)
	defer cancel()

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
		log.Println("Cannot unmarshal request body to flight")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	obj, err := FlightsDB.Update(ctx, flight)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("Request timed out")
			http.Error(w, "Request timed out", http.StatusGatewayTimeout)
			return
		} else {
			log.Printf("Error updating flight with ID: %d\n", flight.Id)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	out_bs, err := json.Marshal(obj)
	if err != nil {
		log.Println("Cannot marshal updated flight")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out_bs)
	log.Printf("Successfully updated flight with ID: %d", flight.Id)
}

func DeleteFlight(w http.ResponseWriter, r *http.Request) {
	log.Println("Deleting flight")

	ctx, cancel := context.WithTimeout(context.Background(), config.PSQL_LIMIT_RESPONSE_TIME*time.Second)
	defer cancel()

	id_s := mux.Vars(r)["id"]
	id, err := strconv.Atoi(id_s)
	if err != nil {
		log.Println("Bad request: wrong flight ID")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	obj, err := FlightsDB.Delete(ctx, id)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("Request timed out")
			http.Error(w, "Request timed out", http.StatusGatewayTimeout)
			return
		} else {
			log.Println("Error deleting flight with ID:", id)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
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