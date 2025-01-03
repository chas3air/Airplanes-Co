package router

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/chas3air/Airplanes-Co/Core/Auth/internal/service"
)

var env_url = os.Getenv("MANAGEMENT_CUSTOMERS_URL")
var limitTime = service.GetLimitTime("LIMIT_RESPONSE_TIME")

var httpClient = &http.Client{
	Timeout: limitTime,
}

// SignInHandler handles user sign-in requests.
// It retrieves the login and password from the request, validates them against the customer database,
// and returns the response from the authentication service.
func SignInHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Initiating sign-in process")

	r.ParseForm()
	login := r.Form.Get("login")
	password := r.Form.Get("password")

	url := fmt.Sprintf("%s/login?login=%s&password=%s", env_url, login, password)
	log.Println("url:", url)

	resp, err := httpClient.Get(url)
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
// It reads user data from the request body, sends it to the customer database,
// and returns the response from the registration service.
func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Initiating sign-up process")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Bad request: cannot read request body during sign-up")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	resp, err := httpClient.Post(env_url, "application/json", bytes.NewBuffer(body))
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
