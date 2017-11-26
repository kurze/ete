package main

import (
	"log"
)

func main() {
	config, err := readConfig("conf.yml")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("conf: %+v", config)

	server := newServer(config)
	server.registerEndpoints()
	server.run()
}
