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

		log.Println("backend send to client", string(body))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.StatusCode)
		w.Write(body)
		log.Println("Successfully fetched customers from cache.")
		return
	}

	log.Println("Cache is not used")

	resp, err = httpClient.Get(config.Management_customers_api_url)
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
		log.Println("Bad request: cannot get id")
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	escapeID := url.PathEscape(config.KEY_FOR_CUSTOMERS + ":" + id)
	rawURL := config.Management_cache_api_url + "/" + escapeID
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
		log.Println("Successfully fetched customer by id from cache.")
		return
	}

	resp, err = httpClient.Get(config.Management_customers_api_url + "/" + id + "/")
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
		log.Println("Bad request: cannot read request body")
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	resp, err := httpClient.Post(config.Management_customers_api_url, "application/json", bytes.NewBuffer(body))
	if err != nil {
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
	log.Println("Successfully inserted customer.")
}

func UpdateCustomerHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func DeleteCustomerHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Delete customer backend process...")

	id, ok := mux.Vars(r)["id"]
	if !ok {
		log.Println("Bad request: cannot get id")
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	escapeID := url.PathEscape(id)
	rawURL := config.Management_customers_api_url + "/" + escapeID
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

	if resp.StatusCode != http.StatusOK {
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
	log.Println("Successfully deleted customer.")
}

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Sign-up customer backend process...")

	resp, err := httpClient.Post(config.Auth_api_url+"/auth/signup", "application/json", r.Body)
	if err != nil {
		log.Println("Error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		log.Println("Bad response: cannot read response body")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
	log.Println("Successfully deleted customer.")
}

func SignInHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Sign in customer backend process...")

	r.ParseForm()
	login := r.Form.Get("login")
	password := r.Form.Get("password")
	log.Println("Received login and password from url:", login, "and", password)

	url := fmt.Sprintf("%s/auth/signin?login=%s&password=%s", config.Auth_api_url, login, password)
	resp, err := httpClient.Get(url)
	if err != nil {
		log.Println("Error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	log.Println("Response code:", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		http.Error(w, fmt.Errorf("response code: %v", resp.StatusCode).Error(), resp.StatusCode)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
	log.Println("Successfully sign in")
}

func handleError(w http.ResponseWriter, err error) {
	if urlErr, ok := err.(*url.Error); ok && urlErr.Timeout() {
		log.Println("Request timed out:", urlErr)
		http.Error(w, "Request timed out", http.StatusGatewayTimeout)
	} else {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
