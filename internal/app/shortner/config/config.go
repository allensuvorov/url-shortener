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

const (
	sa  = "SERVER_ADDRESS"
	bu  = "BASE_URL"
	fsp = "FILE_STORAGE_PATH"
)

type URLConfig struct {
	SA  string
	BU  string
	FSP string
}

var UC = URLConfig{}

func getSAfromEnv() {
	if UC.SA != "" {
		return
	}
	UC.SA = DefaultSA
	if s, ok := os.LookupEnv(sa); ok {
		UC.SA = s
	}
}

func getBUfromEnv() {
	if UC.BU != "" {
		return
	}
	UC.BU = DefaultBU
	if s, ok := os.LookupEnv(bu); ok {
		UC.BU = s
	}
}

func getFSPfromEnv() {
	if UC.FSP != "" {
		return
	}
	UC.FSP = DefaultFSP
	if s, ok := os.LookupEnv(fsp); ok {
		UC.FSP = s
	}
}

func BuildConfig() {
	flag.Parse()
	log.Println("config/BuildConfig UC after flags", UC)

	getSAfromEnv()
	getBUfromEnv()
	getFSPfromEnv()
	log.Println("config/BuildConfig UC after env vars", UC)
}
