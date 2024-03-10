package argus

import (
	"fmt"
	"net"

	//"strings"
	"sync"
)

type Connection struct {
	Conn     net.Conn
	IsActive bool
}

var connections = make(map[net.Conn]*Connection)

var (
	connLock sync.Mutex
)

func SetupTCP() {
	// Listen for incoming connections
	listener, err := net.Listen("tcp", "localhost:1337")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	defer listener.Close()

	fmt.Println("Server is listening on port 8080")

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

func SendDataToClients() {

	fmt.Println("sending data to client")

	// Send data to the server
	data := []byte("Hello, Client!")

	if len(connections) == 0 {
		fmt.Println("No active connection")
		return
	}

	for connection := range connections{
		go func(){

			_, err := connection.Write(data)

			if err != nil {
				fmt.Println("Error:", err)
				return
			}
		}()
	}
	
}

func handleClient(clientConn *Connection) {
	defer func (){
		clientConn.Conn.Close()
		connLock.Lock()
		delete(connections, clientConn.Conn)
		fmt.Println("Connection removed")
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
		messageType := <-messageChan
		if messageType == "00" {
			fmt.Println("Message sent to client")
			SendDataToClients()
		}
	}
}
