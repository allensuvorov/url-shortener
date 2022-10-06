package config

import (
	"log"
	"os"
)

type Config struct {
}

func NewURLConfig() Config {
	return Config{}
}

func (c Config) GetBU() string {
	// get BU from local env
	log.Println("Config/GetBU: about to take BASE_URL from local env")
	bu, ok := os.LookupEnv("BASE_URL")
	if ok {
		log.Println("Config/GetBU: BASE_URL from local env is:", bu)
	}
	if !ok {
		log.Printf("Config/GetBU: %s not set\n; passing default", "BASE_URL")
		bu = "http://localhost:8080"
	}

	return bu
}
