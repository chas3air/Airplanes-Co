package fligthtadmininterface

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	flightsfunctions "github.com/chas3air/Airplanes-Co/Client/CLI/internal/Functions/FlightsFunctions"
	"github.com/chas3air/Airplanes-Co/Client/CLI/internal/models"
	"github.com/chas3air/Airplanes-Co/Client/CLI/internal/service"
	"github.com/google/uuid"
)

func FlightsAdminInterface(user *models.Customer) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		clearConsole()
		displayMenu()
		_ = scanner.Scan()
		choice := scanner.Text()

		switch choice {
		case "1":
			fmt.Println("Show all flights")
			localFlights, err := flightsfunctions.GetAllFlights()
			if err != nil {
				fmt.Println("Flights weren't loaded:", err)
				time.Sleep(200 * time.Millisecond)
				break
			}
			displayFlights(localFlights)
			fmt.Println("Press Enter to continue...")
			bufio.NewReader(os.Stdin).ReadString('\n')

		case "2":
			fmt.Println("Show flight")
			id := service.GetInput(scanner, "Enter id")
			flight, err := flightsfunctions.GetFlightById(id)
			if err != nil {
				fmt.Println("Cannot get flight by id")
				time.Sleep(200 * time.Millisecond)
				break
			}

			fmt.Println(flight)
			fmt.Println("Press Enter to continue...")
			bufio.NewReader(os.Stdin).ReadString('\n')

		case "3":
			fmt.Println("Adding flights")
			flight, err := flightsfunctions.CreateFlight()
			if err != nil {
				fmt.Println("Cannot create flight")
				break
			}

			flight, err = flightsfunctions.InsertFlight(flight)
			if err != nil {
				fmt.Println("Cannot insert flight")
			}

			fmt.Println(flight)
			fmt.Println("Press Enter to continue...")
			bufio.NewReader(os.Stdin).ReadString('\n')

		case "4":
			fmt.Println("Update flight")

			id := service.GetInput(scanner, "Enter id")

			flight, err := flightsfunctions.CreateFlight()
			if err != nil {
				fmt.Println("Error creating flight:", err)
				break
			}
			parsedId, err := uuid.Parse(id)
			if err != nil {
				fmt.Println("Cannot parse string to uuid")
				time.Sleep(200 * time.Millisecond)
				break
			}
			flight.Id = parsedId

			_, err = flightsfunctions.UpdateFlight(flight)
			if err != nil {
				fmt.Println("Cannot update flight with id:", id)
				time.Sleep(200 * time.Millisecond)
				break
			}

			fmt.Println(flight)
			fmt.Println("Press Enter to continue...")
			bufio.NewReader(os.Stdin).ReadString('\n')

		case "5":
			fmt.Println("Deleting flight")

			id := service.GetInput(scanner, "Enter id")
			flight, err := flightsfunctions.DeleteFlight(id)
			if err != nil {
				fmt.Println("Cannot detere flight")
			}

			fmt.Println(flight)
			fmt.Println("Press Enter to continue...")
			bufio.NewReader(os.Stdin).ReadString('\n')

		case "6":
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

// TODO: под снос
func processPreparations(prepFlightToInsert, prepFlightToUpdate *[]models.Flight, prepIdFlightToDelete *[]string) {
	for _, v := range *prepFlightToInsert {
		if _, err := flightsfunctions.InsertFlight(v); err != nil {
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
}

func displayMenu() {
	fmt.Println("Select an item")
	fmt.Println("1. Get all flights")
	fmt.Println("2. Get a flight")
	fmt.Println("3. Add flight")
	fmt.Println("4. Update flight")
	fmt.Println("5. Delete flight")
	fmt.Println("6. Logout")
}

func clearConsole() {
	// Implement console clearing logic here, depending on your platform
	// This can be done using ANSI codes or system commands
}
