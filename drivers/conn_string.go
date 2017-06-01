package drivers

import (
	"bytes"
	"errors"
	"fmt"
	"sort"

	"github.com/geckoboard/sql-dataset/models"
)

const (
	defaultHost = "localhost"
	tcpConn     = "tcp"
)

var (
	ErrDatabaseRequired = errors.New("Database is required for a connection")
	ErrUsernameRequired = errors.New("Username is required for a connection")
)

type ConnStringBuilder interface {
	Build(*models.DatabaseConfig) (string, error)
}

func NewConnStringBuilder(driver string) (ConnStringBuilder, error) {
	switch driver {
	case models.PostgresDriver:
		return postgres{}, nil
	case models.MySQLDriver:
		return mysql{}, nil
	case models.SQLiteDriver:
		return sqlite{}, nil
	case models.MSSQLDriver:
		return mssql{}, nil
	default:
		return nil, fmt.Errorf("Unknown driver %s to build connection string", driver)
	}
}

func buildParams(buf *bytes.Buffer, str string) {
	if str == "" || str == "=" {
		return
	}

	if buf.Len() > 0 {
		buf.WriteString("&")
	}

	buf.WriteString(str)
}

func orderKeys(kv map[string]string) []string {
	var keys []string

	for k := range kv {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	return keys
}
