package core

import (
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/khelechy/argus/utils"
)

type Connection struct {
	Conn            net.Conn
	IsActive        bool
	IsAuthenticated bool
}

var connections = make(map[net.Conn]*Connection)

var (
	connLock sync.Mutex

	ConnUsername string
	ConnPassword string

	hasAuthData bool
)

func SetupTCP(host, port string) {
	// Listen for incoming connections

	if len(ConnUsername) > 0 && len(ConnPassword) > 0 {
		hasAuthData = true
	}

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
		if hasAuthData {
			go validateClientAuthentication(connection)
		} else {
			go handleClient(connection)
		}
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

func SendDataToClients(eventMsg string) {

	// Send data to the server
	data := []byte(eventMsg)

	if len(connections) == 0 {
		fmt.Println("No active connection")
		return
	}

	for connection := range connections {
		if !hasAuthData {
			go func(currentConnection net.Conn, currentData []byte) {

				sendDataToConnection(currentConnection, currentData)
			}(connection, data)
		} else {
			if connections[connection].IsAuthenticated {
				go func(currentConnection net.Conn, currentData []byte) {

					sendDataToConnection(currentConnection, currentData)
				}(connection, data)
			}else{
				fmt.Println("One or more connection(s) is not authenticated")
				data = []byte("Unauthorized connection \n")
				go func(currentConnection net.Conn, currentData []byte) {

					sendDataToConnection(currentConnection, currentData)
				}(connection, data)
			}
		}
	}
}

func sendDataToConnection(conn net.Conn, data []byte) {
	_, err := conn.Write(data)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}

func validateClientAuthentication(clientConn *Connection) {
	defer func() {
		clientConn.Conn.Close()
		connLock.Lock()
		delete(connections, clientConn.Conn)
		connLock.Unlock()
	}()

	buffer := make([]byte, 1024)

	n, err := clientConn.Conn.Read(buffer)

	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	connectionString := buffer[:n]

	username, password, err := utils.ExtractAuthData(string(connectionString))

	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	// Validiate Auth Credentials
	if !authenticate(username, password, clientConn.Conn) {
		fmt.Println("Authentication failed for: ", clientConn.Conn.RemoteAddr())
		return
	}

	fmt.Println("Authentication successful for:", clientConn.Conn.RemoteAddr())

	connLock.Lock()
	clientConn.IsAuthenticated = true
	connLock.Unlock()

	response := []byte("Authentication successful \n")

	_, err = clientConn.Conn.Write(response)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	handleClient(clientConn)

}

func handleClient(clientConn *Connection) {
	defer func() {
		clientConn.Conn.Close()
		connLock.Lock()
		delete(connections, clientConn.Conn)
		connLock.Unlock()
	}()

	fmt.Printf("%s connected successfully", clientConn.Conn.RemoteAddr())

	// Create a buffer to read data into
	buffer := make([]byte, 1024)

	for {
		// Read data from the client
		_, err := clientConn.Conn.Read(buffer)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

	}
}

func authenticate(username string, password string, conn net.Conn) bool {

	errorMessage := "Authentication failed, invalid username or password\n"

	if username != ConnUsername {

		response := []byte(errorMessage)

		_, _ = conn.Write(response)

		return false
	}

	if password != ConnPassword {

		response := []byte(errorMessage)

		_, _ = conn.Write(response)

		return false
	}

	return true
}
