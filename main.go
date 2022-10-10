package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	fmt.Println("From main function..")
	PrintWriter()
	PrintLoadManager()
	PrintReadReplica()

	log.Println("Starting listening on port 1234")
	port := ":1234"

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("Listening on %s", port)
	srv := NewGRPCServer()

	if err := srv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
