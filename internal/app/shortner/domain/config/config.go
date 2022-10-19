package config

import (
	"flag"
	"os"
)

/*
func init() {
	// log.Println("Config/getConfigFromCLI, passed flag: a,b,f")
	sa = flag.String("a", "", "SERVER_ADDRESS")
	bu = flag.String("b", "", "BASE_URL")
	fsp = flag.String("f", "", "FILE_STORAGE_PATH")
}

// declare config vars
var (
	sa  *string
	bu  *string
	fsp *string
	//default values
	dsa  = ":8080"
	dbu  = "http://localhost:8080"
	dfsp = ""
)

// declare config struct
type URLConfig struct {
	SA  *string //SERVER_ADDRESS
	BU  *string //BASE_URL
	FSP *string //FILE_STORAGE_PATH
}


// new confit struct instance
var UC = URLConfig{}

// func getConfigFromCLI() {

	flag.Parse()
	// log.Println("Config/getConfigFromCLI: CLI flag declared and parsed - completed")
	// log.Println("Config/getConfigFromCLI: CLI flag:", *sa, *bu, *fsp)
}
*/

var (
	sa  string
	bu  string
	fsp string
	//default values
	dsa  = ":8080"
	dbu  = "http://localhost:8080"
	dfsp = ""
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
	if sa != "" {
		return
	}
	// if empty set default
	UC.SA = dsa
	if s, ok := os.LookupEnv("SERVER_ADDRESS"); ok {
		UC.SA = s
	}
}

func getBUfromEnv() {
	// get bu from env if empty
	if bu != "" {
		return
	}
	UC.BU = dbu
	if s, ok := os.LookupEnv("BASE_URL"); ok {
		UC.BU = s
	}
}

func getFSPfromEnv() {
	// get fsp from env if empty
	if fsp != "" {
		return
	}
	UC.FSP = dfsp
	if s, ok := os.LookupEnv("FILE_STORAGE_PATH"); ok {
		UC.FSP = s
	}
}

func BuildConfig() {
	// get config vars from CLI flags
	// getConfigFromCLI()

	flag.StringVar(&UC.SA, "a", dsa, "SERVER_ADDRESS")
	flag.StringVar(&UC.BU, "b", dbu, "BASE_URL")
	flag.StringVar(&UC.FSP, "f", dfsp, "FILE_STORAGE_PATH")

	flag.Parse()

	// get config from local var if was not set by flag
	getSAfromEnv()
	getBUfromEnv()
	getFSPfromEnv()
	// log.Println("Config/BuildConfig: final URL config object:", *UC.SA, *UC.BU, *UC.FSP)
}
