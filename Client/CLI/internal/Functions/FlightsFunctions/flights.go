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
	"text/tabwriter"

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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Cannot read response body:", err)
		return nil, err
	}

	log.Println("Body response from backend is:", string(body))

	bodyStr := string(body)

	var flights []models.Flight
	err = json.Unmarshal([]byte(bodyStr), &flights)
	if err != nil {
		log.Println("Cannot unmarshall response body to flight:", err)
		return nil, err
	}

	return flights, nil
}

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
		log.Println(flight)
		return models.Flight{}, err
	}

	return flight, nil
}

func CreateFlight() (models.Flight, error) {
	scanner := bufio.NewScanner(os.Stdin)

	id := uuid.New()
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
		Id:             id,
		FromWhere:      fromWhere,
		Destination:    destination,
		FlightTime:     flightTime,
		FlightDuration: flightDuration,
	}

	return flight, nil
}

// /flights/delete/id
// DeleteFlight removes a flight from the system by its ID.
// Returns the deleted flight and an error if there was a problem.
func DeleteFlight(id string) (models.Flight, error) {
	req, err := http.NewRequest(http.MethodDelete, config.Backend_url+"/flights/delete/"+id, nil)
	if err != nil {
		return models.Flight{}, nil
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

func PrintFlights(flights []models.Flight) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)

	fmt.Fprintln(w, "ID\tFrom\tDestination\tFlight Time\tDuration\tCosts")
	fmt.Fprintln(w, "-----------------------------------------------------------")

	for _, flight := range flights {
		costs := fmt.Sprintf("%v", flight.FlightSeatsCosts)
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%d\t%s\n",
			flight.Id,
			flight.FromWhere,
			flight.Destination,
			flight.FlightTime.Format("2006-01-02 15:04:05"),
			flight.FlightDuration,
			costs)
	}

	w.Flush()
}
