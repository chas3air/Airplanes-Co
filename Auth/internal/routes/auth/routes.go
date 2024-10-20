package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/chas3air/Airplanes-Co/Auth/internal/models"
)

func SignInHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Sign in process")

	env_url := os.Getenv("DATABASE_URL")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Bad request: cannot read request body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var pre_cust models.PreCust
	err = json.Unmarshal(body, &pre_cust)
	if err != nil {
		log.Println("Cannot unmarshal request body to models.PreCust")
		log.Println("Body:", string(body))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	url := fmt.Sprintf(env_url+"/get/lp"+"?login=%s&password=%s", pre_cust.Login, pre_cust.Password)

	resp, err := http.Get(url)
	if err != nil {
		log.Println("Error:", err.Error())
		http.Error(w, err.Error(), resp.StatusCode)
		return
	}
	defer resp.Body.Close()

	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Cannot read responce body")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bs)
	log.Println("Successfully sign in.")
}

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Sign up process")

	env_url := os.Getenv("DATABASE_URL")

	// тут рабоатет модель models.Customer

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Bad request: cannot read request body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	resp, err := http.Post(env_url+"/insert", "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Println("Error sending request")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("Error, status:", resp.StatusCode)
		http.Error(w, fmt.Errorf("error, status: %d", resp.StatusCode).Error(), resp.StatusCode)
		return
	}

	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Cannot read response body")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bs)
	log.Println("Successfully sing up.")
}
