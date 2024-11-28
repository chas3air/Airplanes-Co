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
				time.Sleep(200 * time.Millisecond)
				break
			}

			flightsfunctions.PrintFlights(localFlights)
			fmt.Println("Press Enter to continue...")
			bufio.NewReader(os.Stdin).ReadString('\n')

		case "2":
		case "3":
		case "4":
		case "5":
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
