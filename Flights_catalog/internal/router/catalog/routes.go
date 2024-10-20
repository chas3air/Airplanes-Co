package catalog

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
func GetFlightsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Fetching flights")

	client := http.Client{
		Timeout: limitTime,
	}

	response, err := client.Get(dal_flights_url + "/get")
	if err != nil {
		log.Printf("Error connecting to %s: %v\n", dal_flights_url+"/get", err)
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
	log.Println("Successfully fetched flights")
}

// GetFlightByIDHandler handles GET requests to fetch a flight by its ID.
// It retrieves the flight details from the DAL and returns them in JSON format.
func GetFlightByIDHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Fetching flight by ID")

	id_s := mux.Vars(r)["id"]

	client := http.Client{
		Timeout: limitTime,
	}

	response, err := client.Get(fmt.Sprintf("%s/get/%s", dal_flights_url, id_s))
	if err != nil {
		log.Printf("Error connecting to %s: %v\n", dal_flights_url+"/get/"+id_s, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Printf("Received non-OK status code: %d\n", response.StatusCode)
		http.Error(w, "Status code is not available", response.StatusCode)
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
// It retrieves the list of flights from the DAL (to be implemented with search logic).
func SearchFlightsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Searching flights")

	client := http.Client{
		Timeout: limitTime,
	}

	response, err := client.Get(dal_flights_url + "/get")
	if err != nil {
		log.Printf("Error connecting to %s: %v\n", dal_flights_url+"/get", err)
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

	// TODO: Implement actual search logic based on request parameters

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body) // Placeholder response for now
	log.Println("Successfully searched flights.")
}

// StatsFlightsHandler handles requests for flight statistics.
// (Implementation to be completed based on requirements.)
func StatsFlightsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Fetching flight statistics (not implemented)")

	// (Implementation goes here)
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}
