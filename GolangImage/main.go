
package main

import (
	"context"
	"log"
	"os/signal"
	"CameraAscciEngine/camera"
	"syscall"
	// "quick/database"
)

func gracefulShutdown(done chan bool) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")
	stop() // Allow Ctrl+C to force shutdown

	done <- true
}

func main() {
	err := camera.RunCam()
	if err != nil {
		panic(err)
	}

	// Create a done channel to signal when the shutdown is complete
	done := make(chan bool, 1)

	// Run graceful shutdown in a separate goroutine
	go gracefulShutdown(done)

	// Wait for the graceful shutdown to complete
	// <-done
	log.Println("Graceful shutdown complete.")
}
