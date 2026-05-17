package main

import (
	"log"
	"user-profiles/internal/app"
)

func main() {
	err := app.Run()
	if err != nil {
		log.Fatalln(err)
	}
}
