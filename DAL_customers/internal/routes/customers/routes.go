package customers

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/chas3air/Airplanes-Co/DAL_customers/internal/models"
	"github.com/chas3air/Airplanes-Co/DAL_customers/internal/service"
	"github.com/chas3air/Airplanes-Co/DAL_customers/internal/storage"
	"github.com/gorilla/mux"
)

var CustomersDB = storage.MustGetInstanceOfCustomersStorage("psql")
var limitTime = service.GetLimitTime()

func GetCustomers(w http.ResponseWriter, r *http.Request) {
	log.Println("Fetching all customers")

	ctx, cancel := context.WithTimeout(context.Background(), limitTime)
	defer cancel()

	entities, err := CustomersDB.GetAll(ctx)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("Request timed out")
			http.Error(w, "Request timed out", http.StatusGatewayTimeout)
			return
		} else {
			log.Printf("Error retrieving customers: %s\n", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	customers, ok := entities.([]models.Customer)
	if !ok {
		log.Println("Invalid data type")
		http.Error(w, "Invalid data type", http.StatusInternalServerError)
		return
	}

	bs, err := json.Marshal(customers)
	if err != nil {
		log.Printf("Cannot marshal object: %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bs)
	log.Println("Successfully fetched customers.")
}

func GetCustomerById(w http.ResponseWriter, r *http.Request) {
	log.Println("Fetching customer by ID")

	ctx, cancel := context.WithTimeout(context.Background(), limitTime)
	defer cancel()

	id := mux.Vars(r)["id"]

	entity, err := CustomersDB.GetById(ctx, id)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("Request timed out")
			http.Error(w, "Request timed out", http.StatusGatewayTimeout)
			return
		} else {
			log.Printf("Error retrieving customer by ID: %s\n", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	customer, ok := entity.(models.Customer)
	if !ok {
		log.Println("Invalid data type")
		http.Error(w, "Invalid data type", http.StatusInternalServerError)
		return
	}

	bs, err := json.Marshal(customer)
	if err != nil {
		log.Printf("Cannot marshal object: %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bs)
	log.Println("Successfully fetched customer by ID.")
}

func GetCustomerByLoginAndPassword(w http.ResponseWriter, r *http.Request) {
	log.Println("Fetching customer by login and password")

	r.ParseForm()
	login := r.Form.Get("login")
	password := r.Form.Get("password")

	ctx, cancel := context.WithTimeout(context.Background(), limitTime)
	defer cancel()

	entity, err := CustomersDB.GetByLoginAndPassword(ctx, login, password)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("Request timed out")
			http.Error(w, "Request timed out", http.StatusGatewayTimeout)
			return
		} else {
			log.Printf("Error retrieving customer by ID: %s\n", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	customer, ok := entity.(models.Customer)
	if !ok {
		log.Println("Invalid data type")
		http.Error(w, "Invalid data type", http.StatusInternalServerError)
		return
	}

	bs, err := json.Marshal(customer)
	if err != nil {
		log.Printf("Cannot marshal object: %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bs)
	log.Println("Successfully fetched customer by ID.")
}

func InsertCustomer(w http.ResponseWriter, r *http.Request) {
	log.Println("Inserting customer")

	ctx, cancel := context.WithTimeout(context.Background(), limitTime)
	defer cancel()

	bs, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Bad request: cannot read request body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var customer models.Customer
	err = json.Unmarshal(bs, &customer)
	if err != nil {
		log.Println("Cannot unmarshal request body to customer")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	obj, err := CustomersDB.Insert(ctx, customer)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("Request timed out")
			http.Error(w, "Request timed out", http.StatusGatewayTimeout)
			return
		} else {
			log.Printf("Error inserting customer with ID: %d\n", customer.Id)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	out_bs, err := json.Marshal(obj)
	if err != nil {
		log.Printf("Cannot marshal inserted customer: %s\n", err.Error())
		http.Error(w, "Cannot marshal object", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out_bs)
	log.Printf("Successfully inserted customer with ID: %d\n", customer.Id)
}

func UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	log.Println("Updating customer")

	ctx, cancel := context.WithTimeout(context.Background(), limitTime)
	defer cancel()

	bs, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Bad request: cannot read request body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var customer models.Customer
	err = json.Unmarshal(bs, &customer)
	if err != nil {
		log.Printf("Cannot unmarshal request body to customer: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	obj, err := CustomersDB.Update(ctx, customer)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("Request timed out")
			http.Error(w, "Request timed out", http.StatusGatewayTimeout)
			return
		} else {
			log.Printf("Error updating customer with ID %d: %s", customer.Id, err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	out_bs, err := json.Marshal(obj)
	if err != nil {
		log.Println("Cannot marshal updated customer")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out_bs)
	log.Printf("Successfully updated customer with ID: %d", customer.Id)
}

func DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	log.Println("Deleting customer")

	ctx, cancel := context.WithTimeout(context.Background(), limitTime)
	defer cancel()

	id := mux.Vars(r)["id"]

	obj, err := CustomersDB.Delete(ctx, id)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("Request timed out")
			http.Error(w, "Request timed out", http.StatusGatewayTimeout)
			return
		} else {
			log.Printf("Error deleting customer with ID %v: %s", id, err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	out_bs, err := json.Marshal(obj)
	if err != nil {
		log.Println("Cannot marshal deleted customer")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out_bs)
	log.Printf("Successfully deleted customer with ID: %v", id)
}
