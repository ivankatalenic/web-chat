package config

import (
	"github.com/ivankatalenic/web-chat/internal/helpers"
)

type auth struct {
	Username string
	Password string
}

// Auth variable holds the configuration used for the authorization
var Auth = auth{
	Username: helpers.GetEnvVarOrDefault("USER", "user"),
	Password: helpers.GetEnvVarOrDefault("PASSWORD", "SecurePass1234"),
}
