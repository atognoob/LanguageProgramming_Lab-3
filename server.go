package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

type Client struct {
	conn net.Conn
	name string
}

var (
	clients    = make(map[string]Client)
	clientLock sync.Mutex
)

func RunServer(host, port string) {
	address := fmt.Sprintf("%s:%s", host, port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println("Error when starting server:", err)
		return
	}
	defer listener.Close()
	fmt.Printf("Server is running at %s\n", address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	clientLock.Lock()
	if _, exists := clients[name]; exists {
		conn.Write([]byte("Name already exists. Please choose another name.\n"))
		clientLock.Unlock()
		return
	}
	clients[name] = Client{conn, name}
	clientLock.Unlock()

	fmt.Printf("%s is connected.\n", name)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			clientLock.Lock()
			delete(clients, name)
			clientLock.Unlock()
			fmt.Printf("%s disconnected\n", name)
			return
		}

		message = strings.TrimSpace(message)

		if strings.HasPrefix(message, "@") {
			sendToSpecificClient(name, message)
		} else {
			broadcastMessage(name, message)
		}
	}
}

func sendToSpecificClient(sender, message string) {
	parts := strings.SplitN(message, " ", 2)
	if len(parts) < 2 {
		return
	}

	target := parts[0][1:]
	content := parts[1]

	clientLock.Lock()
	targetClient, exists := clients[target]
	clientLock.Unlock()

	if exists {
		targetClient.conn.Write([]byte(fmt.Sprintf("[%s -> %s]: %s\n", sender, target, content)))
	} else {
		senderClient, _ := clients[sender]
		senderClient.conn.Write([]byte(fmt.Sprintf("User '%s' does not exist.\n", target)))
	}
}

func broadcastMessage(sender, message string) {
	clientLock.Lock()
	for _, client := range clients {
		if client.name != sender {
			client.conn.Write([]byte(fmt.Sprintf("[%s -> everyone]: %s\n", sender, message)))
		}
	}
	clientLock.Unlock()
}
