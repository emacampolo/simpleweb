package main

import (
	"bootcamp-web/internal/application"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, syscall.SIGINT, syscall.SIGTERM)

	app := application.New()

	go func() {
		if err := app.Start(); err != nil {
			// If the error is ErrServerClosed,
			// then it was a graceful shutdown request.
			if !errors.Is(err, http.ErrServerClosed) {
				log.Fatal(err)
			}
		}
	}()

	// Block until a signal is received.
	<-shutdownChan
	log.Printf("Shutting down...\n")
	if err := app.Stop(); err != nil {
		log.Fatal(err)
	}
}
