package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

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

	if config.DatabaseConfig.Driver == models.MysqlDriver {
		ConfigureMySQLDSN(config)
	}

	if config.RefreshTimeSec == 0 {
		fmt.Println("No refresh timer specified will process once and exit\n")
		processAllDatasets(config)
	} else {
		fmt.Printf("Refresh timer specified run every %d seconds until interrupted\n\n", config.RefreshTimeSec)
		for {
			processAllDatasets(config)
			time.Sleep(time.Duration(config.RefreshTimeSec) * time.Second)
		}
	}
}

func processAllDatasets(config *models.Config) (hasErrored bool) {
	for _, ds := range config.Datasets {
		datasetRecs, err := ds.BuildDataset(config.DatabaseConfig)
		if err != nil {
			fmt.Printf("Dataset '%s' errored: %s\n", ds.Name, err)
			hasErrored = true
			continue
		}

		err = PushData(ds, datasetRecs, config.GeckoboardAPIKey)
		if err != nil {
			fmt.Printf("Dataset '%s' errored: %s\n", ds.Name, err)
			hasErrored = true
			continue
		}

		fmt.Printf("Dataset '%s' successfully completed\n", ds.Name)
	}

	return hasErrored
}
