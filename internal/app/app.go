package app

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"test-api/internal/app/config"
	"test-api/internal/db"
	"test-api/internal/db/postgresql"
	"test-api/internal/http/handler"
	"test-api/internal/http/middleware"
	logs "test-api/internal/log"
	"test-api/internal/service"
	"time"

	"github.com/gofiber/fiber/v3"
)

type app struct {
	server    *fiber.App
	logger    *logs.Logger
	resources []io.Closer
	db        *db.DbManager
	config    *config.Config
}

func NewApp(config *config.Config) *app {
	return &app{
		config: config,
	}
}

func (a *app) Run() error {
	log.Println("INFO: Starting service...")

	var err error

	if a.db, err = db.Connect(a.config.Database); err != nil {
		return err
	}
	a.resources = append(a.resources, a.db)

	a.logger, err = logs.NewLogger(a.config.Logger)
	if err != nil {
		return err
	}

	a.resources = append(a.resources, a.logger)

	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		return errors.New("APP_PORT var is missing")
	}

	a.server = fiber.New()

	a.initialize(a.server)

	log.Println("INFO: Starting HTTP server...")

	err = a.server.Listen(":" + appPort)
	if err != nil {
		if err != http.ErrServerClosed {
			log.Printf("ERR: Error creating listener: %v\n", err)
			return err
		}
		return nil
	}

	return nil
}

func (a *app) Shutdown() error {
	log.Println("INFO: Server shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		// Release all resources
		var wg sync.WaitGroup
		for _, res := range a.resources {
			wg.Add(1)
			go func(res io.Closer) {
				if err := res.Close(); err != nil {
					log.Printf("ERR: Error closing resource: %s\n", err.Error())
				}
				wg.Done()
			}(res)
		}
		wg.Wait()
		cancel()
	}()

	if err := a.server.ShutdownWithContext(ctx); err != nil {
		return err
	}

	return nil
}

func (a *app) initialize(router *fiber.App) {
	router.Use(
		middleware.LoggerMiddleware(a.logger),
	)

	handler.NewArticleHandler(
		service.NewNewsService(
			postgresql.NewArticleRepository(a.db),
		),
	).Register(router)
}
