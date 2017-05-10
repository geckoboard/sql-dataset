package drivers

import (
	"bytes"
	"errors"
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

type DSNBuilder interface {
	Build(*models.DatabaseConfig) (string, error)
}

func NewDSNBuilder(driver string) DSNBuilder {
	var builder DSNBuilder

	switch driver {
	case models.PostgresDriver:
		builder = postgres{}
	case models.MysqlDriver:
		builder = mysql{}
	default:
		builder = sqlite{}
	}

	return builder
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
	keys := make([]string, 0, len(kv))

	for k := range kv {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	return keys
}
