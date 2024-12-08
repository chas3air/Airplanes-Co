package flights_routes

import (
	"bytes"
	"fmt"
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
	log.Println("Get flights backend process...")

	resp, err := httpClient.Get(config.Management_cache_api_url + "/" + config.KEY_FOR_FLIGHTS)
	if err != nil {
		handleError(w, err)
		return
	}
	defer resp.Body.Close()

	log.Println("Response code is:", resp.StatusCode)

	if resp.StatusCode == http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println("Error reading response body:", err)
			http.Error(w, "Error reading response body", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.StatusCode)
		w.Write(body)
		log.Println("Successfully fetched flights from cache.")
		return
	}

	log.Println("Cache is not used")

	resp, err = httpClient.Get(config.Management_flights_api_url)
	if err != nil {
		log.Println("Error fetching flights from management service:", err)
		handleError(w, err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		http.Error(w, "Error reading response body", http.StatusInternalServerError)
		return
	}

	message := models.Message{
		Key: config.KEY_FOR_FLIGHTS,
		Value: models.CacheItem{
			Value:      string(body),
			Expiration: config.VALUE_EXPIRATION_TIME,
			SetedTime:  time.Now(),
		},
	}

	service.SaveToCache(message)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
	log.Println("Successfully fetched all flights from management service.")
}

func GetFlightByIdHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Get flight by id backend process...")

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
			log.Println("Error reading response body:", err)
			http.Error(w, "Error reading response body", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.StatusCode)
		w.Write(body)
		log.Println("Successfully fetched flight by id from cache.")
		return
	}

	resp, err = httpClient.Get(config.Management_flights_api_url + "/" + id)
	if err != nil {
		log.Println("Error fetching flight from management service:", err)
		handleError(w, err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		http.Error(w, "Error reading response body", http.StatusInternalServerError)
		return
	}

	message := models.Message{
		Key: config.KEY_FOR_FLIGHTS + ":" + id,
		Value: models.CacheItem{
			Value:      string(body),
			Expiration: 5,
			SetedTime:  time.Now(),
		},
	}
	service.SaveToCache(message)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
	log.Println("Successfully fetched flight by id:", id)
}

func InsertFlightHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Insert flight backend process...")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Bad request: cannot read request body")
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	resp, err := httpClient.Post(config.Management_flights_api_url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Println("Cannot do post request to management flights:", err)
		handleError(w, err)
		return
	}
	defer resp.Body.Close()

	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Cannot read response body:", err)
		http.Error(w, "Error reading response body", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(bs)
	log.Println("Successfully inserted flight.")
}

func UpdateFlightHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Update flight backend process...")

	req, err := http.NewRequest(http.MethodPatch, config.Management_flights_api_url, r.Body)
	if err != nil {
		log.Println("Cannot create request:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		handleError(w, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("Response code:", resp.StatusCode)
		http.Error(w, fmt.Errorf("response code: %v", resp.StatusCode).Error(), resp.StatusCode)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Bad response: cannot read response body:", err)
		http.Error(w, "Error reading response body", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
	log.Println("Successfully updated flight.")
}

func DeleteFlightHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Delete flight backend process...")

	id, ok := mux.Vars(r)["id"]
	if !ok {
		log.Println("Bad request: cannot get id")
		http.Error(w, "Bad request: cannot read id", http.StatusBadRequest)
		return
	}

	escapedID := url.PathEscape(id)
	rawURL := config.Management_flights_api_url + "/" + escapedID
	req, err := http.NewRequest(http.MethodDelete, rawURL, nil)
	if err != nil {
		log.Println("Cannot create request:", err)
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
		log.Println("Bad response: cannot read response body:", err)
		http.Error(w, "Error reading response body", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
	log.Println("Successfully deleted flight.")
}

func handleError(w http.ResponseWriter, err error) {
	if urlErr, ok := err.(*url.Error); ok && urlErr.Timeout() {
		log.Println("Request timed out:", urlErr)
		http.Error(w, "Request timed out", http.StatusGatewayTimeout)
	} else {
		log.Println("Error occurred:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
