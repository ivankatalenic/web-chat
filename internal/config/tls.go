package config

import "os"

type tlsConfig struct {
	CertFilePath string
	KeyFilePath string
}

// TLS holds the configuration for the TLS connection
var TLS = tlsConfig{
	CertFilePath: os.Getenv("TLS_CERT_PATH"),
	KeyFilePath:  os.Getenv("TLS_KEY_PATH"),
}
