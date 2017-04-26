package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"

	"github.com/geckoboard/sql-dataset/models"
	"github.com/go-sql-driver/mysql"
)

const mysqlTLSKey = "customCert"

func registerMysqlTLSConfig(keyFile, certFile, caFile string) (*x509.CertPool, []tls.Certificate, error) {
	rootCertPool := x509.NewCertPool()

	pem, err := ioutil.ReadFile(caFile)
	if err != nil {
		return nil, nil, err
	}

	if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
		return nil, nil, fmt.Errorf("Failed to append PEM, is it a valid ca cert ?")
	}

	clientCert := make([]tls.Certificate, 0, 1)
	if certFile != "" && keyFile != "" {
		certs, err := tls.LoadX509KeyPair(certFile, caFile)
		if err != nil {
			return nil, nil, fmt.Errorf("Error loading x509 key pair: %s", err)
		}

		clientCert = append(clientCert, certs)
	}

	return rootCertPool, clientCert, nil
}

// ConfigureMySQLDSN does some extra setup such as parseTime param
// and configuring TLS for mysql
func ConfigureMySQLDSN(c *models.Config) error {
	conf, err := mysql.ParseDSN(c.DatabaseConfig.URL)
	if err != nil {
		return fmt.Errorf("Failed to parse dsn: %s", err)
	}

	// This needs to be true to parse time
	conf.ParseTime = true

	if c.DatabaseConfig.TLSConfig != nil {
		tlsConf := c.DatabaseConfig.TLSConfig

		rootCertPool, clientCert, err := registerMysqlTLSConfig(
			tlsConf.KeyFile,
			tlsConf.CertFile,
			tlsConf.CAFile,
		)

		if err != nil {
			return err
		}

		mysql.RegisterTLSConfig(mysqlTLSKey, &tls.Config{
			RootCAs:      rootCertPool,
			Certificates: clientCert,
		})

		conf.TLSConfig = mysqlTLSKey
	}

	c.DatabaseConfig.URL = conf.FormatDSN()

	return nil
}
