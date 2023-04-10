package config

import (
	"crypto/tls"
	"crypto/x509"
	"os"
)

var (
	DebugMode  bool   = false
	LogEncoder string = "plain"
	ConfigFile string = ".kpeering/config.toml"

	DefautlTunnelPort = 10022
)

func LoadServerTlsConfig(serverCertPath, serverKeyPath, serverName string) (*tls.Config, error) {
	crt, err := os.ReadFile(serverCertPath)
	if err != nil {
		return nil, err
	}

	key, err := os.ReadFile(serverKeyPath)
	if err != nil {
		return nil, err
	}

	return CreateServerTlsConfig(crt, key, serverName)
}

func CreateServerTlsConfig(serverCert []byte, serverKey []byte, serverName string) (*tls.Config, error) {
	cert, err := tls.X509KeyPair(serverCert, serverKey)
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		ServerName:   serverName,
	}, nil
}

func LoadClientTlsConfig(caCertPath, serverName string) (*tls.Config, error) {
	crt, err := os.ReadFile(caCertPath)
	if err != nil {
		return nil, err
	}

	return CreateClientTlsConfig(crt, serverName)
}

func CreateClientTlsConfig(caCert []byte, serverName string) (*tls.Config, error) {
	rootCAs := x509.NewCertPool()
	rootCAs.AppendCertsFromPEM(caCert)

	return &tls.Config{
		RootCAs:            rootCAs,
		InsecureSkipVerify: false,
		ServerName:         serverName,
	}, nil
}
