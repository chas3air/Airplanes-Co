package guestinterface

import (
	"bufio"
	"fmt"
	"os"

	customersfunctions "github.com/chas3air/Airplanes-Co/Client/CLI/internal/Functions/CustomersFunctions"
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
		fmt.Println("2. Show flight by id")
		fmt.Println("3. Sign up")
		fmt.Println("4. Sign in")
		fmt.Println("5. Exit")

		if !scanner.Scan() {
			fmt.Println("Error reading input. Exiting...")
			return
		}
		choice := scanner.Text()

		switch choice {
		case "1":
			fmt.Println("Show all flights")
			flights, err := flightsfunctions.GetAllFlights()
			if err != nil {
				fmt.Println("Error:", err)
				waitForEnter()
				continue
			}

			flightsfunctions.PrintFlights(flights, "ID")
			waitForEnter()

		case "2":
			fmt.Println("Show flight by id")
			id := service.GetInput(scanner, "Enter id")
			flight, err := flightsfunctions.GetFlightById(id)
			if err != nil {
				fmt.Println("Cannot get flight by id")
				waitForEnter()
				break
			}

			flightsfunctions.PrintFlights([]models.Flight{flight})
			waitForEnter()

		case "3":
			fmt.Println("Sign up")
			login := service.GetInput(scanner, "Enter login")
			password := service.GetInput(scanner, "Enter password")
			surname := service.GetInput(scanner, "Enter surname")
			name := service.GetInput(scanner, "Enter name")

			customer, err := customersfunctions.SignUpCustomer(
				models.Customer{
					Login:    login,
					Password: password,
					Surname:  surname,
					Name:     name,
				},
			)

			if err != nil {
				fmt.Println("Error:", err)
				waitForEnter()
				break
			}

			*user = customer
			fmt.Println("Sign up successfully")
			waitForEnter()
			return

		case "4":
			login := service.GetInput(scanner, "Enter login")
			password := service.GetInput(scanner, "Enter password")

			customer, err := customersfunctions.SignInCustomer(login, password)
			if err != nil {
				fmt.Println("Error:", err)
				waitForEnter()
				break
			}

			*user = customer
			waitForEnter()
			return

		case "5":
			fmt.Println("Exiting the program.")
			return

		default:
			fmt.Println("Error: Invalid option. Please try again.")
			waitForEnter()
		}
	}
}

func waitForEnter() {
	fmt.Println("Press Enter to continue...")
	bufio.NewReader(os.Stdin).ReadString('\n')
}
