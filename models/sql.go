package models

import (
	"database/sql"
	"fmt"
	"time"

	"gopkg.in/guregu/null.v3"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

const dateFormat = "2006-01-02"

// DatasetRows holds a slice of map[string]interface{}
// which is used to send to a geckoboard dataset
type DatasetRows []map[string]interface{}

// BuildDataset calls queryDatasource to query the datasource for a
// dataset entry and builds up a slice of rows ready for processing by the client
func (ds Dataset) BuildDataset(dc *DatabaseConfig) (DatasetRows, error) {
	var datasetRecs DatasetRows
	recs, err := ds.queryDatasource(dc)

	if err != nil {
		return nil, err
	}

	//TODO: Allow nulls on the datatypes that support it in Geckoboard
	for _, row := range recs {
		data := make(map[string]interface{})

		for i, col := range row.([]interface{}) {
			f := ds.Fields[i]
			k := f.KeyValue()

			switch f.Type {
			case NumberType:
				if f.FloatPrecision != 0 {
					val := col.(*null.Float).Float64
					if f.FloatPrecision == 32 {
						data[k] = float32(val)
					} else {
						data[k] = val
					}
				} else {
					data[k] = col.(*null.Int).Int64
				}
			case MoneyType:
				data[k] = col.(*null.Int).Int64
			case PercentageType:
				val := col.(*null.Float).Float64

				if f.FloatPrecision == 32 {
					data[k] = float32(val)
				} else {
					data[k] = val
				}
			case StringType:
				data[k] = col.(*null.String).String
			case DateType:
				d := col.(*null.Time)
				if d.Valid {
					data[k] = d.Time.Format("2006-01-02")
				} else {
					data[k] = nil
				}
			case DatetimeType:
				d := col.(*null.Time)
				if d.Valid {
					data[k] = d.Time.Format(time.RFC3339)
				} else {
					data[k] = nil
				}
			}
		}

		datasetRecs = append(datasetRecs, data)
	}

	return datasetRecs, nil
}

func (ds Dataset) queryDatasource(dc *DatabaseConfig) (records []interface{}, err error) {
	db, err := sql.Open(dc.Driver, dc.URL)
	if err != nil {
		return nil, fmt.Errorf("Database open failed: %s", err)
	}

	rows, err := db.Query(ds.SQL)

	if err != nil {
		return nil, fmt.Errorf("Database query failed: %s", err)
	}

	defer rows.Close()

	for rows.Next() {
		var rvp []interface{}
		for _, v := range ds.Fields {
			rvp = append(rvp, v.fieldTypeMapping())
		}

		err = rows.Scan(rvp...)

		if err != nil {
			return nil, fmt.Errorf("Scan failed: %s", err)
		}

		records = append(records, rvp)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return records, nil
}

func (f Field) fieldTypeMapping() interface{} {
	switch f.Type {
	case NumberType:
		if f.FloatPrecision != 0 {
			var x null.Float
			return &x
		}
		var x null.Int
		return &x
	case MoneyType:
		var x null.Int
		return &x
	case PercentageType:
		var x null.Float
		return &x
	case StringType:
		var x null.String
		return &x
	case DateType, DatetimeType:
		var x null.Time
		return &x
	}

	return nil
}
