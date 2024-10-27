package router

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/chas3air/Airplanes-Co/Flights_catalog/internal/service"
	"github.com/gorilla/mux"
)

var dal_flights_url = os.Getenv("DAL_FLIGHTS_URL")
var limitTime = service.GetLimitTime()

// GetFlightsHandler handles GET requests to fetch all flights.
// It retrieves the list of flights from the DAL and returns it in JSON format.
// The response contains a JSON array of flight objects.
func GetFlightsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Starting to fetch all flights")

	client := http.Client{
		Timeout: limitTime,
	}

	response, err := client.Get(dal_flights_url)
	if err != nil {
		log.Printf("Error connecting to %s: %v\n", dal_flights_url, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Printf("Received non-OK status code: %d\n", response.StatusCode)
		http.Error(w, fmt.Sprintf("Received status code: %d", response.StatusCode), response.StatusCode)
		return
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
	log.Println("Successfully fetched all flights")
}

// GetFlightByIDHandler handles GET requests to fetch a flight by its ID.
// It retrieves the flight details from the DAL and returns them in JSON format.
// The response contains a JSON object representing the flight.
func GetFlightByIDHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Starting to fetch flight by ID")

	id := mux.Vars(r)["id"]

	client := http.Client{
		Timeout: limitTime,
	}

	response, err := client.Get(fmt.Sprintf("%s/%s", dal_flights_url, id))
	if err != nil {
		log.Printf("Error connecting to %s: %v\n", dal_flights_url+"/"+id, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Printf("Received non-OK status code: %d\n", response.StatusCode)
		http.Error(w, "Flight not found", http.StatusNotFound)
		return
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
	log.Println("Successfully fetched flight by ID.")
}

// SearchFlightsHandler handles requests to search for flights.
// It retrieves the list of flights from the DAL and returns them in JSON format.
// The search logic is to be implemented based on request parameters.
func SearchFlightsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Starting to search for flights")

	// TODO: Implement actual search logic based on request parameters

	client := http.Client{
		Timeout: limitTime,
	}

	response, err := client.Get(dal_flights_url) // Placeholder for actual search URL
	if err != nil {
		log.Printf("Error connecting to %s: %v\n", dal_flights_url, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Printf("Received non-OK status code: %d\n", response.StatusCode)
		http.Error(w, fmt.Sprintf("Received status code: %d", response.StatusCode), response.StatusCode)
		return
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
	log.Println("Successfully searched for flights.")
}

// StatsFlightsHandler handles requests for flight statistics.
// Currently not implemented, but should return relevant statistics about flights.
func StatsFlightsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Fetching flight statistics (not implemented)")

	// (Implementation goes here)
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}
