package models

import (
	"database/sql"
	"fmt"
	"time"

	"gopkg.in/guregu/null.v3"

	_ "github.com/denisenkom/go-mssqldb"
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
func (ds Dataset) BuildDataset(dc *DatabaseConfig, db *sql.DB) (DatasetRows, error) {
	datasetRecs := DatasetRows{}
	recs, err := ds.queryDatasource(dc, db)

	if err != nil {
		return nil, err
	}

	for _, row := range recs {
		data := make(map[string]interface{})

		for i, col := range row.([]interface{}) {
			f := ds.Fields[i]
			k := f.KeyValue()

			switch f.Type {
			case NumberType, MoneyType, PercentageType:
				data[k] = col.(*Number).Value(f.Optional)
			case StringType:
				data[k] = col.(*null.String).String
			case DateType:
				d := col.(*null.Time)
				if d.Valid {
					data[k] = d.Time.Format(dateFormat)
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

func (ds Dataset) queryDatasource(dc *DatabaseConfig, db *sql.DB) (records []interface{}, err error) {
	rows, err := db.Query(ds.SQL)

	if err != nil {
		return nil, fmt.Errorf(errFailedSQLQuery, err)
	}

	defer rows.Close()

	for rows.Next() {
		var rvp []interface{}
		for _, v := range ds.Fields {
			rvp = append(rvp, v.fieldTypeMapping())
		}

		err = rows.Scan(rvp...)

		if err != nil {
			return nil, fmt.Errorf(errParseSQLResultSet, err)
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
	case NumberType, MoneyType, PercentageType:
		var x Number
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
