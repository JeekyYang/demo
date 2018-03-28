package main

import (
	"fmt"
	"log"
	"demo/server"
)

func main() {
	s, err := server.NewServer()
	if err != nil {
		log.Fatalf("failed to create server, %+v", err)
		return
	}

	s.Start()

	fmt.Println("Demo project start!")
}
