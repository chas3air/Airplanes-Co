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

//TODO: переделать
func GetTicketsHandler(w http.ResponseWriter, r *http.Request) {
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
			handleError(w, err)
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
		handleError(w, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("Failed to fetch tickets from management service:", resp.StatusCode)
		http.Error(w, "Failed to fetch tickets", resp.StatusCode)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		handleError(w, err)
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
			handleError(w, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.StatusCode)
		w.Write(body)
		log.Println("Successfully fetched ticket by id from cache.")
		return
	}

	resp, err = httpClient.Get(config.Management_tickets_api_url + "/" + id)
	if err != nil {
		handleError(w, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("Failed to fetch ticket from management service:", resp.StatusCode)
		http.Error(w, "Failed to fetch ticket", resp.StatusCode)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		handleError(w, err)
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
	log.Println("Successfully fetched ticket by id from management service.")
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
		handleError(w, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		log.Println("Failed to insert ticket:", resp.StatusCode)
		http.Error(w, "Failed to insert ticket", resp.StatusCode)
		return
	}

	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		handleError(w, err)
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

	if resp.StatusCode >= 400 {
		log.Println("Response code:", resp.StatusCode)
		http.Error(w, fmt.Errorf("response code: %v", resp.StatusCode).Error(), resp.StatusCode)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		handleError(w, err)
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

	if resp.StatusCode >= 400 {
		log.Println("Failed to delete ticket:", resp.StatusCode)
		http.Error(w, "Failed to delete ticket", resp.StatusCode)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
	log.Println("Successfully deleted ticket")

	go sendtoCache(body)
}

func InsertTicketToCartHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Inserting ticket to cart backend process...")

	var ticket models.Ticket
	err := json.NewDecoder(r.Body).Decode(&ticket)
	if err != nil {
		log.Println("Bad request: cannot read request body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	resp, err := httpClient.Get(config.Management_customers_api_url + "/" + ticket.Owner.Id.String() + "/")
	if err != nil {
		log.Println("Error fetching customer data:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("Failed to fetch customer data:", resp.StatusCode)
		http.Error(w, "Failed to fetch customer data", resp.StatusCode)
		return
	}

	var customer models.Customer
	err = json.NewDecoder(resp.Body).Decode(&customer)
	if err != nil {
		log.Println("Bad response: cannot read response body")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ticket.Owner = customer

	bs, err := json.Marshal(ticket)
	if err != nil {
		log.Println("Cannot marshall object to json")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err = httpClient.Post(config.Cart_api_url, "application/json", bytes.NewBuffer(bs))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		log.Println("Failed to insert ticket to cart:", resp.StatusCode)
		http.Error(w, "Failed to insert ticket to cart", resp.StatusCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	_, err = w.Write(bs)
	if err != nil {
		log.Println("Error writing response:", err)
		return
	}
	log.Println("Successfully inserted ticket to cart.")

	go sendtoCache(bs)
}

func GetTicketsFromCartHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Fetching all tickets from cart backend process...")

	resp, err := httpClient.Get(config.Cart_api_url)
	if err != nil {
		handleError(w, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("Failed to fetch tickets from management service:", resp.StatusCode)
		http.Error(w, "Failed to fetch tickets", resp.StatusCode)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		handleError(w, err)
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

func GetPurchasedTicketHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Fetching all purchased tickets backend process...")

	ownerId := r.URL.Query().Get("ownerId")
	if ownerId == "" {
		log.Println("Bad request: ownerId is required")
		http.Error(w, "ownerId is required", http.StatusBadRequest)
		return
	}

	resp, err := httpClient.Get(config.Purchased_tickets_api_url + "?ownerId=" + ownerId)
	if err != nil {
		handleError(w, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		log.Println("Received response:", resp.Status)
		http.Error(w, "Error response code", resp.StatusCode)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		handleError(w, err)
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
	w.WriteHeader(http.StatusOK)
	w.Write(body)
	log.Println("Successfully fetched all purchased tickets.")
}

func PayForTickets(w http.ResponseWriter, r *http.Request) {
	log.Println("Pay for tickets backend process...")

	var card_info models.Card
	err := json.NewDecoder(r.Body).Decode(&card_info)
	if err != nil {
		log.Println("Bad request: cannot parse body to cardInfo")
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	payInfo := models.PayInfo{
		Cost:     66,
		CardInfo: card_info,
	}
	body, err := json.Marshal(payInfo)
	if err != nil {
		log.Println("Cannot marshall ticket")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := httpClient.Post(config.Payment_api_url+"/pay", "application/json", bytes.NewBuffer(body))
	if err != nil {
		handleError(w, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		log.Println("received response: ", resp.Status)
		http.Error(w, "Error response code", resp.StatusCode)
		return
	}

	req, err := http.NewRequest(http.MethodDelete, config.Cart_api_url+"/clear", nil)
	if err != nil {
		handleError(w, err)
		return
	}

	resp, err = httpClient.Do(req)
	if err != nil {
		handleError(w, err)
		return
	}
	defer resp.Body.Close()

	/*
		тут нужно переделать dal-tickets так чтобы он принимал не одно значение а несколько
	*/

	if resp.StatusCode >= 400 {
		log.Println("received response: ", resp.Status)
		http.Error(w, "Error response code", resp.StatusCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
	log.Println("Successfully fetched all tickets from cart.")
}

func handleError(w http.ResponseWriter, err error) {
	if urlErr, ok := err.(*url.Error); ok {
		if urlErr.Timeout() {
			log.Println("Request timed out:", urlErr)
			http.Error(w, "Request timed out", http.StatusGatewayTimeout)
		} else {
			log.Println("Service unavailable:", urlErr)
			http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		}
	} else {
		log.Println("General error:", err)
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
