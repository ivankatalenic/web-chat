package main

import "os"

type serverConfig struct {
	Host string
}

var ServerConfig = serverConfig{
	Host: os.Getenv("HOST"),
}
