package app

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	customersfunctions "github.com/chas3air/Airplanes-Co/Client/CLI/internal/Functions/CustomersFunctions"
	flightsfunctions "github.com/chas3air/Airplanes-Co/Client/CLI/internal/Functions/FlightsFunctions"
	"github.com/chas3air/Airplanes-Co/Client/CLI/internal/config"
	"github.com/chas3air/Airplanes-Co/Client/CLI/internal/models"
	"github.com/chas3air/Airplanes-Co/Client/CLI/internal/service"
)

func Run() {
	var current_customer models.Customer

	for {
		switch current_customer.Role {
		case config.FlightsAdmin:
			flightsAdminInterface(&current_customer)
		case config.CustomersAdmin:
			customersAdminInterface(&current_customer)
		case config.GeneralAdmin:
			generalAdminInterface(&current_customer)
		case config.User:
			customersInterface(&current_customer)
		case config.Guest:
			guestInterface(&current_customer)
		}
	}
}

func flightsAdminInterface(user *models.Customer) {
	scanner := bufio.NewScanner(os.Stdin)
	var localFlights []models.Flight
	var prepFlightToInsert = make([]models.Flight, 0, 5)
	var prepFlightToUpdate = make([]models.Flight, 0, 5)
	var prepIdFlightToDelete = make([]string, 0, 10)
	var err error
	var mut sync.Mutex

	localFlights, err = flightsfunctions.GetAllFlights()
	if err != nil {
		fmt.Println("Flights wasnt been loaded")
	}

	go func() {
		for {
			mut.Lock()
			for _, v := range prepFlightToInsert {
				_, err := flightsfunctions.PostFlight(v)
				if err != nil {
					log.Println(err)
				}
			}

			for _, v := range prepFlightToUpdate {
				_, err := flightsfunctions.UpdateFlight(v)
				if err != nil {
					log.Println()
				}
			}

			for _, v := range prepIdFlightToDelete {
				_, err := flightsfunctions.DeleteFlight(v)
				if err != nil {
					log.Println(err)
				}
			}

			prepFlightToInsert = make([]models.Flight, 0, 10)
			prepFlightToUpdate = make([]models.Flight, 0, 10)
			prepIdFlightToDelete = make([]string, 0, 10)
			mut.Unlock()
			time.Sleep(3 * time.Second)
		}
	}()

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
			fmt.Println("Show all flights")
			for _, v := range localFlights {
				fmt.Println(v)
			}
			bufio.NewReader(os.Stdin).ReadString('\n')

		case "2":
			fmt.Println("Adding flights")
			flight := models.Flight{}
			localFlights = append(localFlights, flight)
			prepFlightToInsert = append(prepFlightToInsert, flight)

			time.Sleep(200 * time.Millisecond)

		case "3":
			fmt.Println("Update flight")
			for _, v := range localFlights {
				fmt.Println(v)
			}

			id, err := service.GetInt(scanner, "Enter id")
			if err != nil {
				fmt.Println(err.Error())
				break
			}
			if id < 0 || id > len(localFlights) {
				fmt.Println("Undefined element")
				break
			}

			flight, err := flightsfunctions.CreateFlight()
			if err != nil {
				log.Println(err)
				break
			}

			prepFlightToUpdate = append(prepFlightToUpdate, flight)

			for i, v := range localFlights {
				if v.Id == localFlights[id+1].Id {
					localFlights[i].FromWhere = flight.FromWhere
					localFlights[i].Destination = flight.Destination
					localFlights[i].FlightTime = flight.FlightTime
					localFlights[i].FlightDuration = flight.FlightDuration

					fmt.Println("Updated flight:", v.String())
					break
				}
			}
			time.Sleep(200 * time.Millisecond)

		case "4":
			fmt.Println("Deleting flight")

			for _, v := range localFlights {
				fmt.Println(v)
			}

			id, err := service.GetInt(scanner, "Enter id")
			if err != nil {
				fmt.Println(err.Error())
				break
			}
			if id < 0 || id > len(localFlights) {
				fmt.Println("Undefined element")
				break
			}

			prepIdFlightToDelete = append(prepIdFlightToDelete, localFlights[id+1].Id.String())
			time.Sleep(200 * time.Millisecond)

		case "5":
			fmt.Println("Logout")

			if err := service.Logout(); err != nil {
				break
			}
			*user = models.Customer{}
			time.Sleep(200 * time.Millisecond)

		default:
			fmt.Println("Error number of item")
			time.Sleep(200 * time.Millisecond)
		}
	}
}

func customersAdminInterface(user *models.Customer) {
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

func customersInterface(user *models.Customer) {
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

func guestInterface(user *models.Customer) {
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
			flights, err := flightsfunctions.GetAllFlights()
			if err != nil {
				fmt.Println(err)
				break
			}

			for _, v := range flights {
				fmt.Println(v.String())
			}

			bufio.NewReader(os.Stdin).ReadString('\n')
		case "2":
		case "3":
		default:
			fmt.Println("Error number of item")
		}
	}
}
