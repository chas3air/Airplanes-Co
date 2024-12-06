package customeradmininterface

import (
	"bufio"
	"fmt"
	"log"
	"os"

	customersfunctions "github.com/chas3air/Airplanes-Co/Client/CLI/internal/Functions/CustomersFunctions"
	"github.com/chas3air/Airplanes-Co/Client/CLI/internal/models"
	"github.com/chas3air/Airplanes-Co/Client/CLI/internal/service"
	"github.com/google/uuid"
)

func CustomersAdminInterface(user *models.Customer) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		service.ClearConsole()
		displayMenu()
		_ = scanner.Scan()
		choice := scanner.Text()

		switch choice {
		case "1":
			fmt.Println("Get all customers")

			localCustomers, err := customersfunctions.GetAllCustomers()
			if err != nil {
				fmt.Println("Customers weren't loaded:", err)
				bufio.NewReader(os.Stdin).ReadString('\n')
				break
			}

			customersfunctions.PrintCustomers(localCustomers)
			fmt.Println("Press Enter to continue...")
			bufio.NewReader(os.Stdin).ReadString('\n')

		case "2":
			fmt.Println("Show customer by id")
			id := service.GetInput(scanner, "Enter id")
			customer, err := customersfunctions.GetCustomerById(id)
			if err != nil {
				fmt.Println("Cannot get customer by id")
				bufio.NewReader(os.Stdin).ReadString('\n')
				break
			}

			customersfunctions.PrintCustomers([]models.Customer{customer})
			fmt.Println("Press Enter to continue...")
			bufio.NewReader(os.Stdin).ReadString('\n')

		case "3":
			fmt.Println("Adding customer")
			customer, err := customersfunctions.CreateCustomer()
			if err != nil {
				fmt.Println("Cannot create customer")
				break
			}

			log.Println("Created customer:", customer.String())

			customer, err = customersfunctions.InsertCustomer(customer)
			if err != nil {
				fmt.Println("Cannot insert customer")
				break
			}

			customersfunctions.PrintCustomers([]models.Customer{customer})
			fmt.Println("Press Enter to continue...")
			bufio.NewReader(os.Stdin).ReadString('\n')

		case "4":
			fmt.Println("Update customers")

			id := service.GetInput(scanner, "Enter id")

			customer, err := customersfunctions.CreateCustomer()
			if err != nil {
				fmt.Println("Error creating customer")
				break
			}
			parsedId, err := uuid.Parse(id)
			if err != nil {
				fmt.Println("Cannot parse string to uuid")
				bufio.NewReader(os.Stdin).ReadString('\n')
				break
			}
			customer.Id = parsedId

			_, err = customersfunctions.UpdateCustomer(customer)
			if err != nil {
				log.Println("Cannot update customer with id:", id)
				bufio.NewReader(os.Stdin).ReadString('\n')
				break
			}

			fmt.Println(customer)
			fmt.Println("Press Enter to continue...")
			bufio.NewReader(os.Stdin).ReadString('\n')

		case "5":
			fmt.Println("Deleting customer")

			id := service.GetInput(scanner, "Enter id")
			customer, err := customersfunctions.DeleteCustomer(id)
			if err != nil {
				log.Println("Cannot delete customer")
				bufio.NewReader(os.Stdin).ReadString('\n')
				break
			}

			fmt.Println(customer)
			log.Println("Press Enter to cintinue...")
			bufio.NewReader(os.Stdin).ReadString('\n')

		case "6":
			fmt.Println("Logout")

			if err := service.Logout(); err != nil {
				fmt.Println("Error loging out")
				bufio.NewReader(os.Stdin).ReadString('\n')
				break
			}

			*user = models.Customer{}
			return

		default:
			fmt.Println("Error number of item")
			bufio.NewReader(os.Stdin).ReadString('\n')
		}
	}
}

func displayMenu() {
	fmt.Println("Select an item")
	fmt.Println("1. Get all customers")
	fmt.Println("2. Get customer by id")
	fmt.Println("3. Add customer")
	fmt.Println("4. Update customer")
	fmt.Println("5. Delete customer")
	fmt.Println("6. Logout")
}
