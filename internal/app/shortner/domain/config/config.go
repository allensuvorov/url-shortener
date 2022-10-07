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

// get config vars from CLI flags
// func init() {
// 	SA = flag.String("a", ":8080", "SERVER_ADDRESS")
// 	BU = flag.String("b", "http://localhost:8080", "BASE_URL")
// }

func BuildConfig() {
	// get config vars from CLI flags

	sa = flag.String("a", ":8080", "SERVER_ADDRESS")
	bu = flag.String("b", "http://localhost:8080", "BASE_URL")
	flag.Parse()
	log.Println("Config/BuildConfig: completed")

}

// new congit struct instance
var UC = URLConfig{
	SA: getSA(),
	BU: getBU(),
}

// declare config struct
type URLConfig struct {
	SA *string
	BU *string
	// FSP *string
}

func getSA() *string {
	log.Println("Config/GetSA: started")
	log.Println("Config/GetSA: Port from Flag is", *sa)
	// get server address from local env if not in cli flags
	if len(*sa) == 0 {
		saStr, ok := os.LookupEnv("SERVER_ADDRESS")
		if !ok {
			log.Printf("%s not set\n; passing default", "SERVER_ADDRESS")
			saStr = ":8080"
		}
		*sa = saStr
	}
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
