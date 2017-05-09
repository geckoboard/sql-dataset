package drivers

import (
	"bytes"
	"fmt"

	"github.com/geckoboard/sql-dataset/models"
)

type sqlite struct{}

func (s sqlite) Build(dc *models.DatabaseConfig) (dsn string, err error) {
	var buf bytes.Buffer

	if dc.Database == "" {
		return "", ErrDatabaseRequired
	}

	if dc.Password != "" {
		buildParams(&buf, fmt.Sprintf("password=%s", dc.Password))
	}

	keys := orderKeys(dc.Params)
	for _, k := range keys {
		buildParams(&buf, fmt.Sprintf("%s=%s", k, dc.Params[k]))
	}

	if buf.Len() > 0 {
		dsn = fmt.Sprintf("file:%s?%s", dc.Database, buf.String())
	} else {
		dsn = dc.Database
	}

	return dsn, err
}
