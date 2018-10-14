package main

import (
	"./redis"
	"flag"
	"github.com/pkg/errors"
	"log"
)

// main
func main() {

	//redis.client()
	connect := flag.String("connect", "", "IP address of process to join. If empty, go into listen mode.")
	flag.Parse()

	// If the connect flag is set, go into client mode.
	if *connect != "" {
		err := redis.ClientFunc(*connect)
		if err != nil {
			log.Println("Error:", errors.WithStack(err))
		}
		log.Println("Client done.")
		return
	}

	// Else go into server mode.
	err := redis.ServerFunc()
	if err != nil {
		log.Println("Error:", errors.WithStack(err))
	}

	log.Println("Server done.")
}

// The Lshortfile flag includes file name and line number in log messages.
func init() {
	log.SetFlags(log.Lshortfile)
}
