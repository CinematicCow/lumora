package main

import (
	"fmt"
	"log"

	"github.com/CinematicCow/lumora/internal/core"
)

func main() {
	fmt.Println("Hello from lumora!")

	const key = "userID"
	const value = "123"

	db, err := core.Open("./test-data")
	if err != nil {
		log.Print(err)
	}

	if err := db.Put(key, []byte(value)); err != nil {
		log.Print(err)
	}

	val, err := db.Get(key)
	if err != nil {
		log.Print(err)
	}

	fmt.Printf("Key: %s | Value: %s\nFetched Value: %s", key, value, val)

}
