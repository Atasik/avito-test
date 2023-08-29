package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"segmenter/internal/config"
	"segmenter/internal/handler"
	"segmenter/internal/repository"
	"segmenter/internal/server"
	"segmenter/internal/service"
	"segmenter/pkg/postgres"
	"syscall"
	"time"
)

const (
	interval = 30 * time.Second
)

// @title Avito Backend Trainee Assignment
// @version 1.0
// @description тех. задание с отбора на стажировку в Avito

// @host localhost:8080
// @BasePath /
func Run(configPath string) {
	cfg, err := config.InitConfig(configPath)
	if err != nil {
		log.Fatal("Error occurred while loading config: ", err.Error())
	}

	db, err := postgres.NewPostgresqlDB(cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.Username,
		cfg.Postgres.DBName, cfg.Postgres.Password, cfg.Postgres.SSLMode)
	if err != nil {
		log.Fatal("Error occurred while loading DB: ", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)

	h := &handler.Handler{
		Services: services,
	}
	mux := h.InitRoutes()

	srv := server.NewServer(cfg, mux)
	ticker := time.NewTicker(interval)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		if err := srv.Run(); err != nil {
			log.Println("error happened: ", err.Error())
		}
	}()

	// TODO: refactor
	go func() {
		for {
			select {
			case <-ticker.C:
				if err := services.DeleteExpiredSegments(); err != nil {
					log.Println("error happened: ", err.Error())
				}
				log.Println("database was updated")
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	log.Println("Application is running")

	<-quit

	log.Println("Application is shutting down")

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Printf("error occurred on server shutting down: %s", err.Error())
	}

	if err := db.Close(); err != nil {
		log.Printf("error occurred on db connection close: %s", err.Error())
	}
}
