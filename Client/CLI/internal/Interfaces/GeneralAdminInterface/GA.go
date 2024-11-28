package ga_interface

import (
	"bufio"
	"fmt"
	"os"
	"time"

	flightsfunctions "github.com/chas3air/Airplanes-Co/Client/CLI/internal/Functions/FlightsFunctions"
	"github.com/chas3air/Airplanes-Co/Client/CLI/internal/models"
	"github.com/chas3air/Airplanes-Co/Client/CLI/internal/service"
	"github.com/google/uuid"
)

func GeneralAdminInterface(user *models.Customer) {
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
				bufio.NewReader(os.Stdin).ReadString('\n')
				break
			}

			flightsfunctions.PrintFlights(localFlights)
			fmt.Println("Press Enter to continue...")
			bufio.NewReader(os.Stdin).ReadString('\n')

		case "2":
			fmt.Println("Show flight by id")
			id := service.GetInput(scanner, "Enter id")
			flight, err := flightsfunctions.GetFlightById(id)
			if err != nil {
				fmt.Println("Cannot get flight by id")
				bufio.NewReader(os.Stdin).ReadString('\n')
				break
			}

			flightsfunctions.PrintFlights([]models.Flight{flight})
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

			flightsfunctions.PrintFlights([]models.Flight{flight})
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
				bufio.NewReader(os.Stdin).ReadString('\n')
				break
			}
			flight.Id = parsedId

			_, err = flightsfunctions.UpdateFlight(flight)
			if err != nil {
				fmt.Println("Cannot update flight with id:", id)
				bufio.NewReader(os.Stdin).ReadString('\n')
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

		case "7":
		case "8":
		case "9":
		case "10":
		case "11":
		case "12":
		case "13":
		case "14":
		case "15":
		case "16":
			fmt.Println("Logout")
			if err := service.Logout(); err != nil {
				fmt.Println("Error logging out:", err)
				bufio.NewReader(os.Stdin).ReadString('\n')
				break
			}
			*user = models.Customer{}

			fmt.Println("Logout goes successfully")
			bufio.NewReader(os.Stdin).ReadString('\n')
			return

		default:
			fmt.Println("Error: invalid menu option")
			bufio.NewReader(os.Stdin).ReadString('\n')
		}
		time.Sleep(200 * time.Millisecond)
	}
}

func displayMenu() {
	fmt.Println("Select an item")
	fmt.Println("1. Get all flights")
	fmt.Println("2. Get flight by id")
	fmt.Println("3. Add flight")
	fmt.Println("4. Update flight")
	fmt.Println("5. Delete flight")

	fmt.Println("6. Get all customers")
	fmt.Println("7. Get customersby id")
	fmt.Println("8. Add customer")
	fmt.Println("9. Update customer")
	fmt.Println("10. Delete customer")

	fmt.Println("11. Get all tickets")
	fmt.Println("12. Get tickets by id")
	fmt.Println("13. Add tickets")
	fmt.Println("14. Update tickets")
	fmt.Println("15. Delete tickets")

	fmt.Println("16. Logout")
}
