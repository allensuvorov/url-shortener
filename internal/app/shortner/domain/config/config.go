package config

import (
	"flag"
	"log"
	"os"
)

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

func BuildConfig() {
	// get config vars from CLI flags
	sa = flag.String("a", "", "SERVER_ADDRESS")
	bu = flag.String("b", "", "BASE_URL")

	flag.Parse()
	log.Println("Config/BuildConfig: CLI flag declared and parsed completed")

	// get config from local var if was not set by flag
	if len(*sa) == 0 {
		sa = getSA()
	}
	UC.SA = sa
	UC.BU = getBU()
}

func getSA() *string {
	log.Println("Config/GetSA: started")
	log.Println("Config/GetSA: Port from Flag is", *sa)

	saStr, ok := os.LookupEnv("SERVER_ADDRESS")
	if !ok {
		log.Printf("%s not set\n; passing default", "SERVER_ADDRESS")
		saStr = ":8080"
	}
	*sa = saStr

	return sa
}

func getBU() *string {
	// get BU from flag

	if len(*bu) == 0 {

		// get BU from local env
		log.Println("Config/GetBU: about to take BASE_URL from local env")
		buStr, ok := os.LookupEnv("BASE_URL")
		if ok {
			log.Println("Config/GetBU: BASE_URL from local env is:", buStr)
		}
		if !ok {
			log.Printf("Config/GetBU: %s not set\n; passing default", "BASE_URL")
			buStr = "http://localhost:8080"
		}
		*bu = buStr
	}
	return bu
}
