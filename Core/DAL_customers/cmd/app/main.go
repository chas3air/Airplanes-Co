package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/chas3air/Airplanes-Co/Core/DAL_customers/internal/app"
)

func main() {
	go app.Run()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	<-signalChan
	log.Println("Exit")
}
