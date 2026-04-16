package main

import (
	"log"
	"user-profiles/internal/app/user-profiles"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
