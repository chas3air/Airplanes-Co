package app

import (
	"fmt"
	"os"

	customeradmininterface "github.com/chas3air/Airplanes-Co/Client/CLI/internal/Interfaces/CustomerAdminInterface"
	fai "github.com/chas3air/Airplanes-Co/Client/CLI/internal/Interfaces/FligthtAdminInterface"
	ga_interface "github.com/chas3air/Airplanes-Co/Client/CLI/internal/Interfaces/GeneralAdminInterface"
	guestinterface "github.com/chas3air/Airplanes-Co/Client/CLI/internal/Interfaces/GuestInterface"
	ui "github.com/chas3air/Airplanes-Co/Client/CLI/internal/Interfaces/UserInterface"
	"github.com/chas3air/Airplanes-Co/Client/CLI/internal/config"
	"github.com/chas3air/Airplanes-Co/Client/CLI/internal/models"
)

func Run() {
	var current_customer models.Customer

	for {
		switch current_customer.Role {
		case config.FlightsAdmin:
			fai.FlightsAdminInterface(&current_customer)
		case config.CustomersAdmin:
			customeradmininterface.CustomersAdminInterface(&current_customer)
		case config.GeneralAdmin:
			ga_interface.GeneralAdminInterface(&current_customer)
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
