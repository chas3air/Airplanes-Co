package app

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
)

func Run() {
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	_ = cancel
	r := mux.NewRouter()

	// тут вызов функций

	go func() {
		if err := http.ListenAndServe(":84", r); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Received shutdown signal, shitting down...")
}
