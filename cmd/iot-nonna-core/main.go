package main

import (
	"log"
	"net/http"
	"os"

	"github.com/chiaf1/iot-nonna-core/internal/config"
	"github.com/chiaf1/iot-nonna-core/internal/db"
	"github.com/chiaf1/iot-nonna-core/internal/handler"
	"github.com/chiaf1/iot-nonna-core/internal/repository"
	"github.com/chiaf1/iot-nonna-core/internal/router"
)

const CONFIG_PATH = "./config.yaml"

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
	log.Fatal(http.ListenAndServe(":3000", r))
}
