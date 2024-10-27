package router

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/chas3air/Airplanes-Co/Management_customers/internal/service"
	"github.com/gorilla/mux"
)

var database_url = os.Getenv("DAL_CUSTOMERS_URL")
var limitTime = service.GetLimitTime()

var httpClient = &http.Client{
	Timeout: limitTime,
}

// GetAllCustomersHandler handles a GET request to fetch all customers.
// It retrieves all customer data from the database and returns it in JSON format.
// If an error occurs during the process, it responds with an appropriate HTTP status code.
func GetCustomersHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Fetching all customers")

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
	log.Println("Successfully fetched all customers.")
}

// GetCustomerByIdHandler handles a GET request to retrieve a customer by ID.
// It takes the customer ID from the URL and returns the customer data in JSON format.
// If the customer is not found or an error occurs, it responds with an appropriate HTTP status code.
func GetCustomerByIdHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Retrieving customer by ID")

	id_s := mux.Vars(r)["id"]

	resp, err := httpClient.Get(database_url + "/" + id_s + "/")
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
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
	log.Println("Successfully retrieved customer by ID.")
}

// GetCustomerByLoginAndPasswordHandler handles a GET request to retrieve a customer by login and password.
// It expects the login and password to be sent as query parameters.
// If successful, it returns the customer data in JSON format; otherwise, it responds with an error.
func GetCustomerByLoginAndPasswordHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Retrieving customer by login and password")

	r.ParseForm()
	login := r.Form.Get("login")
	password := r.Form.Get("password")

	url := fmt.Sprintf("%s/login?login=%s&password=%s", database_url, login, password)

	resp, err := httpClient.Get(url)
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
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
	log.Println("Successfully signed in.")
}

// InsertCustomerHandler handles a POST request to add a new customer.
// It expects customer data in JSON format in the request body.
// If successful, it responds with the inserted customer data; otherwise, it responds with an error.
func InsertCustomerHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Inserting customer")

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
	log.Println("Successfully inserted customer.")
}

// UpdateCustomerHandler handles a PATCH request to update customer information.
// It expects the updated customer data in JSON format in the request body.
// If successful, it responds with the updated customer data; otherwise, it responds with an error.
func UpdateCustomerHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Updating customer")

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
	log.Println("Successfully updated customer.")
}

// DeleteCustomerHandler handles a DELETE request to remove a customer by ID.
// It takes the customer ID from the URL and responds with confirmation of the deletion.
// If an error occurs, it responds with an appropriate HTTP status code.
func DeleteCustomerHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Deleting customer")

	id_s := mux.Vars(r)["id"]

	req, err := http.NewRequest(http.MethodDelete, database_url+"/"+id_s+"/", nil)
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
	log.Println("Successfully deleted customer.")
}
