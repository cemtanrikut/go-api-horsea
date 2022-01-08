package main

import (
	"log"

	"github.com/cemtanrikut/go-api-horsea/router"
)

func main() {
	log.Println("Starting the application")
	router.MuxUserHandler()

}
