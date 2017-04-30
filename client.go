package main

import (
	"fmt"

	gb "github.com/geckoboard/geckoboard-go"
	"github.com/geckoboard/sql-dataset/models"
)

var (
	gbClient  *gb.Client
	batchRows = 500

	errMoreRowsToSend = "Sent only the first %d rows, %d rows existed " +
		"to support sending more change dataset update type " +
		"from 'replace' to 'append' to support upto 5000 rows"
)

func PushData(ds models.Dataset, rows models.DatasetRows, key string) (err error) {
	if gbClient == nil {
		gbClient = gb.New(gb.Config{Key: key})
	}

	// Create & push dataset schema
	dataset := gb.DataSet{
		ID:       ds.Name,
		UniqueBy: ds.UniqueBy,
		Fields:   remapFields(ds),
	}

	err = dataset.FindOrCreate(gbClient)
	if err != nil {
		return err
	}

	return batchRequests(ds.UpdateType, dataset, rows)
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

func batchRequests(updateType models.DatasetType, dataset gb.DataSet, rows models.DatasetRows) (err error) {
	switch updateType {
	case models.Replace:
		if len(rows) > batchRows {
			err = dataset.SendAll(gbClient, rows[0:batchRows])
			if err == nil {
				// Error that there were more rows to send
				err = fmt.Errorf(errMoreRowsToSend, batchRows, len(rows))
			}
		} else {
			err = dataset.SendAll(gbClient, rows[:len(rows)])
		}
	case models.Append:
		grps := len(rows) / batchRows

		for i := 0; i <= grps; i++ {
			if i == grps {
				if (batchRows*i)+1 <= len(rows) {
					err = dataset.Append(gbClient, rows[batchRows*i:])
				}
			} else {
				err = dataset.Append(gbClient, rows[batchRows*i:batchRows*(i+1)])
			}
		}
	}

	return err
}
