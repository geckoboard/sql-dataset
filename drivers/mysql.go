package drivers

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"

	"github.com/geckoboard/sql-dataset/models"
	msql "github.com/go-sql-driver/mysql"
)

/*
SSL Supported Modes
https://github.com/go-sql-driver/mysql#tls

false - no SSL (default)
true - use ssl connection
skip-verify - use self-signed or invalid server cert
customCert - name of the registered tls config (automatic when supplying ssl (ca,cert,key)
*/

type mysql struct{}

const (
	mysqlTLSKey = "customCert"
	mysqlPort   = "3306"
)

func (m mysql) Build(dc *models.DatabaseConfig) (string, error) {
	var buf bytes.Buffer

	m.setDefaults(dc)

	if dc.Database == "" {
		return "", errDatabaseRequired
	}

	if dc.Username == "" {
		return "", errUsernameRequired
	}

	if dc.TLSConfig != nil {
		str, err := m.registerTLS(dc.TLSConfig)

		if err != nil {
			return "", err
		}

		dc.Params["tls"] = str
	}

	keys := orderKeys(dc.Params)
	for _, k := range keys {
		buildParams(&buf, fmt.Sprintf("%s=%s", k, dc.Params[k]))
	}

	return fmt.Sprintf("%s?%s", m.buildConnString(dc), buf.String()), nil
}

func (m mysql) loadCerts(keyFile, certFile, caFile string) (*x509.CertPool, []tls.Certificate, error) {
	var rootCertPool *x509.CertPool

	if caFile != "" {
		rootCertPool = x509.NewCertPool()

		pem, err := ioutil.ReadFile(caFile)
		if err != nil {
			return nil, nil, err
		}

		if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
			return nil, nil, errAppendPEMFailed
		}
	}

	var clientCert []tls.Certificate

	if certFile != "" && keyFile != "" {
		certs, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			return nil, nil, fmt.Errorf("There was an error while "+
				"loading your x509 key pair: %s", err)
		}

		clientCert = append(clientCert, certs)
	}

	return rootCertPool, clientCert, nil
}

func (m mysql) registerTLS(tlsConfig *models.TLSConfig) (string, error) {
	if tlsConfig.KeyFile == "" && tlsConfig.CertFile == "" &&
		tlsConfig.CAFile == "" && tlsConfig.SSLMode != "" {
		return tlsConfig.SSLMode, nil
	}

	rootCertPool, clientCert, err := m.loadCerts(
		tlsConfig.KeyFile,
		tlsConfig.CertFile,
		tlsConfig.CAFile,
	)

	if err != nil {
		return "", err
	}

	return mysqlTLSKey, msql.RegisterTLSConfig(mysqlTLSKey, &tls.Config{
		RootCAs:      rootCertPool,
		Certificates: clientCert,
	})
}

func (m mysql) buildConnString(dc *models.DatabaseConfig) string {
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

	return fmt.Sprintf("%s@%s(%s)/%s", auth, dc.Protocol, netHost, dc.Database)
}

func (m mysql) setDefaults(dc *models.DatabaseConfig) {
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
		dc.Port = mysqlPort
	}

	dc.Params["parseTime"] = "true"
}
