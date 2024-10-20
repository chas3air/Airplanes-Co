package router

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

var database_url = os.Getenv("DAL_CUSTOMERS_URL")

// GetAllCustomersHandler handles a GET request to fetch all customers.
// Returns a list of customers in JSON format.
func GetAllCustomersHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Fetching all customers")

	resp, err := http.Get(database_url + "/get")
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
// Takes the customer ID from the URL and returns the customer data in JSON format.
func GetCustomerByIdHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Retrieving customer by ID")

	id_s := mux.Vars(r)["id"]

	resp, err := http.Get(database_url + "/get/" + id_s)
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
func GetCustomerByLoginAndPasswordHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Retrieving customer by login and password")

	r.ParseForm()
	login := r.Form.Get("login")
	password := r.Form.Get("password")

	url := fmt.Sprintf("%s/get/lp?login=%s&password=%s", database_url, login, password)

	resp, err := http.Get(url)
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
// Takes customer data in JSON format and returns confirmation of the insertion.
func InsertCustomerHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Inserting customer")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Bad request: cannot read request body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	resp, err := http.Post(database_url+"/insert", "application/json", bytes.NewBuffer(body))
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
// Takes customer data in JSON format and returns the updated customer data.
func UpdateCustomerHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Updating customer")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Bad request: cannot read request body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	req, err := http.NewRequest(http.MethodPatch, database_url+"/update", bytes.NewBuffer(body))
	if err != nil {
		log.Println("Error creating request")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
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
// Takes the customer ID from the URL and returns confirmation of the deletion.
func DeleteCustomerHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Deleting customer")

	id_s := mux.Vars(r)["id"]

	req, err := http.NewRequest(http.MethodDelete, database_url+"/delete/"+id_s, nil)
	if err != nil {
		log.Println("Error creating request")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
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
