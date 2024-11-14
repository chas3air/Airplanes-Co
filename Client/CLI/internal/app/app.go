package app

import (
	"bufio"
	"fmt"
	"os"
	"time"

	customersfunctions "github.com/chas3air/Airplanes-Co/Client/CLI/internal/Functions/CustomersFunctions"
	fai "github.com/chas3air/Airplanes-Co/Client/CLI/internal/Interfaces/FligthtAdminInterface"
	guestinterface "github.com/chas3air/Airplanes-Co/Client/CLI/internal/Interfaces/GuestInterface"
	ui "github.com/chas3air/Airplanes-Co/Client/CLI/internal/Interfaces/UserInterface"
	"github.com/chas3air/Airplanes-Co/Client/CLI/internal/config"
	"github.com/chas3air/Airplanes-Co/Client/CLI/internal/models"
)

func Run() {
	var current_customer models.Customer
	current_customer.Role = config.FlightsAdmin

	for {
		switch current_customer.Role {
		case config.FlightsAdmin:
			fai.FlightsAdminInterface(&current_customer)
		case config.CustomersAdmin:
			customersAdminInterface(&current_customer)
		case config.GeneralAdmin:
			generalAdminInterface(&current_customer)
		case config.User:
			ui.UserInterface(&current_customer)
		default:
			guestinterface.GuestInterface(&current_customer)
		}

		fmt.Println("Do you want to exit?(Y)")
		var action string
		fmt.Scanln(&action)

		if action == "Y" {
			fmt.Println("Exiting program.")
			os.Exit(0)
		} else {
			fmt.Println("You are in yet")
		}
	}
}

func customersAdminInterface(user *models.Customer) {
	scanner := bufio.NewScanner(os.Stdin)
	// var localCustomers []models.Customer
	// var prepCustomerToInsert = make([]models.Customer, 0, 5)
	// var prepCustomerToUpdate = make([]models.Customer, 0, 5)
	// var prepIdCustomerToDelete = make([]string, 0, 5)
	var err error
	//var mut sync.Mutex

	localCustomers, err := customersfunctions.GetAllCustomers()
	if err != nil {
		fmt.Println("customers weren't loaded", err)
		return
	}
	_ = localCustomers

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
			fmt.Println("Get all customers")

			customers, err := customersfunctions.GetAllCustomers()
			if err != nil {
				fmt.Println(err)
				break
			}
			for _, v := range customers {
				fmt.Println(v.String())
			}
			bufio.NewReader(os.Stdin).ReadString('\n')

		case "2":
			fmt.Println("Adding customer")

			time.Sleep(200 * time.Millisecond)

		case "3":
			fmt.Println("Update customers")
			time.Sleep(200 * time.Millisecond)

		case "4":
			fmt.Println("Deleting customer")
			time.Sleep(200 * time.Millisecond)

		case "5":
			fmt.Println("Logout")
			time.Sleep(200 * time.Millisecond)

		default:
			fmt.Println("Error number of item")
			time.Sleep(200 * time.Millisecond)

		}
	}
}

func generalAdminInterface(user *models.Customer) {
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
