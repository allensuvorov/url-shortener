package config

import (
	"flag"
	"log"
	"os"
)

// default values
var (
	Dsa  = ":8080"
	Dbu  = "http://localhost:8080"
	Dfsp = ""
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
	UC.SA = Dsa
	if s, ok := os.LookupEnv("SERVER_ADDRESS"); ok {
		UC.SA = s
	}
}

func getBUfromEnv() {
	// get bu from env if empty
	if UC.BU != "" {
		return
	}
	UC.BU = Dbu
	if s, ok := os.LookupEnv("BASE_URL"); ok {
		UC.BU = s
	}
}

func getFSPfromEnv() {
	// get fsp from env if empty
	if UC.FSP != "" {
		return
	}
	UC.FSP = Dfsp
	if s, ok := os.LookupEnv("FILE_STORAGE_PATH"); ok {
		UC.FSP = s
	}
}

func BuildConfig() {
	// get config vars from CLI flags
	flag.Parse()
	log.Println("config/BuildConfig", UC)

	// get config from local var if was not set by flag
	getSAfromEnv()
	getBUfromEnv()
	getFSPfromEnv()
}
