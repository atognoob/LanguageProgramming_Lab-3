package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	mode := flag.String("mode", "", "Select mode: server or client")
	host := flag.String("host", "127.0.0.1", "Server IP address")
	port := flag.String("port", "8080", "Connection port")
	username := flag.String("username", "", "Username (client only)")

	flag.Parse()

	switch *mode {
	case "server":
		RunServer(*host, *port)
	case "client":
		if *username == "" {
			fmt.Println("Please provide a username using the -username parameter.")
			os.Exit(1)
		}
		RunClient(*host, *port, *username)
	default:
		fmt.Println("Please select mode with -mode=server or -mode=client.")
		os.Exit(1)
	}
}
