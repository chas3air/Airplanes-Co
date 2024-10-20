package router

import (
	"bytes"
	"io"
	"log"
	"net/http"

	"github.com/chas3air/Airplanes-Co/Management_flights/internal/config"
	"github.com/gorilla/mux"
)

var database_url = config.DATABASE_URL

// GetAllFlightsHandler handles a GET request to fetch all flights.
// Returns a list of flights in JSON format.
func GetAllFlightsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Fetching all flights")

	resp, err := http.Get(database_url + "/get")
	if err != nil {
		log.Println("Cannot send request to", database_url)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Bad response: cannot read response body")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
	log.Println("Successfully fetched all flights.")
}

// GetFlightByIdHandler handles a GET request to retrieve a flight by ID.
// Takes the flight ID from the URL and returns the flight data in JSON format.
func GetFlightByIdHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Retrieving flight by ID")

	id_s := mux.Vars(r)["id"]

	resp, err := http.Get(database_url + "/get/" + id_s)
	if err != nil {
		log.Println("Cannot send request to", database_url)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Cannot read response body")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
	log.Println("Successfully retrieved flight by ID.")
}

// InsertFlightsHandler handles a POST request to add a new flight.
// Takes flight data in JSON format and returns confirmation of the insertion.
func InsertFlightsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Inserting flight")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Bad request: cannot read request body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close() // Close the request body after reading

	resp, err := http.Post(database_url+"/insert", "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Println("Error sending request")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Cannot read response body")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bs)
	log.Println("Successfully inserted flight.")
}

// UpdateFlightHandler handles a PATCH request to update flight information.
// Takes flight data in JSON format and returns the updated flight data.
func UpdateFlightHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Updating flight")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Bad request: cannot read request body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	req, err := http.NewRequest(http.MethodPatch, database_url+"/update", bytes.NewBuffer(body))
	if err != nil {
		log.Println("Error creating request")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error sending request")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(responseBody)
	log.Println("Successfully updated flight.")
}

// DeleteFlightHandler handles a DELETE request to remove a flight by ID.
// Takes the flight ID from the URL and returns confirmation of the deletion.
func DeleteFlightHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Deleting flight")

	id_s := mux.Vars(r)["id"]

	req, err := http.NewRequest(http.MethodDelete, database_url+"/delete/"+id_s, nil)
	if err != nil {
		log.Println("Error creating request")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error sending request")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Cannot read response body")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(responseBody)
	log.Println("Successfully deleted flight.")
}
