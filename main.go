package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/geckoboard/sql-dataset/models"
)

var configFile = flag.String("config", "sql-dataset.yml", "Config file to load")

func main() {
	flag.Parse()

	config, err := models.LoadConfig(*configFile)

	if err != nil {
		log.Fatal(err)
	}

	if errs := config.Validate(); errs != nil {
		fmt.Println("\nFollowing errors occurred with the config;")

		for i, err := range errs {
			fmt.Println(" -", err)

			if i == len(errs)-1 {
				fmt.Println("")
			}
		}

		os.Exit(1)
	}
}
