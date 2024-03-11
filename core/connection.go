package core

import (
	"fmt"
	"net"
	"sync"
	"log"
)

type Connection struct {
	Conn     net.Conn
	IsActive bool
}

var connections = make(map[net.Conn]*Connection)

var (
	connLock sync.Mutex
)

func SetupTCP(host, port string) {
	// Listen for incoming connections

	addr := fmt.Sprintf("%s:%s", host, port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal("Error:", err)
		return
	}

	defer listener.Close()

	fmt.Printf("Speak lord, your server is listening on port %s\n", port)

	go HandleBroadcast()

	for {
		// Accept incoming connections
		newConn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		connection := &Connection{
			Conn:     newConn,
			IsActive: true,
		}

		connLock.Lock()
		connections[newConn] = connection
		connLock.Unlock()

		// Handle client connection in a goroutine
		go handleClient(connection)
	}
}

func SendDataToClients(eventMsg string) {

	// Send data to the server
	data := []byte(eventMsg)

	if len(connections) == 0 {
		fmt.Println("No active connection")
		return
	}

	for connection := range connections{
		go func(currentConnection net.Conn){

			_, err := currentConnection.Write(data)

			if err != nil {
				fmt.Println("Error:", err)
				return
			}
		}(connection)
	}
	
}

func handleClient(clientConn *Connection) {
	defer func (){
		clientConn.Conn.Close()
		connLock.Lock()
		delete(connections, clientConn.Conn)
		connLock.Unlock()
	}()

	// Create a buffer to read data into
	buffer := make([]byte, 1024)

	for {
		// Read data from the client
		n, err := clientConn.Conn.Read(buffer)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		// Process and use the data (here, we'll just print it)
		fmt.Printf("Received: %s\n", buffer[:n])
	}
}

func HandleBroadcast() {

	fmt.Println("Listening for filewatcher message signals to broadcast")
	for {
		eventMsg := <-messageChan
		if len(eventMsg) > 0 {
			SendDataToClients(eventMsg)
		}
	}
}
