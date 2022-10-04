package storage

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/allensuvorov/urlshortner/internal/app/shortner/domain/hashmap"
)

func write(h, u, fsp string) error {
	log.Printf("Storage: saving to path - %s", fsp)

	// Create and open file
	file, err := os.OpenFile(fsp, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		log.Fatal(err)
	}
	// Close file at the end
	defer file.Close()

	// Write to file
	enc := json.NewEncoder(file)                                // will be encoding to file
	if err := enc.Encode(map[string]string{h: u}); err != nil { // add map to buff
		fmt.Println(err)
		return nil
	}

	return nil
}

func restore(fsp string) hashmap.URLHashMap {
	log.Println("Storage: restoring data from file")
	um := make(hashmap.URLHashMap) // url map

	// open file
	file, err := os.OpenFile(fsp, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		log.Fatal(err)
	}
	dec := json.NewDecoder(file)

	// Go over the data
	for dec.More() {
		t := map[string]string{}
		if err := dec.Decode(&t); err != nil {
			log.Fatal(err)
		}
		log.Println("t =", t)

		// push data to url map
		for k, v := range t {
			um[k] = v
		}
	}
	log.Println(um)
	return um
}
