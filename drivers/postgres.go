package drivers

import (
	"bytes"
	"fmt"
	"net"

	"github.com/geckoboard/sql-dataset/models"
)

type postgres struct{}

const (
	postgresPort = "5432"
	connPrefix   = "postgres://"
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
		return "", ErrDatabaseRequired
	}

	if dc.Username == "" {
		return "", ErrUsernameRequired
	}

	p.setDefaults(dc)
	p.buildTLSParams(dc)

	keys := orderKeys(dc.Params)
	for _, k := range keys {
		buildParams(&buf, fmt.Sprintf("%s=%s", k, dc.Params[k]))
	}

	var paramSplit string
	if buf.Len() > 0 {
		paramSplit = "?"
	}

	return fmt.Sprintf("%s%s%s%s",
		connPrefix,
		p.buildConnString(dc),
		paramSplit,
		buf.String(),
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
		auth = dc.Username
	} else {
		auth = fmt.Sprintf("%s:%s", dc.Username, dc.Password)
	}

	if dc.Protocol == tcpConn {
		netHost = net.JoinHostPort(dc.Host, dc.Port)
	} else {
		netHost = dc.Host
	}

	return fmt.Sprintf("%s@%s/%s", auth, netHost, dc.Database)
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
