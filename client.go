package main

import (
	gb "github.com/geckoboard/geckoboard-go"
	"github.com/geckoboard/sql-dataset/models"
)

var gbClient *gb.Client

func PushData(ds models.Dataset, rows models.DatasetRows, key string) (err error) {
	if gbClient == nil {
		gbClient = gb.New(gb.Config{Key: key})
	}

	// Create & push dataset schema
	dataset := gb.DataSet{
		ID:     ds.Name,
		Fields: remapFields(ds),
	}

	err = dataset.FindOrCreate(gbClient)
	if err != nil {
		return err
	}

	// Push dataset data based on the update type
	if ds.UpdateType == models.Replace {
		err = dataset.SendAll(gbClient, rows)
	} else {
		err = dataset.Append(gbClient, rows)
	}

	if err != nil {
		return err
	}

	return nil
}

func remapFields(ds models.Dataset) (fields gb.Fields) {
	fields = make(gb.Fields)

	for _, f := range ds.Fields {
		fields[f.KeyValue()] = gb.Field{
			Name:         f.Name,
			Type:         string(f.Type),
			CurrencyCode: f.CurrencyCode,
		}
	}

	return fields
}
