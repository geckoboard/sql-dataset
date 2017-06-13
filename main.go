package main

import (
	"database/sql"
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

	dsn, err := b.Build(dc)

	if err != nil {
		fmt.Println("Error occurred building connection string:", err)
		os.Exit(1)
	}

	client := NewClient(config.GeckoboardAPIKey)
	db, err := newDBConnection(dc.Driver, dsn)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if config.RefreshTimeSec == 0 {
		fmt.Printf("No refresh timer specified will process once and exit\n")
		processAllDatasets(config, client, db)
	} else {
		fmt.Printf("Refresh timer specified run every %d seconds until interrupted\n\n", config.RefreshTimeSec)
		for {
			processAllDatasets(config, client, db)
			time.Sleep(time.Duration(config.RefreshTimeSec) * time.Second)
		}
	}
}

func processAllDatasets(config *models.Config, client *Client, db *sql.DB) (hasErrored bool) {
	for _, ds := range config.Datasets {
		datasetRecs, err := ds.BuildDataset(config.DatabaseConfig, db)
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

func newDBConnection(driver, url string) (*sql.DB, error) {
	pool, err := sql.Open(driver, url)

	if err != nil {
		return nil, fmt.Errorf("Database open failed: %s", err)
	}

	pool.SetMaxOpenConns(5)

	return pool, err
}
