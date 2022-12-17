package storage

import (
	"encoding/json"
	"log"
	"os"

	"github.com/allensuvorov/urlshortner/internal/app/shortner/domain/entity"
	"github.com/allensuvorov/urlshortner/internal/app/shortner/domain/hashmap"
)

func write(ue entity.URLEntity, fsp string) error {
	log.Printf("Storage/File: saving to path - %s", fsp)

	file, err := os.OpenFile(fsp, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	enc := json.NewEncoder(file)

	err = enc.Encode(ue)
	if err != nil {
		log.Println(err)
		return nil
	}

	return nil
}

func restore(fsp string) inMemory {
	log.Println("File/restore: restoring data from file")
	um := make(hashmap.URLHashMap)
	ca := make(hashmap.ClientActivity)
	dd := make(hashmap.Deleted)

	file, err := os.OpenFile(fsp, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		log.Fatal(err)
	}
	dec := json.NewDecoder(file)

	for dec.More() {
		t := entity.URLEntity{}
		if err := dec.Decode(&t); err != nil {
			log.Fatal(err)
		}
		log.Println("File/restore: restoring URL entry from file:", t)

		um[t.Hash] = t.URL
		dd[t.Hash] = t.Deleted
		_, ok := ca[t.ClientID]
		if !ok {
			ca[t.ClientID] = make(map[string]bool)
		}
		ca[t.ClientID][t.Hash] = true
	}

	if err := file.Close(); err != nil {
		log.Fatal(err)
	}

	log.Println("File/restore: all restored data in map:", um, ca)
	return inMemory{um, ca, dd}
}
