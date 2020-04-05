package config

import (
	"github.com/ivankatalenic/web-chat/internal/helpers"
)

type serverConfig struct {
	Host string
}

// Server hold the configuration for a server
var Server = serverConfig{
	Host: helpers.GetEnvVarOrDefault("HOST", "localhost"),
}
