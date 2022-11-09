package storage

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/allensuvorov/urlshortner/internal/app/shortner/domain/entity"
	"github.com/allensuvorov/urlshortner/internal/app/shortner/domain/hashmap"
)

func write(ue entity.DTO, fsp string) error {
	log.Printf("Storage/File: saving to path - %s", fsp)

	// Create and open file
	file, err := os.OpenFile(fsp, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		log.Fatal(err)
	}
	// Close file at the end
	defer file.Close()

	// Write to file
	enc := json.NewEncoder(file) // will be encoding to file

	// h, u, id
	// let's assume we don't have to record attemps of shortening existing urls
	// type record struct {
	// 	url map[string]string
	// 	history map[string]map[string]bool
	// }

	err = enc.Encode(ue)
	// err := enc.Encode(map[string]string{h: u});
	if err != nil { // add map to buff
		fmt.Println(err)
		return nil
	}

	return nil
}

func restore(fsp string) inMemory {
	log.Println("File/restore: restoring data from file")
	um := make(hashmap.URLHashMap) // url map
	ca := make(hashmap.ClientActivity)

	// open file
	file, err := os.OpenFile(fsp, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		log.Fatal(err)
	}
	dec := json.NewDecoder(file)

	// Go over the data
	for dec.More() {
		t := entity.DTO{}
		if err := dec.Decode(&t); err != nil {
			log.Fatal(err)
		}
		log.Println("File/restore: restoring URL entry from file:", t)

		// push data to maps
		// for _, v := range t {
		um[t.Hash] = t.URL

		_, ok := ca[t.ClientID]
		if !ok {
			ca[t.ClientID] = make(map[string]bool)
		}
		ca[t.ClientID][t.Hash] = true
		// }
	}

	if err := file.Close(); err != nil {
		log.Fatal(err)
	}

	log.Println("File/restore: all restored data in map:", um, ca)
	return inMemory{um, ca}
}
