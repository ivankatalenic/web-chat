package config

import "os"

// TLSConfig holds paths to the files used for TLS configuration
type TLSConfig struct {
	CertFilePath string
	KeyFilePath string
}

// TLS holds the configuration for the TLS connection
var TLS = TLSConfig{
	CertFilePath: os.Getenv("TLS_CERT_PATH"),
	KeyFilePath:  os.Getenv("TLS_KEY_PATH"),
}
