package drivers

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/url"

	chouse "github.com/ClickHouse/clickhouse-go"
	"github.com/geckoboard/sql-dataset/models"
)

type clickhouse struct{}

const (
	clickhouseTLSKey = "clickhouseCert"
	clickhousePort   = "9000"
)

func (ch clickhouse) Build(dc *models.DatabaseConfig) (conn string, err error) {
	link := url.URL{
		Scheme: "tcp",
	}
	if dc.Host == "" {
		return "", errHostnameRequired
	}
	port := clickhousePort
	if dc.Port != "" {
		port = dc.Port
	}
	link.Host = dc.Host + ":" + port
	params := url.Values{}
	if dc.Username != "" {
		params.Add("username", dc.Username)
		if dc.Password != "" {
			params.Add("password", dc.Password)
		}
	}
	if dc.Database != "" {
		params.Add("database", dc.Database)
	}
	if dc.TLSConfig != nil {
		tlsName, err := ch.registerTLS(dc.TLSConfig)
		if err != nil {
			return "", err
		}
		params.Add("secure", "true")
		params.Add("tls_config", tlsName)
	}
	for k, v := range dc.Params {
		params.Add(k, v)
	}
	if len(params) > 0 {
		return link.String() + "?" + params.Encode(), nil
	}
	return link.String(), nil
}

func (ch clickhouse) registerTLS(tlsConfig *models.TLSConfig) (string, error) {
	rootCertPool, clientCert, err := ch.loadCerts(
		tlsConfig.KeyFile,
		tlsConfig.CertFile,
		tlsConfig.CAFile,
	)
	if err != nil {
		return "", err
	}
	return clickhouseTLSKey, chouse.RegisterTLSConfig(clickhouseTLSKey, &tls.Config{
		RootCAs:      rootCertPool,
		Certificates: clientCert,
	})
}

func (ch clickhouse) loadCerts(keyFile, certFile, caFile string) (*x509.CertPool, []tls.Certificate, error) {
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
