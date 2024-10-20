package auth

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

var env_url = os.Getenv("DAL_CUSTOMERS_URL")

// SignInHandler handles user sign-in requests.
// It retrieves login and password from the request, and validates them against the customer database.
func SignInHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Initiating sign-in process")

	r.ParseForm()
	login := r.Form.Get("login")
	password := r.Form.Get("password")

	url := fmt.Sprintf("%s/get/lp?login=%s&password=%s", env_url, login, password)

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error during sign-in: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Cannot read response body during sign-in")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(bs)
	log.Println("Successfully signed in.")
}

// SignUpHandler handles user sign-up requests.
// It reads user data from the request body and inserts it into the customer database.
func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Initiating sign-up process")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Bad request: cannot read request body during sign-up")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	resp, err := http.Post(env_url+"/insert", "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Error sending request during sign-up: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Sign-up error, status code: %d", resp.StatusCode)
		http.Error(w, fmt.Sprintf("error, status: %d", resp.StatusCode), resp.StatusCode)
		return
	}

	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Cannot read response body during sign-up")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bs)
	log.Println("Successfully signed up.")
}
