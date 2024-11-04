package app

import (
	"bufio"
	"fmt"
	"os"

	flightsadmininterface "github.com/chas3air/Airplanes-Co/Client/CLI/internal/Interfaces/FlightsAdminInterface"
	"github.com/chas3air/Airplanes-Co/Client/CLI/internal/config"
	"github.com/chas3air/Airplanes-Co/Client/CLI/internal/models"
	"github.com/chas3air/Airplanes-Co/Client/CLI/internal/service"
)

func Run() {
	var current_customer models.Customer

	for {
		switch current_customer.Role {
		case config.FlightsAdmin:
			FlightsAdminInterface(&current_customer)
		case config.CustomersAdmin:
			CustomersAdminInterface(&current_customer)
		case config.GeneralAdmin:
			GeneralAdminInterface(&current_customer)
		case config.User:
			CustomersInterface(&current_customer)
		case config.Guest:
			GuestInterface(&current_customer)
		}
	}
}

func FlightsAdminInterface(user *models.Customer) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		// тут очистка консоли
		fmt.Println("Select an item")
		fmt.Println("1. Get all flights")
		fmt.Println("2. Add flight")
		fmt.Println("3. Update flight")
		fmt.Println("4. Delete flight")
		fmt.Println("5. Logout")
		_ = scanner.Scan()
		choise := scanner.Text()

		switch choise {
		case "1":
			flights, err := flightsadmininterface.GetAllFlights()
			if err != nil {
				fmt.Println(err.Error())
				break
			}

			for _, v := range flights {
				fmt.Println(v)
			}

		case "2":
			flight, err := flightsadmininterface.AddFlight()
			if err != nil {
				fmt.Println(err.Error())
				break
			}

			fmt.Println("Added flight:", flight.String())

		case "3":
			flights, err := flightsadmininterface.GetAllFlights()
			if err != nil {
				fmt.Println(err.Error())
				break
			}

			for _, v := range flights {
				fmt.Println(v)
			}

			id, err := service.GetInt(scanner, "Enter id")
			if err != nil {
				fmt.Println(err.Error())
				break
			}
			if id < 0 || id > len(flights) {
				fmt.Println("Undefined element")
				break
			}

			flight, err := flightsadmininterface.UpdateFlight(flights[id-1].Id.String())
			if err != nil {
				fmt.Println(err)
				break
			}

			fmt.Println("Updated flight:", flight.String())

		case "4":
			flights, err := flightsadmininterface.GetAllFlights()
			if err != nil {
				fmt.Println(err.Error())
				break
			}

			for _, v := range flights {
				fmt.Println(v)
			}

			id, err := service.GetInt(scanner, "Enter id")
			if err != nil {
				fmt.Println(err.Error())
				break
			}
			if id < 0 || id > len(flights) {
				fmt.Println("Undefined element")
				break
			}

		case "5":
			if err := service.Logout(); err != nil {
				break
			}
			*user = models.Customer{}

		default:
			fmt.Println("Error number of item")
		}
	}
}
func CustomersAdminInterface(user *models.Customer) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("Select an item")
		fmt.Println("1. Get all customers")
		fmt.Println("2. Add customer")
		fmt.Println("3. Update customer")
		fmt.Println("4. Delete customer")
		fmt.Println("5. Logout")
		_ = scanner.Scan()
		choise := scanner.Text()

		switch choise {
		case "1":
		case "2":
		case "3":
		case "4":
		case "5":
		default:
			fmt.Println("Error number of item")
		}
	}
}

func GeneralAdminInterface(user *models.Customer) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("Select an item")
		fmt.Println("1. Get all flights")
		fmt.Println("2. Add flight")
		fmt.Println("3. Update flight")
		fmt.Println("4. Delete flight")
		fmt.Println("5. Get all customers")
		fmt.Println("6. Add customer")
		fmt.Println("7. Update customer")
		fmt.Println("8. Delete customer")
		fmt.Println("9. Logout")
		_ = scanner.Scan()
		choise := scanner.Text()

		switch choise {
		case "1":
		case "2":
		case "3":
		case "4":
		case "5":
		case "6":
		case "7":
		case "8":
		case "9":
		default:
			fmt.Println("Error number of item")
		}
	}
}

func CustomersInterface(user *models.Customer) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("Select an item")
		fmt.Println("1. Show all flights")
		fmt.Println("2. Buy ticket")
		fmt.Println("3. Show cart")
		fmt.Println("4. Pay for cart")
		fmt.Println("5. Manage tickets")

		fmt.Println("6. Logout")
		_ = scanner.Scan()
		choise := scanner.Text()

		switch choise {
		case "1":
		case "2":
		case "3":
		case "4":
		case "5":
		case "6":
		default:
			fmt.Println("Error number of item")
		}
	}
}

func GuestInterface(user *models.Customer) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("Select an item")
		fmt.Println("1. Show all flights")
		fmt.Println("2. Login")
		fmt.Println("3. Logout")
		_ = scanner.Scan()
		choise := scanner.Text()

		switch choise {
		case "1":
		case "2":
		case "3":
		default:
			fmt.Println("Error number of item")
		}
	}
}