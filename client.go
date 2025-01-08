package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func RunClient(host, port, username string) {
	address := fmt.Sprintf("%s:%s", host, port)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Printf("Unable to connect to server at %s. Error: %v\n", address, err)
		return
	}
	defer conn.Close()

	fmt.Println("Successfully connected to server.")
	conn.Write([]byte(username + "\n"))

	go receiveMessages(conn)

	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)

		if text == "" {
			continue
		}

		conn.Write([]byte(text + "\n"))
	}
}

func receiveMessages(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Disconnect from server.")
			return
		}
		fmt.Print(message)
	}
}
