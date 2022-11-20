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
	DefaultDSN = ""
)

const (
	sa  = "SERVER_ADDRESS"
	bu  = "BASE_URL"
	fsp = "FILE_STORAGE_PATH"
	dsn = "DATABASE_DSN"
)

type URLConfig struct {
	SA  string
	BU  string
	FSP string
	DSN string
}

var UC = URLConfig{}

func getSAfromEnv() {
	if UC.SA != "" {
		log.Println("config/getSAfromEnv: sa in flag is:", UC.SA)
		return
	}
	UC.SA = DefaultSA
	
	if s, ok := os.LookupEnv(sa); ok {
		log.Println("config/getSAfromEnv: sa in env is:", s)
		UC.SA = s
	}
	log.Println("config/getSAfromEnv: finished")
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

func getDSNfromEnv() {
	if UC.DSN != "" {
		return
	}
	UC.DSN = DefaultDSN
	if s, ok := os.LookupEnv(dsn); ok {
		UC.DSN = s
	}
}

func BuildConfig() {
	flag.Parse()
	log.Println("config/BuildConfig UC after flags", UC)

	getSAfromEnv()
	getBUfromEnv()
	getFSPfromEnv()
	getDSNfromEnv()
	log.Println("config/BuildConfig UC after env vars", UC)
}
