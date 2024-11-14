package guestinterface

import (
	"bufio"
	"fmt"
	"os"

	flightsfunctions "github.com/chas3air/Airplanes-Co/Client/CLI/internal/Functions/FlightsFunctions"
	"github.com/chas3air/Airplanes-Co/Client/CLI/internal/models"
	"github.com/chas3air/Airplanes-Co/Client/CLI/internal/service"
)

func GuestInterface(user *models.Customer) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		service.ClearConsole()
		fmt.Println("\nSelect an item")
		fmt.Println("1. Show all flights")
		fmt.Println("2. Logout")
		fmt.Println("3. Exit")

		_ = scanner.Scan()
		choice := scanner.Text()

		switch choice {
		case "1":
			flights, err := flightsfunctions.GetAllFlights()
			if err != nil {
				fmt.Println(err)
				break
			}

			flightsfunctions.PrintFlights(flights)

			fmt.Println("Press Enter to continue...")
			bufio.NewReader(os.Stdin).ReadString('\n')

		case "2":
			// Add login functionality here
			fmt.Println("Login functionality not implemented yet.")

		case "3":
			fmt.Println("Exiting the program.")
			return // Exit the loop and terminate the function

		default:
			fmt.Println("Error: Invalid option. Please try again.")
		}
	}
}
