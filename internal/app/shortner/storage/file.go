package storage

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func write(h, u, fsp string) error {
	log.Printf("Storage: saving to path - %s", fsp)

	// Create and open file
	fileName := "urls.txt"
	file, err := os.OpenFile(fsp+fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
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
