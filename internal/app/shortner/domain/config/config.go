package config

import (
	"flag"
	"log"
	"os"
)

func init() {
	log.Println("Config/getConfigFromCLI, passed flag: a")
	sa = flag.String("a", "", "SERVER_ADDRESS")

	log.Println("Config/getConfigFromCLI, passed flag: b")
	bu = flag.String("b", "", "BASE_URL")
}

// declare config vars
var (
	sa *string
	bu *string
	//default values
	dsa = ":8080"
	dbu = "http://localhost:8080"
)

// declare config struct
type URLConfig struct {
	SA *string
	BU *string
	// FSP *string
}

// new congit struct instance
var UC = URLConfig{}

func getConfigFromCLI() {

	flag.Parse()
	log.Println("Config/BuildConfig: CLI flag declared and parsed - completed")
}

func getSAfromEnv() {
	// get sa from env if empty
	if len(*sa) == 0 {
		s, ok := os.LookupEnv("SERVER_ADDRESS")
		// if empty set default
		if !ok {
			log.Printf("%s not set\n; passing default", "SERVER_ADDRESS")
			s = dsa
		}
		*sa = s
	}
	UC.SA = sa
}

func getBUfromEng() {
	// get bu from env if empty
	if len(*bu) == 0 {

		// get BU from local env
		log.Println("Config/GetBU: about to take BASE_URL from local env")
		s, ok := os.LookupEnv("BASE_URL")
		if ok {
			log.Println("Config/GetBU: BASE_URL from local env is:", s)
		}
		if !ok {
			log.Printf("Config/GetBU: %s not set\n; passing default", "BASE_URL")
			s = dbu
		}
		*bu = s
	}
	UC.BU = bu
}

func BuildConfig() {
	// get config vars from CLI flags
	getConfigFromCLI()

	// get config from local var if was not set by flag
	getSAfromEnv()
	getBUfromEng()

}
