package flights_routes

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/chas3air/Airplanes-Co/Core/Backend/internal/config"
	"github.com/chas3air/Airplanes-Co/Core/Backend/internal/models"
	"github.com/chas3air/Airplanes-Co/Core/Backend/internal/service"
	"github.com/gorilla/mux"
)

var limitTime = service.GetLimitTime("LIMIT_RESPONSE_TIME")

var httpClient = &http.Client{
	Timeout: limitTime,
}

func GetFlightsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Get flights via management-flights")

	resp, err := httpClient.Get(config.Management_cache_api_url + "/" + config.KEY_FOR_FLIGHTS)
	if err != nil {
		handleError(w, err)
		return
	}
	defer resp.Body.Close()
	log.Println("response code is:", resp.StatusCode)

	if resp.StatusCode == http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		log.Println("backend send to client:", string(body))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.StatusCode)
		w.Write(body)
		log.Println("Successfully fetched flights from cache.")
		return
	}

	log.Println("Cache is not used")

	resp, err = httpClient.Get(config.Management_flights_api_url)
	if err != nil {
		handleError(w, err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	message := models.Message{
		Key: "flights",
		Value: models.CacheItem{
			Value:      string(body),
			Expiration: config.VALUE_EXPIRATION_TIME,
			SetedTime:  time.Now(),
		},
	}

	saveToCache(message)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
	log.Println("Successfully fetched all flights from management service.")
}

func GetFlightByIdHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Get flight by id via management-flights")

	id, ok := mux.Vars(r)["id"]
	if !ok {
		log.Println("Bad request: cannot get id")
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	escapedID := url.PathEscape(config.KEY_FOR_FLIGHTS + ":" + id)
	rawURL := config.Management_cache_api_url + "/" + escapedID
	resp, err := httpClient.Get(rawURL)
	if err != nil {
		handleError(w, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.StatusCode)
		w.Write(body)
		log.Println("Successfully fetched flight from cache.")
		return
	}

	resp, err = httpClient.Get(config.Management_flights_api_url + "/" + id)
	if err != nil {
		handleError(w, err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	message := models.Message{
		Key: "flights:" + id,
		Value: models.CacheItem{
			Value:      string(body), // Используем body напрямую
			Expiration: config.VALUE_EXPIRATION_TIME,
			SetedTime:  time.Now(),
		},
	}

	saveToCache(message)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
	log.Println("Successfully fetched flight by id:", id)
}

func InsertFlightHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Insert flight via management-flights")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Bad request: cannot read request body")
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	resp, err := httpClient.Post(config.Management_flights_api_url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		handleError(w, err)
		return
	}
	defer resp.Body.Close()

	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Cannot read response body")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(bs)
	log.Println("Successfully inserted flight.")
}

func UpdateFlightHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Update flight via management-flights")

	req, err := http.NewRequest(http.MethodPatch, config.Management_flights_api_url, r.Body)
	if err != nil {
		log.Println("Cannot create request")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		handleError(w, err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Bad request: cannot read response body")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
	log.Println("Successfully updated flight.")
}

func DeleteFlightHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Delete flight via management-flights")

	id, ok := mux.Vars(r)["id"]
	if !ok {
		log.Println("Bad request: cannot read id")
		http.Error(w, "Bad request: cannot read id", http.StatusBadRequest)
		return
	}

	escapedID := url.PathEscape(id)
	rawURL := config.Management_flights_api_url + "/" + escapedID
	req, err := http.NewRequest(http.MethodDelete, rawURL, nil)
	if err != nil {
		log.Println("Cannot create request")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		handleError(w, err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Bad request: cannot read response body")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
	log.Println("Successfully deleted flight.")
}

// saveToCache saves a message to the cache in a separate goroutine.
func saveToCache(message models.Message) {
	messageJSON, err := json.Marshal(message)
	if err != nil {
		log.Println("Error marshaling JSON:", err)
		return
	}

	resp, err := http.Post(config.Management_cache_api_url, "application/json", bytes.NewBuffer(messageJSON))
	if err != nil {
		log.Println("Error posting to cache:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("Error posting to cache, response code:", resp.StatusCode)
		return
	}

	log.Println("Element successfully set.")
}

// handleError is a helper function to handle errors and return appropriate HTTP status codes.
func handleError(w http.ResponseWriter, err error) {
	if urlErr, ok := err.(*url.Error); ok && urlErr.Timeout() {
		log.Println("Request timed out:", urlErr)
		http.Error(w, "Request timed out", http.StatusGatewayTimeout)
	} else {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
