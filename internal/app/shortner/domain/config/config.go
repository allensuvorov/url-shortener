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

	UC.SA = flag.String("a", ":8080", "SERVER_ADDRESS")
	bu = flag.String("b", "http://localhost:8080", "BASE_URL")
	flag.Parse()
	log.Println("Config/BuildConfig: CLI flag declared and parsed completed")

	if len(*UC.SA) == 0 {
		UC.SA = getSA()
	}
	UC.BU = getBU()
}

func getSA() *string {
	log.Println("Config/GetSA: started")
	log.Println("Config/GetSA: Port from Flag is", *sa)
	// get server address from local env if not in cli flags

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
