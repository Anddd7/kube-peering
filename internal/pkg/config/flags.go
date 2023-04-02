package config

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
)

var (
	DebugMode  bool   = false
	LogEncoder string = "plain"

	DefautlTunnelPort = 10022
)

func LoadServerTlsConfig(serverCertPath, serverKeyPath, serverName string) (*tls.Config, error) {
	crt, err := ioutil.ReadFile(serverCertPath)
	if err != nil {
		return nil, err
	}

	key, err := ioutil.ReadFile(serverKeyPath)
	if err != nil {
		return nil, err
	}

	cert, err := tls.X509KeyPair(crt, key)
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		ServerName:   serverName,
	}, nil
}

func LoadClientTlsConfig(caCertPath, serverName string) (*tls.Config, error) {
	crt, err := ioutil.ReadFile(caCertPath)
	if err != nil {
		return nil, err
	}

	rootCAs := x509.NewCertPool()
	rootCAs.AppendCertsFromPEM(crt)

	return &tls.Config{
		RootCAs:            rootCAs,
		InsecureSkipVerify: false,
		ServerName:         serverName,
	}, nil
}
