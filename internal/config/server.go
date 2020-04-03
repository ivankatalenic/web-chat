package config

import "os"

type serverConfig struct {
	Host string
}

// Server hold the configuration for a server
var Server = serverConfig{
	Host: os.Getenv("HOST"),
}
