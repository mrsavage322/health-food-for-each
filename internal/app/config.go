package app

import (
	"flag"
	"github.com/caarlos0/env/v6"
)

var ConnectionDB DBConnect
var ServerAddress string
var Config ConfigApp

type ConfigApp struct {
	ServerAddress   string `env:"SERVER_ADDRESS"`
	DatabaseAddress string `env:"DATABASE_DSN"`
	ConnectionDB    DBConnect
}

func SetConfig() {
	err := env.Parse(&Config)
	if err != nil {
		return
	}
}

func SetFlags() {
	flag.StringVar(&Config.ServerAddress, "a", "localhost:8080", "Address to run the HTTP server")
	flag.StringVar(&Config.DatabaseAddress, "d", "", "Address to connect with Database")
	flag.Parse()
}
