package routes

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/chas3air/Airplanes-Co/Core/DAL_airplanes/internal/models"
	"github.com/chas3air/Airplanes-Co/Core/DAL_airplanes/internal/service"
	"github.com/chas3air/Airplanes-Co/Core/DAL_airplanes/internal/storage"
	"github.com/gorilla/mux"
)

var AirplanesDB = storage.MustGetInstanceOfAirplanesStorage("psql")
var limitTime = service.GetLimitTime("PSQL_LIMIT_RESPONSE_TIME")

// GetAirplanesHandler handles the HTTP request to retrieve all airplanes.
// It fetches airplane data from the database, marshals it to JSON, and writes it to the response.
// In case of errors, it logs the issue and returns an appropriate HTTP status code.
func GetAirplanesHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Fetching all airplanes")

	ctx, cancel := context.WithTimeout(context.Background(), limitTime)
	defer cancel()

	entities, err := AirplanesDB.GetAll(ctx)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("Request timed out")
			http.Error(w, "Request timed out", http.StatusGatewayTimeout)
			return
		} else {
			log.Println("Error retrieving airplanes")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	airplanes, ok := entities.([]models.Airplane)
	if !ok {
		log.Println("Invalid data type")
		http.Error(w, "Invalid data type", http.StatusInternalServerError)
		return
	}

	bs, err := json.Marshal(airplanes)
	if err != nil {
		log.Println("Cannot marshal object")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bs)
	log.Println("Successfully fetched airplanes.")
}

// InsertAirplaneHandler handles the HTTP request to insert a new airplane.
// It reads the request body, unmarshals it into an airplane model, and attempts to insert it into the database.
// If successful, it responds with the inserted airplane data; otherwise, it logs the error and returns an appropriate status code.
func InsertAirplaneHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Inserting airplane")

	ctx, cancel := context.WithTimeout(context.Background(), limitTime)
	defer cancel()

	bs, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Bad request: cannot read request body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var airplane models.Airplane

	err = json.Unmarshal(bs, &airplane)
	if err != nil {
		log.Println("Cannot unmarshal request body to airplane")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	obj, err := AirplanesDB.Insert(ctx, airplane)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("Request timed out")
			http.Error(w, "Request timed out", http.StatusGatewayTimeout)
			return
		} else {
			log.Printf("Error inserting airplane with ID: %d\n", airplane.Id)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	out_bs, err := json.Marshal(obj)
	if err != nil {
		log.Println("Cannot marshal inserted airplane")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out_bs)
	log.Printf("Successfully inserted airplane with ID: %d\n", airplane.Id)
}

// UpdateAirplaneHandler handles the HTTP request to update an existing airplane.
// It reads the request body, unmarshals it into an airplane model, and attempts to update the corresponding record in the database.
// If the update is successful, it returns the updated airplane data; otherwise, it logs the error and responds with an appropriate status code.
func UpdateAirplaneHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Updating airplane")

	ctx, cancel := context.WithTimeout(context.Background(), limitTime)
	defer cancel()

	bs, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Bad request: cannot read request body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var airplane models.Airplane
	err = json.Unmarshal(bs, &airplane)
	if err != nil {
		log.Println("Cannot unmarshal request body to airplane")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	obj, err := AirplanesDB.Update(ctx, airplane)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("Request timed out")
			http.Error(w, "Request timed out", http.StatusGatewayTimeout)
			return
		} else {
			log.Printf("Error updating airplane with ID: %d\n", airplane.Id)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	out_bs, err := json.Marshal(obj)
	if err != nil {
		log.Println("Cannot marshal updated airplane")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out_bs)
	log.Printf("Successfully updated airplane with ID: %d", airplane.Id)
}

// DeleteAirplaneHandler handles the HTTP request to delete an airplane by its ID.
// It extracts the ID from the request URL, attempts to delete the airplane record from the database,
// and returns the deleted object data if successful. In case of errors, it logs the issue and returns an appropriate status code.
func DeleteAirplaneHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Deleting airplane")

	ctx, cancel := context.WithTimeout(context.Background(), limitTime)
	defer cancel()

	id := mux.Vars(r)["id"]

	obj, err := AirplanesDB.Delete(ctx, id)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("Request timed out")
			http.Error(w, "Request timed out", http.StatusGatewayTimeout)
			return
		} else {
			log.Println("Error deleting airplane with ID:", id)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	out_bs, err := json.Marshal(obj)
	if err != nil {
		log.Println("Cannot marshal deleted ticket")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out_bs)
	log.Printf("Successfully deleted airplane with ID: %v", id)
}
