package config

import (
	"flag"
	"log"
	"os"
)

var (
	DefaultSA  = ":8080"
	DefaultBU  = "http://localhost:8080"
	DefaultFSP = ""
)

// declare config struct
type URLConfig struct {
	SA  string //SERVER_ADDRESS
	BU  string //BASE_URL
	FSP string //FILE_STORAGE_PATH
}

// new confit struct instance
var UC = URLConfig{}

func getSAfromEnv() {
	// get sa from env if empty
	if UC.SA != "" {
		return
	}
	// if empty set default
	UC.SA = DefaultSA
	if s, ok := os.LookupEnv("SERVER_ADDRESS"); ok {
		UC.SA = s
	}
}

func getBUfromEnv() {
	// get bu from env if empty
	if UC.BU != "" {
		return
	}
	UC.BU = DefaultBU
	if s, ok := os.LookupEnv("BASE_URL"); ok {
		UC.BU = s
	}
}

func getFSPfromEnv() {
	// get fsp from env if empty
	if UC.FSP != "" {
		return
	}
	UC.FSP = DefaultFSP
	if s, ok := os.LookupEnv("FILE_STORAGE_PATH"); ok {
		UC.FSP = s
	}
}

func BuildConfig() {
	// get config vars from CLI flags
	flag.Parse()
	log.Println("config/BuildConfig UC after flags", UC)

	// get config from local var if was not set by flag
	getSAfromEnv()
	getBUfromEnv()
	getFSPfromEnv()
	log.Println("config/BuildConfig UC after env vars", UC)
}
