package ui

import (
	"bufio"
	"fmt"
	"os"
	"time"

	flightsfunctions "github.com/chas3air/Airplanes-Co/Client/CLI/internal/Functions/FlightsFunctions"
	"github.com/chas3air/Airplanes-Co/Client/CLI/internal/models"
	"github.com/chas3air/Airplanes-Co/Client/CLI/internal/service"
)

func UserInterface(user *models.Customer) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		service.ClearConsole()
		displayMenu()
		_ = scanner.Scan()
		choice := scanner.Text()

		switch choice {
		case "1":
			fmt.Println("Show all flights")
			localFlights, err := flightsfunctions.GetAllFlights()
			if err != nil {
				fmt.Println("Flights weren't loaded:", err)
				fmt.Println("Press Enter to continue...")
				bufio.NewReader(os.Stdin).ReadString('\n')
				break
			}

			flightsfunctions.PrintFlights(localFlights)
			fmt.Println("Press Enter to continue...")
			bufio.NewReader(os.Stdin).ReadString('\n')

		case "2":
			fmt.Println("Buy ticket")

			from := service.GetInput(scanner, "Enter from where u need flight")
			dest := service.GetInput(scanner, "Enter destination u need")
			flightTimeStr := service.GetInput(scanner, "Enter flight time (format: YYYY-MM-DD HH:MM:SS)")

			flightTime, err := service.ParseTime(flightTimeStr, "2006-01-02 15:04:05")
			if err != nil {
				fmt.Println("Error:", err)
				fmt.Println("Press Enter to continue...")
				bufio.NewReader(os.Stdin).ReadString('\n')
				break
			}

			flights, err := flightsfunctions.GetAllFlights()
			if err != nil {
				fmt.Println("Flights weren't loaded:", err)
				fmt.Println("Press Enter to continue...")
				bufio.NewReader(os.Stdin).ReadString('\n')
				break
			}

			var current_flight models.Flight
			for i, v := range flights {
				if v.FromWhere == from && v.Destination == dest && v.FlightTime == flightTime {
					current_flight = flights[i]
					break
				}
			}

			ticket, err := CreateTicket(current_flight)
			if err != nil {
				fmt.Println("Cannot create ticket to flight")
				fmt.Println("Press Enter to continue...")
				bufio.NewReader(os.Stdin).ReadString('\n')
				break
			}

			ticket.Owner = *user

			fmt.Println("Ticket:", ticket)

			fmt.Println("Press Enter to continur...")
			bufio.NewReader(os.Stdin).ReadString('\n')

		case "3":
			fmt.Println("Show cart")

			fmt.Println("Press Enter to continue...")
			bufio.NewReader(os.Stdin).ReadString('\n')

		case "4":
			fmt.Println("Pay for cart")

			fmt.Println("Press Enter to continue...")
			bufio.NewReader(os.Stdin).ReadString('\n')

		case "5":
			fmt.Println("Manage tickets")

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

func displayMenu() {
	fmt.Println("Select an item")
	fmt.Println("1. Show all flights")
	fmt.Println("2. Buy ticket")
	fmt.Println("3. Show cart")
	fmt.Println("4. Pay for cart")
	fmt.Println("5. Manage tickets")
	fmt.Println("6. Logout")
}
