package catalog

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/chas3air/Airplanes-Co/Flights_catalog/internal/service"
	"github.com/gorilla/mux"
)

/*
	DAL_FLIGHTS_URL         = "http://dal_flights:12001/dal/flights/postgres/get"
	DAL_LIMIT_RESPONSE_TIME = 5
*/

var dal_flights_url = os.Getenv("DAL_FLIGHTS_URL")

func GetFlightsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Fetching flights")

	limitTime, err := service.GetLimitTime()
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	client := http.Client{
		Timeout: limitTime,
	}

	response, err := client.Get(dal_flights_url)
	if err != nil {
		log.Printf("Error of connecting to %s: %s\n", dal_flights_url, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Printf("Status code is not ok, code: %d\n", response.StatusCode)
		http.Error(w, fmt.Sprintf("Status code is not ok, code: %d", response.StatusCode), response.StatusCode)
		return
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println("Cannot read response body:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
	log.Println("Successfully fetched flights")
}

func GetFlightByIDHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Fetching flight by id")

	limitTime, err := service.GetLimitTime()
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id_s := mux.Vars(r)["id"]

	client := http.Client{
		Timeout: limitTime,
	}

	response, err := client.Get(dal_flights_url + "/" + id_s)
	if err != nil {
		fmt.Printf("Error of connecting to %s: %s\n", dal_flights_url, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Println("Status code is not ok, code:", response.StatusCode)
		http.Error(w, "Status code is not available", response.StatusCode)
		return
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println("Cannot read response body")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
	log.Println("Successfully fetched flight by id.")
}

func SearchFlightsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Search flights")

	limitTime, err := service.GetLimitTime()
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	client := http.Client{
		Timeout: limitTime,
	}

	response, err := client.Get(dal_flights_url)
	if err != nil {
		log.Printf("Error of connecting to %s: %s\n", dal_flights_url, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Printf("Status code is not ok, code: %d\n", response.StatusCode)
		http.Error(w, fmt.Sprintf("Received status code: %d", response.StatusCode), response.StatusCode)
		return
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println("Cannot read response body:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = body
}

func StatsFlightsHandler(w http.ResponseWriter, r *http.Request) {

}
