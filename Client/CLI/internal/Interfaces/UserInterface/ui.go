package ui

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"

	flightsfunctions "github.com/chas3air/Airplanes-Co/Client/CLI/internal/Functions/FlightsFunctions"
	ticketsfunctions "github.com/chas3air/Airplanes-Co/Client/CLI/internal/Functions/TicketsFunctions"
	"github.com/chas3air/Airplanes-Co/Client/CLI/internal/config"
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

			flightsfunctions.PrintFlights(localFlights, "ID")
			fmt.Println("Press Enter to continue...")
			bufio.NewReader(os.Stdin).ReadString('\n')

		case "2":
			fmt.Println("Buy ticket")

			from := service.GetInput(scanner, "Enter from where u need flight")
			dest := service.GetInput(scanner, "Enter destination u need")
			var numOfFlight int = 0
			var costNum int = 0

			flights, err := flightsfunctions.GetFlightsWithFromAndDest(from, dest)
			if err != nil {
				fmt.Println("Flights weren't loaded:", err)
				fmt.Println("Press Enter to continue...")
				bufio.NewReader(os.Stdin).ReadString('\n')
				break
			}

			flightsfunctions.PrintFlights(flights, "ID")

			fmt.Println("Enter number of flights you want to order")
			scanner.Scan()

			if numOfFlight, err = strconv.Atoi(scanner.Text()); err != nil {
				fmt.Println("Undefined num of flight:", err)
				fmt.Println("Press Enter to continue...")
				bufio.NewReader(os.Stdin).ReadString('\n')
				break
			}

			if numOfFlight <= 0 || numOfFlight > len(flights) {
				fmt.Println("You entered:", numOfFlight, "but expected a number between 1 and:", len(flights))
				fmt.Println("Press Enter to continue...")
				bufio.NewReader(os.Stdin).ReadString('\n')
				break
			}

			selectedFlight := flights[numOfFlight-1]
			fmt.Println("You selected flight:", selectedFlight)

			fmt.Print("Enter number of cost: ")
			fmt.Scan(&costNum)

			if costNum <= 0 || costNum > len(selectedFlight.FlightSeatsCosts) {
				fmt.Println("You entered:", costNum, "but expected a number between 1 and:", len(selectedFlight.FlightSeatsCosts))
				fmt.Println("Press Enter to continue...")
				bufio.NewReader(os.Stdin).ReadString('\n')
				break
			}

			ticket := CreateTicket(user.Id, selectedFlight, selectedFlight.FlightSeatsCosts[costNum], config.NamesSeats[costNum])
			_, err = ticketsfunctions.SendTicketToTheCart(ticket)
			if err != nil {
				fmt.Println("Error:", err)
				fmt.Println("Press Enter to continue...")
				bufio.NewReader(os.Stdin).ReadString('\n')
				break
			}

			fmt.Println("Press Enter to continur...")
			bufio.NewReader(os.Stdin).ReadString('\n')

		case "3":
			fmt.Println("Show cart")
			ticketsFromCart, err := ticketsfunctions.GetTicketFromCart(user.Id)
			if err != nil {
				fmt.Println("Error:", err)

				fmt.Println("Press Enter to continue...")
				bufio.NewReader(os.Stdin).ReadString('\n')
			}

			fmt.Println("Tickets in cart")
			if len(ticketsFromCart) > 0 {
				ticketsfunctions.PrintTickets(ticketsFromCart, "ID", "Owner")
			}

			fmt.Println("Press Enter to continue...")
			bufio.NewReader(os.Stdin).ReadString('\n')

		case "4":
			fmt.Println("Pay for cart")

			card_number := service.GetInput(scanner, "Enter card number")
			var sum_cost int
			tickets, err := ticketsfunctions.GetTicketFromCart(user.Id)
			if err != nil {
				fmt.Println("Error:", err)
			}
			for _, v := range tickets {
				sum_cost += v.TicketCost
			}

			err = ticketsfunctions.PayForTickets(user.Id, card_number, sum_cost)
			if err != nil {
				fmt.Println("Error:", err)
				fmt.Println("Press Enter to continue...")
				bufio.NewReader(os.Stdin).ReadString('\n')
				break
			}

			fmt.Println("Tickets payed successfully...")

			fmt.Println("Press Enter to continue...")
			bufio.NewReader(os.Stdin).ReadString('\n')

		case "5":
			fmt.Println("Manage tickets")

			purchasedTickets, err := ticketsfunctions.GetPurchasedTickets(*user)
			if err != nil {
				fmt.Println("Error:", err)
			}

			fmt.Println("Yout purcahsed tickets")
			ticketsfunctions.PrintTickets(purchasedTickets, "ID", "Owner")

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
