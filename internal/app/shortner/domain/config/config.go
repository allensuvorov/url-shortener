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

	log.Println("Config/getConfigFromCLI, passed flag: f")
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
	SA  *string
	BU  *string
	FSP *string
}

// new congit struct instance
var UC = URLConfig{}

func getConfigFromCLI() {

	flag.Parse()
	log.Println("Config/getConfigFromCLI: CLI flag declared and parsed - completed")
	log.Println("Config/getConfigFromCLI: CLI flag:", *sa, *bu, *fsp)
}

func getSAfromEnv() {

	// get sa from env if empty
	if len(*sa) == 0 {
		s, ok := os.LookupEnv("SERVER_ADDRESS")
		// if empty set default
		if !ok {
			log.Printf("Config/getSAfromEnv: %s not set\n; passing default", "SERVER_ADDRESS")
			s = dsa
		}
		*sa = s
	}
	UC.SA = sa
}

func getBUfromEnv() {
	// get bu from env if empty
	if len(*bu) == 0 {
		s, ok := os.LookupEnv("BASE_URL")
		if !ok {
			log.Printf("Config/getBUfromEnv: %s not set\n; passing default", "BASE_URL")
			s = dbu
		}
		*bu = s
	}
	UC.BU = bu
}

func getFSPfromEnv() {

	// Set env for testing
	os.Setenv("FILE_STORAGE_PATH", "/Users/allen/go/src/yandex/projects/urlshortner/internal/app/shortner/storage/urls.txt")
	s, _ := os.LookupEnv("FILE_STORAGE_PATH")
	log.Print("Config/getSAfromEnv: set FSP env var for testing:", s)

	if len(*fsp) == 0 {
		s, ok := os.LookupEnv("FILE_STORAGE_PATH")
		log.Println("Config/getFSPfromEnv: fsp in env var is", s)

		if !ok {
			log.Printf("Config/GetFSP: %s not set\n; passing default", "FILE_STORAGE_PATH")
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
	log.Println("Config/BuildConfig: final URL config object:", *UC.SA, *UC.BU, *UC.FSP)
}
