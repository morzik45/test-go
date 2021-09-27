package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	exam "github.com/morzik45/test-go"
	"github.com/morzik45/test-go/config"
	"github.com/morzik45/test-go/logger"
	"github.com/morzik45/test-go/pkg/handler"
	"github.com/morzik45/test-go/pkg/repository"
	"github.com/morzik45/test-go/pkg/service"
)

func main() {
	configPath := flag.String("configFile", "./config/config.json", "Path to config file")
	flag.Parse()
	cfg, err := config.InitConfig(*configPath)
	if err != nil {
		logger.ERROR.Fatalf("error initializing configs: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(cfg.Db)
	if err != nil {
		logger.ERROR.Fatalf("failed to initialize db: %s", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(exam.HttpServer)
	go func() {
		if err := srv.Run(cfg.Http.Port, handlers.InitRoutes()); err != nil {
			if err == http.ErrServerClosed {
				logger.INFO.Println("Http server been closed.")
			} else {
				logger.ERROR.Fatalf("error occured while running http server: %s", err.Error())
			}
		}
	}()

	logger.INFO.Println("App started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logger.INFO.Println("App shutting down")

	if err := srv.Shutdown(context.Background()); err != nil {
		logger.ERROR.Printf("error occured on server shutting down: %s", err.Error())
	}

	if err := db.Close(); err != nil {
		logger.ERROR.Printf("error occured on db connection close: %s", err.Error())
	}
}
