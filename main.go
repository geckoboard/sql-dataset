package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/geckoboard/sql-dataset/drivers"
	"github.com/geckoboard/sql-dataset/models"
)

const errMsg = "Dataset '%s' errored: %s\n"

var (
	configFile = flag.String("config", "sql-dataset.yml", "Config file to load")
	client     *Client
)

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

	// Build the connection string
	dc := config.DatabaseConfig
	b, err := drivers.NewConnStringBuilder(dc.Driver)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	s, err := b.Build(dc)

	if err != nil {
		fmt.Println("Error occurred building connection string:", err)
		os.Exit(1)
	}

	dc.URL = s
	client = NewClient(config.GeckoboardAPIKey)

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
			fmt.Printf(errMsg, ds.Name, err)
			hasErrored = true
			continue
		}

		err = client.FindOrCreateDataset(&ds)
		if err != nil {
			fmt.Println(errMsg, ds.Name, err)
			hasErrored = true
			continue
		}

		err = client.SendAllData(&ds, datasetRecs)
		if err != nil {
			fmt.Printf(errMsg, ds.Name, err)
			hasErrored = true
			continue
		}

		fmt.Printf("Dataset '%s' successfully completed\n", ds.Name)
	}

	return hasErrored
}
