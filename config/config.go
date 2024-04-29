package config

import (
	"os"
)

type Config struct {
	Port         string
	ReadonlyDash bool
}

func newConfig() Config {
	config := Config{
		Port: "8080",
	}
	envPort := os.Getenv("PORT")
	_, ReadonlyDash := os.LookupEnv("READONLY_DASH")
	if len(envPort) > 0 {
		config.Port = envPort
	}
	config.ReadonlyDash = ReadonlyDash
	return config
}

var CLIConfig = newConfig()
