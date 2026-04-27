package main

import (
	"CameraAscciEngine/camera"
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"
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
	ctx := context.Background()
	err := camera.RunCam(ctx)
	if err != nil {
		fmt.Print("\033[H\033[2J") // Clear screen and move to top-left
		return
	}

	// Create a done channel to signal when the shutdown is complete
	done := make(chan bool, 1)

	// Run graceful shutdown in a separate goroutine
	go gracefulShutdown(done)

	// Wait for the graceful shutdown to complete
	// <-done
	log.Println("Graceful shutdown complete.")
}
