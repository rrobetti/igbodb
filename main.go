package main

import (
	"fmt"
	"igbodb/leaderwritter"
	"igbodb/loadmanager"
	"igbodb/readreplica"
	"igbodb/server"
	"log"
	"net"
)

func main() {
	fmt.Println("From main function..")
	leaderwritter.Print()
	loadmanager.Print()
	readreplica.Print()

	log.Println("Starting listening on port 8080")
	port := ":1234"

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("Listening on %s", port)
	srv := server.NewGRPCServer()

	if err := srv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
