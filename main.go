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

var (
	configFile     = flag.String("config", "sql-dataset.yml", "Config file to load")
	displayVersion = flag.Bool("version", false, "Displays version info")
	version        = ""
	gitSHA         = ""
)

func main() {
	flag.Parse()

	if *displayVersion {
		fmt.Printf("Version: %s\nGitSHA: %s\n", version, gitSHA)
		os.Exit(0)
	}

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
	client := NewClient(config.GeckoboardAPIKey)

	if config.RefreshTimeSec == 0 {
		fmt.Printf("No refresh timer specified will process once and exit\n")
		processAllDatasets(config, client)
	} else {
		fmt.Printf("Refresh timer specified run every %d seconds until interrupted\n\n", config.RefreshTimeSec)
		for {
			processAllDatasets(config, client)
			time.Sleep(time.Duration(config.RefreshTimeSec) * time.Second)
		}
	}
}

func processAllDatasets(config *models.Config, client *Client) (hasErrored bool) {
	for _, ds := range config.Datasets {
		datasetRecs, err := ds.BuildDataset(config.DatabaseConfig)
		if err != nil {
			printErrorMsg(ds.Name, err)
			hasErrored = true
			continue
		}

		err = client.FindOrCreateDataset(&ds)
		if err != nil {
			printErrorMsg(ds.Name, err)
			hasErrored = true
			continue
		}

		err = client.SendAllData(&ds, datasetRecs)
		if err != nil {
			printErrorMsg(ds.Name, err)
			hasErrored = true
			continue
		}

		fmt.Printf("Dataset '%s' successfully completed\n", ds.Name)
	}

	return hasErrored
}

func printErrorMsg(name string, err error) {
	fmt.Printf("Dataset '%s' errored: %s\n", name, err)
}
