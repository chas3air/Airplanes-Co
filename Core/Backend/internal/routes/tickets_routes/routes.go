package tickets_routes

import (
	"bytes"
	"encoding/json"
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

func GetTicketHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Get tickets backend process...")

	resp, err := httpClient.Get(config.Management_cache_api_url + "/" + config.KEY_FOR_TICKETS)
	if err != nil {
		handleError(w, err)
		return
	}
	defer resp.Body.Close()
	log.Println("response code is:", resp.StatusCode)

	if resp.StatusCode == http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println("Error:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Println("backend sent to client:", string(body))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.StatusCode)
		w.Write(body)
		log.Println("Successfully fetched tickets from cache.")
		return
	}

	log.Println("Cache is not used")

	resp, err = httpClient.Get(config.Management_tickets_api_url)
	if err != nil {
		log.Println("Error:", err)
		handleError(w, err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	message := models.Message{
		Key: config.KEY_FOR_TICKETS,
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
	log.Println("Successfully fetched all tickets from management service.")
}

func GetTicketByIdHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Get ticket by id backend process...")

	id, ok := mux.Vars(r)["id"]
	if !ok {
		log.Println("Bad request: cannot get id")
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	escapedID := url.PathEscape(config.KEY_FOR_TICKETS + ":" + id)
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
			log.Println("Error:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.StatusCode)
		w.Write(body)
		log.Println("Successfully fetched ticket by id from cache.")
	}

	resp, err = httpClient.Get(config.Management_tickets_api_url + "/" + id)
	if err != nil {
		handleError(w, err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	message := models.Message{
		Key: config.KEY_FOR_TICKETS,
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
	log.Println("Successfully fetched ticket by id from cache.")
}

func InsertTicketHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Insert ticket backend process...")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Bad request: cannot read request body")
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	resp, err := httpClient.Post(config.Management_tickets_api_url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Println("Cannot do post request to management tickets")
		handleError(w, err)
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
	w.WriteHeader(resp.StatusCode)
	w.Write(bs)
	log.Println("Successfully inserted ticket.")

	go sendtoCache(body)
}

func UpdateTicketHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Update ticket backend process...")

	req, err := http.NewRequest(http.MethodPatch, config.Management_tickets_api_url, r.Body)
	if err != nil {
		log.Println("Cannot create request")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		handleError(w, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 == 4 {
		log.Println("Response code:", resp.StatusCode)
		http.Error(w, fmt.Errorf("response code: %v", resp.StatusCode).Error(), resp.StatusCode)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Bad response: cannot read response body")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
	log.Println("Successfully updated ticket")

	go sendtoCache(body)
}

func DeleteTicketHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Delete ticket backend process...")

	id, ok := mux.Vars(r)["id"]
	if !ok {
		log.Println("Bad request: cannot get id")
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	escapedID := url.PathEscape(id)
	rawURL := config.Management_tickets_api_url + "/" + escapedID
	req, err := http.NewRequest(http.MethodDelete, rawURL, nil)
	if err != nil {
		log.Println("Cannot create request")
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		log.Println("Bad response: cannot read response body")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
	log.Println("Successfully deleted ticket")

	go sendtoCache(body)
}

func handleError(w http.ResponseWriter, err error) {
	if urlErr, ok := err.(*url.Error); ok && urlErr.Timeout() {
		log.Println("Request timed out:", urlErr)
		http.Error(w, "Request timed out", http.StatusGatewayTimeout)
	} else {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func sendtoCache(body []byte) {
	var ticket models.Ticket
	if err := json.Unmarshal(body, &ticket); err != nil {
		log.Println("Error unmarshalling response:", err)
		return
	}

	message := models.Message{
		Key: config.KEY_FOR_TICKETS + ":" + ticket.Id.String(),
		Value: models.CacheItem{
			Value:      string(body),
			Expiration: config.VALUE_EXPIRATION_TIME,
			SetedTime:  time.Now(),
		},
	}
	go service.SaveToCache(message)
}
