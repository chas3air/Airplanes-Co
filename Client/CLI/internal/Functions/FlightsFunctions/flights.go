package flightsfunctions

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/chas3air/Airplanes-Co/Client/CLI/internal/config"
	"github.com/chas3air/Airplanes-Co/Client/CLI/internal/models"
	"github.com/chas3air/Airplanes-Co/Client/CLI/internal/service"
	"github.com/google/uuid"
)

var limitTime = service.GetLimitTime()

// /flights/get
// GetAllFlights retrieves a list of all flights from the API.
// Returns an array of flights and an error if the request fails.
func GetAllFlights() ([]models.Flight, error) {
	resp, err := http.Get(config.Backend_url + "/flights/get")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get flights: %s", resp.Status)
	}

	var flights []models.Flight
	err = json.NewDecoder(resp.Body).Decode(&flights)
	if err != nil {
		return nil, err
	}

	return flights, nil
}

// /flihgts/get/id
// GetFlightById retrieves a flight by its ID from the backend service.
// Returns the flight details or an error if the request fails or the flight is not found.
func GetFlightById(id string) (models.Flight, error) {
	resp, err := http.Get(config.Backend_url + "/flights/get/" + id)
	if err != nil {
		log.Println("Error:", err)
		return models.Flight{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return models.Flight{}, fmt.Errorf("failed to get flight: %s", resp.Status)
	}

	var flight models.Flight
	err = json.NewDecoder(resp.Body).Decode(&flight)
	if err != nil {
		return models.Flight{}, err
	}

	return flight, nil
}

// /flights/insert + body(json)
// postFlight sends the new flight data to the server for addition.
// Returns the added flight and an error if there was a problem.
func InsertFlight(flight models.Flight) (models.Flight, error) {
	bs, err := json.Marshal(flight)
	if err != nil {
		fmt.Println("Error marshaling flight:", err)
		return models.Flight{}, err
	}

	client := &http.Client{
		Timeout: limitTime,
	}

	resp, err := client.Post(config.Backend_url+"/flights/insert", "application/json", bytes.NewBuffer(bs))
	if err != nil {
		return models.Flight{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return models.Flight{}, fmt.Errorf("failed to post flight: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Cannot read response body:", err)
		return models.Flight{}, err
	}

	var outFlight models.Flight
	if err = json.Unmarshal(body, &outFlight); err != nil {
		fmt.Println("Cannot unmarshal response body:", err)
		return models.Flight{}, err
	}

	return outFlight, nil
}

// /flights/update + body(json)
// UpdateFlight updates the data of an existing flight by its UUID.
// Returns the updated flight and an error if there was a problem.
func UpdateFlight(flight models.Flight) (models.Flight, error) {
	bs, err := json.Marshal(flight)
	if err != nil {
		fmt.Println("Error marshaling flight:", err)
		return models.Flight{}, err
	}

	req, err := http.NewRequest(http.MethodPatch, config.Backend_url+"/flights/update", bytes.NewBuffer(bs))
	if err != nil {
		return models.Flight{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: limitTime,
	}

	resp, err := client.Do(req)
	if err != nil {
		return models.Flight{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return models.Flight{}, fmt.Errorf("failed to patch flight: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Cannot read response body:", err)
		return models.Flight{}, err
	}

	var outFlight models.Flight
	if err = json.Unmarshal(body, &outFlight); err != nil {
		fmt.Println("Cannot unmarshal response body:", err)
		return models.Flight{}, err
	}

	return outFlight, nil
}

// /flights/delete/id
// DeleteFlight removes a flight from the system by its ID.
// Returns the deleted flight and an error if there was a problem.
func DeleteFlight(id string) (models.Flight, error) {
	req, err := http.NewRequest(http.MethodDelete, config.Backend_url+"/flights/delete/"+id, nil)
	if err != nil {
		return models.Flight{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: limitTime,
	}

	resp, err := client.Do(req)
	if err != nil {
		return models.Flight{}, err
	}
	defer resp.Body.Close()

	resp_body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.Flight{}, err
	}

	var flight models.Flight
	err = json.Unmarshal(resp_body, &flight)
	if err != nil {
		return models.Flight{}, err
	}
	return flight, nil
}

func PrintFlights(flights []models.Flight, exclude ...string) {
	excludeMap := make(map[string]struct{})
	for _, col := range exclude {
		excludeMap[col] = struct{}{}
	}

	headersNames := [6]string{"ID", "From", "Destination", "FlightTime", "Duration", "Costs"}
	headersWidth := [6]int{36, 15, 15, 19, 15, 17}

	counter := 0
	for _, header := range headersNames {
		fmt.Printf("| %-*s ", headersWidth[counter], header)
		counter++
	}
	fmt.Println("|")
	fmt.Println(strings.Repeat("-", 136))

	for _, flight := range flights {
		if _, ok := excludeMap["ID"]; ok {
			fmt.Printf("| %-36s ", "None")
		} else {
			fmt.Printf("| %-36s ", flight.Id)
		}

		if _, ok := excludeMap["From"]; ok {
			fmt.Printf("| %-15s ", "None")
		} else {
			fmt.Printf("| %-15s ", flight.FromWhere)
		}

		if _, ok := excludeMap["Destination"]; ok {
			fmt.Printf("| %-15s ", "None")
		} else {
			fmt.Printf("| %-15s ", flight.Destination)
		}

		if _, ok := excludeMap["FlightTime"]; ok {
			fmt.Printf("| %-19s ", "None")
		} else {
			flightTime := flight.FlightTime.Format("2006-01-02 15:04:05")
			fmt.Printf("| %-19s ", flightTime)
		}

		if _, ok := excludeMap["Duration"]; ok {
			fmt.Printf("| %-15s ", "None")
		} else {
			fmt.Printf("| %-15d ", flight.FlightDuration)
		}

		if _, ok := excludeMap["Costs"]; ok {
			fmt.Printf("| %-17s ", "None")
		} else {
			costs := fmt.Sprintf("%v", flight.FlightSeatsCosts)
			fmt.Printf("| %-17s ", costs)
		}

		fmt.Println("|")
	}

	fmt.Println(strings.Repeat("-", 133))
}

func CreateFlight() (models.Flight, error) {
	scanner := bufio.NewScanner(os.Stdin)

	fromWhere := service.GetInput(scanner, "Enter from where will arrive plane")
	destination := service.GetInput(scanner, "Enter destination")
	flightTimeStr := service.GetInput(scanner, "Enter flight time (format: YYYY-MM-DD HH:MM:SS)")

	flightTime, err := service.ParseTime(flightTimeStr, "2006-01-02 15:04:05")
	if err != nil {
		return models.Flight{}, err
	}

	flightDuration, err := service.GetInt(scanner, "Enter flight duration")
	if err != nil {
		return models.Flight{}, err
	}

	flight := models.Flight{
		Id:               uuid.New(),
		FromWhere:        fromWhere,
		Destination:      destination,
		FlightTime:       flightTime,
		FlightDuration:   flightDuration,
		FlightSeatsCosts: make([]int, 0),
	}

	log.Println("Id of flight generated on client:", flight.Id)

	return flight, nil
}
