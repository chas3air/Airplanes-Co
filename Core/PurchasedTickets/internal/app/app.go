package app

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/chas3air/Airplanes-Co/Core/PurchasedTickets/internal/routes"
	"github.com/gorilla/mux"
)

func Run() {
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	_ = cancel
	r := mux.NewRouter()

	// example: /purchased-tickets?ownerId=ownerId
	r.HandleFunc("/purchased-tickets", routes.GetPurchasedTicketsHandler).Methods(http.MethodGet)

	go func() {
		if err := http.ListenAndServe(":84", r); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Received shutdown signal, shitting down...")
}
