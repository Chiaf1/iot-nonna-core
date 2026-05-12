package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/chiaf1/iot-nonna-core/internal/config"
	"github.com/chiaf1/iot-nonna-core/internal/db"
	"github.com/chiaf1/iot-nonna-core/internal/handler"
	"github.com/chiaf1/iot-nonna-core/internal/repository"
	"github.com/chiaf1/iot-nonna-core/internal/router"
)

const CONFIG_PATH = "./config.yaml"

// @title 			IOT nonna core API
// @version 		1.0
// @description 	API for management of IoT system database
// @host			localhost:3030
// @BasePath		/
func main() {
	confPath := os.Getenv("CONFIG_PATH")
	if confPath == "" {
		confPath = CONFIG_PATH
	}
	// 1. Load configs from file
	var conf config.Config
	err := conf.Load(confPath)
	if err != nil {
		log.Fatal(err)
	}
	err = conf.Validate()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Config Loaded")

	// 2. Connet to db
	// DB connection pool creation
	dbPool, err := db.NewPgPool(conf.DB)
	if err != nil {
		log.Fatalf("Error opening connection to db: %v", err)
	}
	log.Println("Connection established with db")

	// 3. Launch migration of db if enabled
	if os.Getenv("RUN_MIGRATIONS") == "true" {
		if err := db.RunMigrations(conf.DB.DbURL); err != nil {
			log.Fatal(err)
		}
	}
	// 4. Launch seeding of db if enabled
	if os.Getenv("RUN_SEEDING") == "true" {
		if err := db.RunSeeding(dbPool, conf.DB.Query_timeout_read); err != nil {
			log.Fatal(err)
		}
	}

	// 5. Creating the repo and handler structures
	repo := repository.NewRepo(dbPool, conf.DB.Query_timeout_read, conf.DB.Query_timeout_write)
	h := handler.NewHandler(repo)
	log.Println("Repo and Handler created, starting web server...")

	// 6. Chi router creation
	r := router.Setup(h)

	// Creating an HTTP server to handle graceful shutdowns
	srv := http.Server{
		Addr:    ":3030",
		Handler: r,
	}

	// Starting the server in a go routine
	go func() {
		log.Println("Server listening on :3030")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen error: %v", err)
		}
	}()

	//Waiting for shutdown signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutdown signal received")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Println("Server exited cleanly")
}
