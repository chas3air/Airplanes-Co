package tickets_routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
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

func GetTicketsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Get tickets backend process...")

	resp, err := httpClient.Get(config.Management_cache_api_url + "/" + config.KEY_FOR_TICKETS)
	if err != nil {
		handleError(w, err, "Error fetching tickets from cache")
		return
	}
	defer resp.Body.Close()

	log.Println("Response code is:", resp.StatusCode)

	if resp.StatusCode == http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			handleError(w, err, "Error reading response body from cache")
			return
		}

		log.Println("Backend sent to client:", string(body))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.StatusCode)
		w.Write(body)
		log.Println("Successfully fetched tickets from cache.")
		return
	}

	log.Println("Cache is not used")

	resp, err = httpClient.Get(config.Management_tickets_api_url)
	if err != nil {
		handleError(w, err, "Error fetching tickets from management service")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		handleError(w, fmt.Errorf("failed to fetch tickets from management service: %v", resp.StatusCode), "Failed to fetch tickets")
		return
	}

	var tickets []models.Ticket
	err = json.NewDecoder(resp.Body).Decode(&tickets)
	if err != nil {
		handleError(w, err, "cannot parse response body to tickets slice")
		return
	}

	for i, v := range tickets {
		log.Println("FlightId:", v.FlightInfo.Id.String())

		resp, err := httpClient.Get(config.Management_flights_api_url + "/" + v.FlightInfo.Id.String())
		if err != nil {
			handleError(w, err, "cannor get ticket flight with id")
			return
		}
		defer resp.Body.Close()

		var flight models.Flight
		err = json.NewDecoder(resp.Body).Decode(&flight)
		if err != nil {
			handleError(w, err, "cannot parse response body to flight object")
			return
		}
		///
		log.Println("OwnerId:", v.Owner.Id.String())
		resp, err = httpClient.Get(config.Management_customers_api_url + "/" + v.Owner.Id.String() + "/")
		if err != nil {
			handleError(w, err, "cannot ticket owner with id")
			return
		}
		defer resp.Body.Close()

		var customer models.Customer
		err = json.NewDecoder(resp.Body).Decode(&customer)
		if err != nil {
			handleError(w, err, "cannot parse response body to customer object")
			return
		}

		tickets[i].FlightInfo = flight
		tickets[i].Owner = customer
	}

	bs, err := json.Marshal(tickets)
	if err != nil {
		handleError(w, err, "cannot marshall ticket slice to byte string")
		return
	}

	message := models.Message{
		Key: config.KEY_FOR_TICKETS,
		Value: models.CacheItem{
			Value:      string(bs),
			Expiration: config.VALUE_EXPIRATION_TIME,
			SetedTime:  time.Now(),
		},
	}

	service.SaveToCache(message)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(bs)
	log.Println("Successfully fetched all tickets from management service.")
}

func GetTicketByIdHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Get ticket by id backend process...")

	id, ok := mux.Vars(r)["id"]
	if !ok {
		handleError(w, fmt.Errorf("bad request: cannot get id"), "Bad request")
		return
	}
	escapedID := url.PathEscape(config.KEY_FOR_TICKETS + ":" + id)
	rawURL := config.Management_cache_api_url + "/" + escapedID
	resp, err := httpClient.Get(rawURL)
	if err != nil {
		handleError(w, err, "Error fetching ticket by id from cache")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			handleError(w, err, "Error reading response body from cache")
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
		handleError(w, err, "Error fetching ticket by id from management service")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		handleError(w, fmt.Errorf("failed to fetch ticket from management service: %v", resp.StatusCode), "Failed to fetch ticket")
		return
	}

	var ticket models.Ticket
	err = json.NewDecoder(resp.Body).Decode(&ticket)
	if err != nil {
		handleError(w, err, "cannot parse response body to tickets slice")
		return
	}

	log.Println("FlightId:", ticket.FlightInfo.Id.String())
	resp, err = httpClient.Get(config.Management_flights_api_url + "/" + ticket.FlightInfo.Id.String())
	if err != nil {
		handleError(w, err, "cannor get ticket flight with id: "+ticket.FlightInfo.Id.String())
		return
	}
	defer resp.Body.Close()

	var flight models.Flight
	err = json.NewDecoder(resp.Body).Decode(&flight)
	if err != nil {
		handleError(w, err, "cannot parse response body to flight object")
		return
	}
	///
	log.Println("OwnerId:", ticket.Owner.Id.String())
	resp, err = httpClient.Get(config.Management_customers_api_url + "/" + ticket.Owner.Id.String() + "/")
	if err != nil {
		handleError(w, err, "cannot get ticket owner with id: "+ticket.Owner.Id.String())
		return
	}
	defer resp.Body.Close()

	var customer models.Customer
	err = json.NewDecoder(resp.Body).Decode(&customer)
	if err != nil {
		handleError(w, err, "cannot parse response body to customer object")
		return
	}

	bs, err := json.Marshal(ticket)
	if err != nil {
		handleError(w, err, "cannot marshal ticket to byte string")
		return
	}

	message := models.Message{
		Key: config.KEY_FOR_TICKETS,
		Value: models.CacheItem{
			Value:      string(bs),
			Expiration: config.VALUE_EXPIRATION_TIME,
			SetedTime:  time.Now(),
		},
	}

	service.SaveToCache(message)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(bs)
	log.Println("Successfully fetched ticket by id from management service.")
}

func InsertTicketHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Insert ticket backend process...")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		handleError(w, err, "Bad request: cannot read request body")
		return
	}
	defer r.Body.Close()

	resp, err := httpClient.Post(config.Management_tickets_api_url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		handleError(w, err, "Error inserting ticket into management service")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		handleError(w, fmt.Errorf("failed to insert ticket: %v", resp.StatusCode), "Failed to insert ticket")
		return
	}

	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		handleError(w, err, "Error reading response body from management service")
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

	id, ok := mux.Vars(r)["id"]
	if !ok {
		handleError(w, fmt.Errorf("bad request: cannot get id"), "Bad request")
		return
	}

	req, err := http.NewRequest(http.MethodPatch, config.Management_tickets_api_url+"/"+id, r.Body)
	if err != nil {
		handleError(w, err, "Cannot create request")
		return
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		handleError(w, err, "Error updating ticket in management service")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		handleError(w, fmt.Errorf("failed to update ticket: %v", resp.StatusCode), "Failed to update ticket")
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		handleError(w, err, "Error reading response body from management service")
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
		handleError(w, fmt.Errorf("bad request: cannot get id"), "Bad request")
		return
	}

	escapedID := url.PathEscape(id)
	rawURL := config.Management_tickets_api_url + "/" + escapedID
	req, err := http.NewRequest(http.MethodDelete, rawURL, nil)
	if err != nil {
		handleError(w, err, "Cannot create request")
		return
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		handleError(w, err, "Error deleting ticket from management service")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		handleError(w, fmt.Errorf("failed to delete ticket: %v", resp.StatusCode), "Failed to delete ticket")
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		handleError(w, err, "Error reading response body from management service")
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
		handleError(w, err, "Bad request: cannot read request body")
		return
	}
	defer r.Body.Close()

	resp, err := httpClient.Get(config.Management_customers_api_url + "/" + ticket.Owner.Id.String() + "/")
	if err != nil {
		handleError(w, err, "Error fetching customer data")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		handleError(w, fmt.Errorf("failed to fetch customer data: %v", resp.StatusCode), "Failed to fetch customer data")
		return
	}

	var customer models.Customer
	err = json.NewDecoder(resp.Body).Decode(&customer)
	if err != nil {
		handleError(w, err, "Error reading response body of customer")
		return
	}
	ticket.Owner = customer

	bs, err := json.Marshal(ticket)
	if err != nil {
		handleError(w, err, "Cannot marshal object to JSON")
		return
	}

	resp, err = httpClient.Post(config.Cart_api_url, "application/json", bytes.NewBuffer(bs))
	if err != nil {
		handleError(w, err, "Error inserting ticket to cart")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		handleError(w, fmt.Errorf("failed to insert ticket to cart: %v", resp.StatusCode), "Failed to insert ticket to cart")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	_, err = w.Write(bs)
	if err != nil {
		handleError(w, err, "Error writing response")
		return
	}
	log.Println("Successfully inserted ticket to cart.")

	go sendtoCache(bs)
}

func GetTicketsFromCartHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Fetching all tickets from cart backend process...")

	ownerId := r.URL.Query().Get("ownerId")
	log.Println("ownerId:", ownerId)

	resp, err := httpClient.Get(config.Cart_api_url + "/" + ownerId)
	if err != nil {
		handleError(w, err, "Error fetching tickets from cart")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		handleError(w, fmt.Errorf("failed to fetch tickets from cart: %v", resp.StatusCode), "Failed to fetch tickets")
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		handleError(w, err, "Error reading response body from cart")
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

func ClearTicketCartHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Clear ticket card backend process...")

	id := mux.Vars(r)["id"]

	req, err := http.NewRequest(http.MethodDelete, config.Cart_api_url+"/"+id+"/clear", nil)
	if err != nil {
		log.Println("Cannot create request")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		handleError(w, err, "Cannoot do request")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		handleError(w, fmt.Errorf("response code: %d", resp.StatusCode), "Response status: "+resp.Status)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		handleError(w, err, "Cannot read response body")
		return
	}
	log.Println("Cleared cart content for current user:", string(body))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
	log.Println("Successfully clear cart.")
}

func GetPurchasedTicketHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Fetching all purchased tickets backend process...")

	ownerId := r.URL.Query().Get("ownerId")
	if ownerId == "" {
		handleError(w, fmt.Errorf("bad request: ownerId is required"), "Bad request")
		return
	}

	resp, err := httpClient.Get(config.Purchased_tickets_api_url + "?ownerId=" + ownerId)
	if err != nil {
		handleError(w, err, "Error fetching purchased tickets")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		handleError(w, fmt.Errorf("error response code: %v", resp.StatusCode), "Error fetching purchased tickets")
		return
	}

	var tickets []models.Ticket
	err = json.NewDecoder(resp.Body).Decode(&tickets)
	if err != nil {
		handleError(w, err, "cannot parse response body to tickets slice")
		return
	}

	for i, v := range tickets {
		log.Println("FlightId:", v.FlightInfo.Id.String())

		resp, err := httpClient.Get(config.Management_flights_api_url + "/" + v.FlightInfo.Id.String())
		if err != nil {
			handleError(w, err, "cannor get ticket flight with id")
			return
		}
		defer resp.Body.Close()

		var flight models.Flight
		err = json.NewDecoder(resp.Body).Decode(&flight)
		if err != nil {
			handleError(w, err, "cannot parse response body to flight object")
			return
		}
		log.Println(flight.String())
		///
		log.Println("OwnerId:", v.Owner.Id.String())
		resp, err = httpClient.Get(config.Management_customers_api_url + "/" + v.Owner.Id.String() + "/")
		if err != nil {
			handleError(w, err, "cannot ticket owner with id")
			return
		}
		defer resp.Body.Close()

		var customer models.Customer
		err = json.NewDecoder(resp.Body).Decode(&customer)
		if err != nil {
			handleError(w, err, "cannot parse response body to customer object")
			return
		}
		log.Println(customer.String())

		tickets[i].FlightInfo = flight
		tickets[i].Owner = customer
	}

	bs, err := json.Marshal(tickets)
	if err != nil {
		handleError(w, err, "cannot marshall ticket slice to byte string")
		return
	}

	message := models.Message{
		Key: config.KEY_FOR_TICKETS,
		Value: models.CacheItem{
			Value:      string(bs),
			Expiration: config.VALUE_EXPIRATION_TIME,
			SetedTime:  time.Now(),
		},
	}

	service.SaveToCache(message)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bs)
	log.Println("Successfully fetched all purchased tickets.")
}

func PayForTickets(w http.ResponseWriter, r *http.Request) {
	log.Println("Pay for tickets backend process...")
	id := r.URL.Query().Get("ownerId")
	if id == "" {
		handleError(w, fmt.Errorf("id is empty string"), "Id is empty")
		return
	}

	var payData = struct {
		CardNumber  string `json:"card_number"`
		BankAccount string `json:"bank_account"`
		Cost        int    `json:"cost"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&payData)
	if err != nil {
		handleError(w, err, "Bad request: cannot parse body to cardInfo")
		return
	}

	payData.BankAccount = os.Getenv("BANK_ACCOUNT")
	body, err := json.Marshal(payData)
	if err != nil {
		handleError(w, err, "Cannot marshal ticket")
		return
	}

	resp, err := httpClient.Post(config.Payment_api_url+"/pay", "application/json", bytes.NewBuffer(body))
	if err != nil {
		handleError(w, err, "Error processing payment")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		handleError(w, fmt.Errorf("error response code: %v", resp.StatusCode), "Payment error")
		return
	}

	req, err := http.NewRequest(http.MethodDelete, config.Cart_api_url+"/"+id+"/clear", nil)
	if err != nil {
		handleError(w, err, "Cannot create request to clear cart")
		return
	}

	resp, err = httpClient.Do(req)
	if err != nil {
		handleError(w, err, "Error clearing cart")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		handleError(w, fmt.Errorf("error response code: %v", resp.StatusCode), "Error clearing cart")
		return
	}

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Bad request: cannot read response body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println("Trying clearing cart for current user:", string(body))

	resp, err = httpClient.Post(config.Management_tickets_api_url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		handleError(w, err, "Cannor do request")
		return
	}
	defer resp.Body.Close()
	log.Println("there1")

	if resp.StatusCode >= 400 {
		handleError(w, fmt.Errorf("response status: %s", resp.Status), "Response status: "+resp.Status)
		return
	}
	log.Println("there2")

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Cannot read response body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Println("there3")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
	log.Println("Successfully processed payment and cleared the cart.")
}

func handleError(w http.ResponseWriter, err error, context string) {
	log.Println(context, ":", err)
	if urlErr, ok := err.(*url.Error); ok {
		if urlErr.Timeout() {
			http.Error(w, "Request timed out", http.StatusGatewayTimeout)
		} else {
			http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		}
	} else {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
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
