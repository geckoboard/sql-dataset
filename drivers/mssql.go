package drivers

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/geckoboard/sql-dataset/models"
)

type mssql struct{}

const (
	mssqlPort = "1433"
)

/*
SSL Supported Modes
https://github.com/denisenkom/go-mssqldb/blob/master/README.md

disable - Data send between client and server is not encrypted.
false - Data sent between client and server is not encrypted beyond the login packet. (Default)
true - Data sent between client and server is encrypted.
*/

func (ms mssql) Build(dc *models.DatabaseConfig) (string, error) {
	ms.setDefaults(dc)

	if dc.Database == "" {
		return "", errDatabaseRequired
	}

	// It might be possible to support Windows single sign on
	// however this means username can be empty. For now lets not support
	// not sure what is involved - I think it needs SPN (Kerberos) :(
	if dc.Username == "" {
		return "", errUsernameRequired
	}

	if err := ms.buildTLSParams(dc); err != nil {
		return "", err
	}

	return ms.buildConnString(dc), nil
}

func (ms mssql) buildConnString(dc *models.DatabaseConfig) string {
	var buf bytes.Buffer
	var password string

	// Shouldn't be the case with password policies
	if dc.Password != "" {
		password = fmt.Sprintf("password={%s}", dc.Password)
	}

	keys := orderKeys(dc.Params)
	for _, k := range keys {
		ms.buildParams(&buf, fmt.Sprintf("%s=%s", k, dc.Params[k]))
	}

	conn := fmt.Sprintf("odbc:server={%s};port=%s;user id={%s};%s;database=%s",
		dc.Host, dc.Port, dc.Username, password, dc.Database)

	if buf.Len() > 0 {
		conn = fmt.Sprintf(conn+";%s", buf.String())
	}

	return conn
}

func (ms mssql) buildTLSParams(dc *models.DatabaseConfig) error {
	tc := dc.TLSConfig

	if tc == nil {
		return nil
	}

	if tc.SSLMode != "" {
		dc.Params["encrypt"] = tc.SSLMode
	}

	if tc.CAFile != "" {
		dc.Params["certificate"] = tc.CAFile
	}

	if tc.KeyFile != "" {
		return errors.New("Key file not supported, only ca_file is for MSSQL Driver")
	}

	if tc.CertFile != "" {
		return errors.New("Cert file not supported, only ca_file is for MSSQL Driver")
	}

	return nil
}

func (ms mssql) setDefaults(dc *models.DatabaseConfig) {
	if dc.Params == nil {
		dc.Params = make(map[string]string)
	}

	if dc.Host == "" {
		dc.Host = defaultHost
	}

	if dc.Port == "" {
		dc.Port = mssqlPort
	}

	// Set the application intent to readonly connection
	dc.Params["ApplicationIntent"] = "ReadOnly"
}

/* buildParams builds a buffer of query parameters, specified by the
   user or system defaults, seperated by semi-colon in the following format.

   ApplicationIntent=ReadOnly;connection timeout=10;dial timeout=2
*/
func (ms mssql) buildParams(buf *bytes.Buffer, str string) {
	if str == "=" {
		return
	}

	if buf.Len() > 0 {
		buf.WriteString(";")
	}

	buf.WriteString(str)
}
