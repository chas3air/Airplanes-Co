package customeradmininterface

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	customersfunctions "github.com/chas3air/Airplanes-Co/Client/CLI/internal/Functions/CustomersFunctions"
	"github.com/chas3air/Airplanes-Co/Client/CLI/internal/models"
)

func customersAdminInterface(user *models.Customer) {
	scanner := bufio.NewScanner(os.Stdin)
	var localCustomers []models.Customer
	var prepCustomerToInsert = make([]models.Customer, 0, 5)
	var prepCustomerToUpdate = make([]models.Customer, 0, 5)
	var prepIdCustomerToDelete = make([]string, 0, 5)
	var err error
	var mut sync.Mutex

	localCustomers, err = customersfunctions.GetAllCustomers()
	if err != nil {
		fmt.Println("customers weren't loaded", err)
		return
	}

	go func() {
		for {
			mut.Lock()
			processPreparations(&prepCustomerToInsert, &prepCustomerToUpdate, &prepIdCustomerToDelete)
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
			fmt.Println("Get all customers")

			for _, v := range localCustomers {
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

func processPreparations(prepCustomerToInsert, prepCustomerToUpdate *[]models.Customer, prepIdCustomerToDelete *[]string) {
	for _, v := range *prepCustomerToInsert {
		if _, err := customersfunctions.PostCustomer(v); err != nil {
			log.Println("Error posting flight:", err)
		}
	}

	for _, v := range *prepCustomerToUpdate {
		if _, err := customersfunctions.UpdateCustomer(v); err != nil {
			log.Println("Error updating flight:", err)
		}
	}

	for _, v := range *prepIdCustomerToDelete {
		if _, err := customersfunctions.DeleteCustomer(v); err != nil {
			log.Println("Error deleting flight:", err)
		}
	}

	*prepCustomerToInsert = make([]models.Customer, 0, 5)
	*prepCustomerToUpdate = make([]models.Customer, 0, 5)
	*prepIdCustomerToDelete = make([]string, 0, 10)
}

func displayMenu() {
	fmt.Println("Select an item")
	fmt.Println("1. Get all customers")
	fmt.Println("2. Add customer")
	fmt.Println("3. Update customer")
	fmt.Println("4. Delete customer")
	fmt.Println("5. Logout")
}

func clearConsole() {
	// Implement console clearing logic here, depending on your platform
	// This can be done using ANSI codes or system commands
}
