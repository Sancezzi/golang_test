package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"test-api/internal/app"
	"test-api/internal/app/config"

	"github.com/joho/godotenv"
)

func init() {
	basePath, err := os.Getwd()
	if err != nil {
		basePath = "."
	}
	envFile := basePath + "/.env"

	if err := godotenv.Overload(envFile); err != nil {
		log.Fatal("No .env file found")
	}
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	rest := app.NewApp(config.GetInstance())

	go func() {
		if err := rest.Run(); err != nil {
			log.Fatalf("FATAL: Cannot run app: %v\n", err)
		}
	}()

	<-ctx.Done()

	if err := rest.Shutdown(); err != nil {
		log.Fatalf("FATAL: Server close: %v\n", err)
	}

	log.Println("INFO: Server closed succesfully")
}
