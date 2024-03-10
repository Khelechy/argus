package argus

import (
	"fmt"
	"log"
	"os"

	//"path/filepath"

	"github.com/fsnotify/fsnotify"
)

type Argus struct {
}
var messageChan chan string
var watcher *fsnotify.Watcher

func Watch(testfile string) {

	messageChan = make(chan string)
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
					go func(){
						messageChan <-"00"
					}()
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
