package fligthtadmininterface

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	flightsfunctions "github.com/chas3air/Airplanes-Co/Client/CLI/internal/Functions/FlightsFunctions"
	"github.com/chas3air/Airplanes-Co/Client/CLI/internal/models"
	"github.com/chas3air/Airplanes-Co/Client/CLI/internal/service"
)

func FlightsAdminInterface(user *models.Customer) {
	scanner := bufio.NewScanner(os.Stdin)
	var localFlights []models.Flight
	var prepFlightToInsert = make([]models.Flight, 0, 5)
	var prepFlightToUpdate = make([]models.Flight, 0, 5)
	var prepIdFlightToDelete = make([]string, 0, 10)
	var err error
	var mut sync.Mutex

	localFlights, err = flightsfunctions.GetAllFlights()
	if err != nil {
		fmt.Println("Flights weren't loaded:", err)
		return
	}

	go func() {
		for {
			mut.Lock()
			processPreparations(&prepFlightToInsert, &prepFlightToUpdate, &prepIdFlightToDelete)
			mut.Unlock()
			time.Sleep(3 * time.Second)
		}
	}()

	for {
		clearConsole()
		displayMenu()
		_ = scanner.Scan()
		choice := scanner.Text()

		switch choice {
		case "1":
			fmt.Println("Show all flights")
			displayFlights(localFlights)

		case "2":
			fmt.Println("Adding flights")
			flight := models.Flight{}
			if err := getFlightInput(scanner, &flight); err != nil {
				fmt.Println("Error adding flight:", err)
				break
			}
			localFlights = append(localFlights, flight)
			prepFlightToInsert = append(prepFlightToInsert, flight)

		case "3":
			fmt.Println("Update flight")
			displayFlights(localFlights)

			id, err := service.GetInt(scanner, "Enter id (1-based index)")
			if err != nil || id < 1 || id > len(localFlights) {
				fmt.Println("Undefined element")
				break
			}

			flight, err := flightsfunctions.CreateFlight()
			if err != nil {
				fmt.Println("Error creating flight:", err)
				break
			}

			prepFlightToUpdate = append(prepFlightToUpdate, flight)
			localFlights[id-1] = flight
			fmt.Println("Updated flight:", flight.String())

		case "4":
			fmt.Println("Deleting flight")
			displayFlights(localFlights)

			id, err := service.GetInt(scanner, "Enter id (1-based index)")
			if err != nil || id < 1 || id > len(localFlights) {
				fmt.Println("Undefined element")
				break
			}

			prepIdFlightToDelete = append(prepIdFlightToDelete, localFlights[id-1].Id.String())

		case "5":
			fmt.Println("Logout")
			if err := service.Logout(); err != nil {
				fmt.Println("Error logging out:", err)
				break
			}
			*user = models.Customer{}
			return

		default:
			fmt.Println("Error: invalid menu option")
		}
		time.Sleep(200 * time.Millisecond)
	}
}

func processPreparations(prepFlightToInsert, prepFlightToUpdate *[]models.Flight, prepIdFlightToDelete *[]string) {
	for _, v := range *prepFlightToInsert {
		if _, err := flightsfunctions.PostFlight(v); err != nil {
			log.Println("Error posting flight:", err)
		}
	}

	for _, v := range *prepFlightToUpdate {
		if _, err := flightsfunctions.UpdateFlight(v); err != nil {
			log.Println("Error updating flight:", err)
		}
	}

	for _, v := range *prepIdFlightToDelete {
		if _, err := flightsfunctions.DeleteFlight(v); err != nil {
			log.Println("Error deleting flight:", err)
		}
	}

	*prepFlightToInsert = make([]models.Flight, 0, 5)
	*prepFlightToUpdate = make([]models.Flight, 0, 5)
	*prepIdFlightToDelete = make([]string, 0, 10)
}

func displayFlights(flights []models.Flight) {
	fmt.Println("Current Flights:")
	for i, flight := range flights {
		fmt.Printf("%d: %s\n", i+1, flight.String())
	}
	fmt.Println("Press Enter to continue...")
	bufio.NewReader(os.Stdin).ReadString('\n')
}

func displayMenu() {
	fmt.Println("Select an item")
	fmt.Println("1. Get all flights")
	fmt.Println("2. Add flight")
	fmt.Println("3. Update flight")
	fmt.Println("4. Delete flight")
	fmt.Println("5. Logout")
}

func getFlightInput(scanner *bufio.Scanner, flight *models.Flight) error {
	flight.FromWhere = service.GetInput(scanner, "Enter from where will arrive plane")
	flight.Destination = service.GetInput(scanner, "Enter destination")
	flightTimeStr := service.GetInput(scanner, "Enter flight time (format: YYYY-MM-DD HH:MM:SS)")

	flightTime, err := service.ParseTime(flightTimeStr, "2006-01-02 15:04:05")
	if err != nil {
		return err
	}
	flight.FlightTime = flightTime

	flightDuration, err := service.GetInt(scanner, "Enter flight duration")
	if err != nil {
		return err
	}
	flight.FlightDuration = flightDuration

	return nil
}

func clearConsole() {
	// Implement console clearing logic here, depending on your platform
	// This can be done using ANSI codes or system commands
}
