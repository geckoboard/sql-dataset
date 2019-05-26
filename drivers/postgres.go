package drivers

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/geckoboard/sql-dataset/models"
)

type postgres struct{}

const (
	postgresPort = "5432"
)

/*
SSL Supported Modes
https://github.com/lib/pq/blob/068cb1c8e4be77b9bdef4d0d91f162160537779e/doc.go

disable - No SSL
require - Always SSL (skip verification)
verify-ca - Always SSL (verify server cert trusted CA)
verify-full - Same as verify-ca and server host on cert matches
*/

func (p postgres) Build(dc *models.DatabaseConfig) (string, error) {
	var buf bytes.Buffer

	if dc.Database == "" {
		return "", errDatabaseRequired
	}

	if dc.Username == "" {
		return "", errUsernameRequired
	}

	p.setDefaults(dc)
	p.buildTLSParams(dc)

	keys := orderKeys(dc.Params)
	for _, k := range keys {
		buf.WriteString(fmt.Sprintf(" %s=%s", k, p.Encode(dc.Params[k])))
	}

	var params string
	if buf.Len() > 0 {
		params = buf.String()
	}

	return fmt.Sprintf("dbname=%s %s%s",
		dc.Database,
		p.buildConnString(dc),
		params,
	), nil
}

func (p postgres) buildTLSParams(dc *models.DatabaseConfig) {
	tc := dc.TLSConfig

	if tc == nil {
		return
	}

	if tc.SSLMode != "" {
		dc.Params["sslmode"] = tc.SSLMode
	}

	if tc.CAFile != "" {
		dc.Params["sslrootcert"] = tc.CAFile
	}

	if tc.KeyFile != "" {
		dc.Params["sslkey"] = tc.KeyFile
	}

	if tc.CertFile != "" {
		dc.Params["sslcert"] = tc.CertFile
	}
}

func (p postgres) buildConnString(dc *models.DatabaseConfig) string {
	var auth, netHost string

	if dc.Password == "" {
		auth = "user=" + dc.Username
	} else {
		auth = fmt.Sprintf("user=%s password=%s", p.Encode(dc.Username), p.Encode(dc.Password))
	}

	if dc.Protocol == tcpConn {
		netHost = fmt.Sprintf("host=%s port=%s", dc.Host, dc.Port)
	} else {
		netHost = fmt.Sprintf("host=%s", dc.Host)
	}

	return auth + " " + netHost
}

func (p postgres) setDefaults(dc *models.DatabaseConfig) {
	if dc.Params == nil {
		dc.Params = make(map[string]string)
	}

	if dc.Host == "" {
		dc.Host = defaultHost
	}

	if dc.Protocol == "" {
		dc.Protocol = tcpConn
	}

	if dc.Port == "" && dc.Protocol == tcpConn {
		dc.Port = postgresPort
	}
}

func (p postgres) Encode(s string) string {
	var changed bool
	new := s

	if strings.Contains(s, `\`) {
		new = strings.Replace(new, `\`, `\\`, -1)
		changed = true
	}

	if strings.Contains(s, " ") {
		new = strings.Replace(new, " ", `\ `, -1)
		changed = true
	}

	if strings.Contains(s, "'") {
		new = strings.Replace(new, "'", `\'`, -1)
		changed = true
	}

	if changed {
		return fmt.Sprintf("'%s'", new)
	}

	return s

}
