package config

import (
	"flag"
	"os"
)

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

// new congit struct instance
var UC = URLConfig{}

func getConfigFromCLI() {

	flag.Parse()
	// log.Println("Config/getConfigFromCLI: CLI flag declared and parsed - completed")
	// log.Println("Config/getConfigFromCLI: CLI flag:", *sa, *bu, *fsp)
}

func getSAfromEnv() {

	// get sa from env if empty
	if *sa == "" {
		s, ok := os.LookupEnv("SERVER_ADDRESS")
		// if empty set default
		if !ok {
			// log.Printf("Config/getSAfromEnv: %s not set\n; passing default", "SERVER_ADDRESS")
			s = dsa
		}
		*sa = s
	}
	UC.SA = sa
}

func getBUfromEnv() {
	// get bu from env if empty
	if *bu == "" {
		s, ok := os.LookupEnv("BASE_URL")
		if !ok {
			// log.Printf("Config/getBUfromEnv: %s not set\n; passing default", "BASE_URL")
			s = dbu
		}
		*bu = s
	}
	UC.BU = bu
}

func getFSPfromEnv() {

	if *fsp == "" {
		s, ok := os.LookupEnv("FILE_STORAGE_PATH")

		if !ok {
			s = dfsp
		}
		*fsp = s
	}
	UC.FSP = fsp
}

func BuildConfig() {
	// get config vars from CLI flags
	getConfigFromCLI()

	// get config from local var if was not set by flag
	getSAfromEnv()
	getBUfromEnv()
	getFSPfromEnv()
	// log.Println("Config/BuildConfig: final URL config object:", *UC.SA, *UC.BU, *UC.FSP)
}
