package routes

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/chas3air/Airplanes-Co/Core/Management_flights/internal/service"
	"github.com/gorilla/mux"
)

var database_url = os.Getenv("DAL_FLIGHTS_URL")
var limitTime = service.GetLimitTime("LIMIT_RESPONSE_TIME")

var httpClient = &http.Client{
	Timeout: limitTime,
}

// GetFlightsHandler handles a GET request to fetch all flights.
// It retrieves all flight data from the database and returns it in JSON format.
// If an error occurs during the process, it responds with an appropriate HTTP status code.
func GetFlightsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Fetching all flights")

	resp, err := httpClient.Get(database_url)
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
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
	log.Println("Successfully fetched all flights.")
}

// GetFlightByIdHandler handles a GET request to retrieve a flight by ID.
// It takes the flight ID from the URL and returns the flight data in JSON format.
// If the flight is not found or an error occurs, it responds with an appropriate HTTP status code.
func GetFlightByIdHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Retrieving flight by ID")

	id_s := mux.Vars(r)["id"]

	resp, err := httpClient.Get(database_url + "/" + id_s)
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
// It expects flight data in JSON format in the request body.
// If successful, it responds with the inserted flight data; otherwise, it responds with an error.
func InsertFlightsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Inserting flight")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Bad request: cannot read request body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	resp, err := httpClient.Post(database_url, "application/json", bytes.NewBuffer(body))
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
// It expects the updated flight data in JSON format in the request body.
// If successful, it responds with the updated flight data; otherwise, it responds with an error.
func UpdateFlightHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Updating flight")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Bad request: cannot read request body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	req, err := http.NewRequest(http.MethodPatch, database_url, bytes.NewBuffer(body))
	if err != nil {
		log.Println("Error creating request")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
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
// It takes the flight ID from the URL and responds with confirmation of the deletion.
// If an error occurs, it responds with an appropriate HTTP status code.
func DeleteFlightHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Deleting flight")

	id_s := mux.Vars(r)["id"]

	req, err := http.NewRequest(http.MethodDelete, database_url+"/"+id_s, nil)
	if err != nil {
		log.Println("Error creating request")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
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
