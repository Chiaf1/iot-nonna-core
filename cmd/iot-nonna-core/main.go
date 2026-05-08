package main

import (
	"log"
	"os"

	"github.com/chiaf1/iot-nonna-core/internal/config"
	"github.com/chiaf1/iot-nonna-core/internal/db"
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
	_, err = db.NewPgPool(conf.DB)
	if err != nil {
		log.Fatalf("Error opening connection to db: %v", err)
	}
	log.Println("Connection established with db")

}
