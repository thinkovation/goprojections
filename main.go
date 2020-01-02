package main

import (
	"fmt"
	"log"
)

func main() {

	server := NewServer("8888")

	fmt.Printf("Starting Server on port %v\n", server.Addr)
	log.Printf("Starting Server on port %v\n", server.Addr)
	go func() {
		// This starts the HTTP server
		err := server.ListenAndServe()

		if err != nil {
			fmt.Println("Error: ", err)
			log.Fatalln("WOAH! Cannot Start Server, exiting %v", err)
		}
	}()

	//wait shutdown
	server.WaitShutdown()

	log.Printf("Service Exiting")
}
