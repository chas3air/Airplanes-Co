package routes

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/chas3air/Airplanes-Co/Core/Management_airplanes/internal/models"
	"github.com/chas3air/Airplanes-Co/Core/Management_airplanes/internal/service"
	"github.com/gorilla/mux"
)

//TODO: доделать

var dal_airplanes = os.Getenv("DAL_AIRPLANES_URL")
var limitTime = service.GetLimitTime("LIMIT_RESPONSE_TIME")
var httpClient = &http.Client{
	Timeout: limitTime,
}

func GetAirplanesHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Fetching all airplanes")

	resp, err := httpClient.Get(dal_airplanes + "/airplanes")
	if err != nil {
		log.Println("Cannot send request to:", dal_airplanes+"/airplanes")
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to fetch airplanes: %s\n", resp.Status)
		http.Error(w, "Failed to fetch airplanes", resp.StatusCode)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Bad request: cannot read response body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
	log.Println("Successfully fetched all airplanes.")
}

func GetAirplaneByIdHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Retrieving airplane by ID")

	id_s := mux.Vars(r)["id"]

	resp, err := httpClient.Get(dal_airplanes + "/airplanes")
	if err != nil {
		log.Println("Cannot send request to", dal_airplanes)
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		log.Printf("Failed to retrieve airplane by ID: %s\n", resp.Status)
		http.Error(w, "Airplane not found", resp.StatusCode)
		return
	}

	var airplanes []models.Airplane
	err = json.NewDecoder(resp.Body).Decode(&airplanes)
	if err != nil {
		log.Println("Cannot read response body")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var body []byte
	for _, airplane := range airplanes {
		if airplane.Id.String() == id_s {
			body, err = json.Marshal(airplane)
			if err != nil {
				log.Println("Cannot marshal airplane to json")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			break
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
	log.Println("Successfully retrieved airplane by ID.")
}

func InsertAirplaneHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Inserting airplane")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Bad request: cannot read request body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	resp, err := httpClient.Post(dal_airplanes+"/airplanes", "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Println("Error sending request")
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to insert airplane: %s\n", resp.Status)
		http.Error(w, "Failed to insert airplane", resp.StatusCode)
		return
	}

	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Cannot read response body")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(bs)
	log.Println("Successfully inserted airplane.")
}

func UpdateAirplaneHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Updating airplane")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Bad request: cannot read request body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	req, err := http.NewRequest(http.MethodPatch, dal_airplanes+"/airplanes", bytes.NewBuffer(body))
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

	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to update airplane: %s\n", resp.Status)
		http.Error(w, "Failed to update airplane", resp.StatusCode)
		return
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(responseBody)
	log.Println("Successfully updated airplane.")
}

func DeleteAirplaneHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Deleting airplane")

	id_s := mux.Vars(r)["id"]

	req, err := http.NewRequest(http.MethodDelete, dal_airplanes+"/airplanes/"+id_s, nil) // Исправлено: добавлен id к URL
	if err != nil {
		log.Println("Error creating request")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Println("Error sending request")
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to delete airplane: %s\n", resp.Status)
		http.Error(w, "Failed to delete airplane", resp.StatusCode)
		return
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Cannot read response body")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(responseBody)
	log.Println("Successfully deleted airplane.")
}
