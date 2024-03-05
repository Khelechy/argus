package argus

import (
	"fmt"
	"log"
	"os"
	"net"
	"sync"
	"bufio"

	//"path/filepath"

	"github.com/fsnotify/fsnotify"
)

type Argus struct {
}

var (
	connection net.Conn
	connLock   sync.Mutex
)

var watcher *fsnotify.Watcher

func Watch(testfile string) {
	// setup watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	// use goroutine to start the watcher
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				// monitor only for write events
				if event.Op&fsnotify.Write == fsnotify.Write {
					fmt.Println("Modified file:", event.Name)
				}
			case err := <-watcher.Errors:
				log.Println("Error:", err)
			}
		}
	}()

	// provide the file name along with path to be watched
	err = watcher.Add(testfile)
	if err != nil {
		log.Fatal(err)
	}

	// if err := filepath.Walk("/Users/skdomino/Desktop/test", watchDirRecursively); err != nil {
	// 	fmt.Println("ERROR", err)
	// }
	<-done
}

func watchDirRecursively(path string, fi os.FileInfo, err error) error {

	// since fsnotify can watch all the files in a directory, watchers only need
	// to be added to each nested directory
	if fi.Mode().IsDir() {
		return watcher.Add(path)
	}

	return nil
}




func SetupTCP() {
	// Listen for incoming connections
    listener, err := net.Listen("tcp", "localhost:8080")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
	
    defer listener.Close()

    fmt.Println("Server is listening on port 8080")

    for {
        // Accept incoming connections
        conn, err := listener.Accept()
        if err != nil {
            fmt.Println("Error:", err)
            continue
        }

		connLock.Lock()
		connection = conn
		connLock.Unlock()

        // Handle client connection in a goroutine
        go handleClient(conn)
    }
}

func SendDataToClient(conn net.Conn){
	// Send data to the server
	data := []byte("Hello, Server!")

	_, err := conn.Write(data)

	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}

func handleClient(conn net.Conn) {
    defer conn.Close()

	for {
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			fmt.Println("Client says:", scanner.Text())
		}
	}

	// // Create a buffer to read data into
	// buffer := make([]byte, 1024)

	//  for {
	// 	// Read data from the client
	// 	n, err := conn.Read(buffer)
	// 	if err != nil {
	// 	fmt.Println("Error:", err)
	// 	return
	//   }
	
	//   // Process and use the data (here, we'll just print it)
	//   fmt.Printf("Received: %s\n", buffer[:n])
	// }
}

