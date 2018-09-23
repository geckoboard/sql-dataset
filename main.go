package main

import (
	"bufio"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/geckoboard/sql-dataset/drivers"
	"github.com/geckoboard/sql-dataset/models"
)

var (
	configFile     = flag.String("config", "sql-dataset.yml", "Config file to load")
	deleteDataset  = flag.String("delete-dataset", "", "Pass a dataset name you want to delete")
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
		fmt.Println("\nThere are errors in your config:")

		for i, err := range errs {
			fmt.Println(" -", err)

			if i == len(errs)-1 {
				fmt.Println("")
			}
		}

		os.Exit(1)
	}

	if *deleteDataset != "" {
		if err := deleteDatasetSwitch(*deleteDataset, config); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		os.Exit(0)
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
		fmt.Println("There was an error while trying to build "+
			"your database connection string:", err)
		os.Exit(1)
	}

	client := NewClient(config.GeckoboardAPIKey)
	db, err := newDBConnection(dc.Driver, dsn)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if config.RefreshTimeSec == 0 {
		processAllDatasets(config, client, db)
	} else {
		fmt.Printf("Running every %d seconds, until interrupted.\n\n", config.RefreshTimeSec)
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

		fmt.Printf("Successfully updated \"%s\"\n", ds.Name)
	}

	return hasErrored
}

func printErrorMsg(name string, err error) {
	fmt.Printf("There was an error while trying to update %s: %s\n", name, err)
}

func newDBConnection(driver, url string) (*sql.DB, error) {
	// Ignore this error which just checks we have the driver loaded
	pool, _ := sql.Open(driver, url)
	err := pool.Ping()

	if err != nil {
		return nil, fmt.Errorf("Failed to open database connection. "+
			"This is the error received: %s", err)
	}

	pool.SetMaxOpenConns(5)

	return pool, err
}

func deleteDatasetSwitch(name string, config *models.Config) error {
	fmt.Printf("Delete dataset %q (y/N): ", name)

	v, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return err
	}

	// Remove newline and carriage return for windows
	v = strings.TrimRight(v, "\n")
	v = strings.TrimRight(v, "\r")

	switch strings.ToLower(v) {
	case "y":
		client := NewClient(config.GeckoboardAPIKey)
		if err := client.DeleteDataset(*deleteDataset); err != nil {
			return err
		}

		fmt.Println("Dataset deleted successfully")
	default:
		fmt.Println("Cancelled action")
	}

	return nil
}
