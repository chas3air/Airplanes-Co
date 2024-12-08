package customers_routes

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

func GetCustomersHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Get customers backend process...")

	resp, err := httpClient.Get(config.Management_cache_api_url + "/" + config.KEY_FOR_CUSTOMERS)
	if err != nil {
		handleError(w, err, "Error fetching customers from cache")
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

		log.Println("Backend sending to client:", string(body))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.StatusCode)
		w.Write(body)
		log.Println("Successfully fetched customers from cache.")
		return
	}

	log.Println("Cache is not used")

	resp, err = httpClient.Get(config.Management_customers_api_url)
	if err != nil {
		handleError(w, err, "Error fetching customers from management service")
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		handleError(w, err, "Error reading response body from management service")
		return
	}

	message := models.Message{
		Key: config.KEY_FOR_CUSTOMERS,
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
	log.Println("Successfully fetched all customers from management service.")
}

func GetCustomerByIdHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Get customer by id backend process...")

	id, ok := mux.Vars(r)["id"]
	if !ok {
		http.Error(w, "Bad request: cannot get id", http.StatusBadRequest)
		return
	}

	escapedID := url.PathEscape(config.KEY_FOR_CUSTOMERS + ":" + id)
	rawURL := config.Management_cache_api_url + "/" + escapedID
	resp, err := httpClient.Get(rawURL)
	if err != nil {
		handleError(w, err, "Error fetching customer by id from cache")
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
		log.Println("Successfully fetched customer by id from cache.")
		return
	}

	resp, err = httpClient.Get(config.Management_customers_api_url + "/" + id + "/")
	if err != nil {
		handleError(w, err, "Error fetching customer by id from management service")
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		handleError(w, err, "Error reading response body from management service")
		return
	}

	message := models.Message{
		Key: config.KEY_FOR_CUSTOMERS + ":" + id,
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
	log.Println("Successfully fetched customer by id:", id)
}

func InsertCustomerHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Insert customer backend process...")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad request: cannot read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	resp, err := httpClient.Post(config.Management_customers_api_url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		handleError(w, err, "Error inserting customer into management service")
		return
	}
	defer resp.Body.Close()

	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		handleError(w, err, "Error reading response body from management service")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(bs)
	log.Println("Successfully inserted customer.")
}

func UpdateCustomerHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Update customer backend process...")

	req, err := http.NewRequest(http.MethodPatch, config.Management_customers_api_url, r.Body)
	if err != nil {
		log.Println("Cannot create request:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		handleError(w, err, "Error updating customer into management service")
		return
	}
	defer resp.Body.Close()

	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		handleError(w, err, "Error reading response body from management service")
		return
	}

	if resp.StatusCode >= 400 {
		log.Println("Response code:", resp.StatusCode)
		http.Error(w, fmt.Errorf("response code: %v", resp.StatusCode).Error(), resp.StatusCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(bs)
	log.Println("Successfully updated customer.")
}

func DeleteCustomerHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Delete customer backend process...")

	id, ok := mux.Vars(r)["id"]
	if !ok {
		http.Error(w, "Bad request: cannot get id", http.StatusBadRequest)
		return
	}

	escapedID := url.PathEscape(id)
	rawURL := config.Management_customers_api_url + "/" + escapedID
	req, err := http.NewRequest(http.MethodDelete, rawURL, nil)
	if err != nil {
		handleError(w, err, "Error creating delete request")
		return
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		handleError(w, err, "Error deleting customer from management service")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		handleError(w, fmt.Errorf("unexpected response code: %v", resp.StatusCode), "Failed to delete customer")
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
	log.Println("Successfully deleted customer.")
}

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Sign-up customer backend process...")

	resp, err := httpClient.Post(config.Auth_api_url+"/auth/signup", "application/json", r.Body)
	if err != nil {
		handleError(w, err, "Error signing up customer")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		handleError(w, fmt.Errorf("unexpected response code: %v", resp.StatusCode), "Sign up failed")
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		handleError(w, err, "Error reading response body from sign up")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
	log.Println("Successfully signed up customer.")
}

func SignInHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Sign in customer backend process...")

	r.ParseForm()
	login := r.Form.Get("login")
	password := r.Form.Get("password")
	log.Println("Received login and password from URL:", login, "and", password)

	url := fmt.Sprintf("%s/auth/signin?login=%s&password=%s", config.Auth_api_url, login, password)
	resp, err := httpClient.Get(url)
	if err != nil {
		handleError(w, err, "Error signing in customer")
		return
	}
	defer resp.Body.Close()
	log.Println("Response code:", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		handleError(w, fmt.Errorf("unexpected response code: %v", resp.StatusCode), "Sign in failed")
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		handleError(w, err, "Error reading response body from sign in")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
	log.Println("Successfully signed in customer.")
}

func handleError(w http.ResponseWriter, err error, context string) {
	log.Println(context, ":", err)
	if urlErr, ok := err.(*url.Error); ok && urlErr.Timeout() {
		http.Error(w, "Request timed out", http.StatusGatewayTimeout)
	} else {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
