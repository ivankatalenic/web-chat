package config

import "os"

type auth struct {
	Username string
	Password string
}

// Auth variable holds the configuration used for the authorization
var Auth = auth{
	Username: os.Getenv("USERNAME"),
	Password: os.Getenv("PASSWORD"),
}
