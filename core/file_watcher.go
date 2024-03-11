package core

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"path/filepath"

	"github.com/khelechy/argus/enums"
	"github.com/khelechy/argus/models"

	"github.com/fsnotify/fsnotify"
)

type Argus struct {
}

var messageChan chan string
var fileWatcher *fsnotify.Watcher

func Watch(watchStructures []models.WatchStructure) {

	messageChan = make(chan string)

	// setup watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	fileWatcher = watcher

	defer watcher.Close()
	defer fileWatcher.Close()

	done := make(chan bool)

	// use goroutine to start the watcher
	go func() {
		for {
			select {
			case event := <-fileWatcher.Events:
				newEvent := &models.Event{}
				newEvent.Timestamp = time.Now()
				newEvent.Name = event.Name
				eventData, err := json.Marshal(event)
				if err != nil {
					fmt.Println(err)
					return
				}
				newEvent.EventMetaData = string(eventData)

				if event.Op&fsnotify.Create == fsnotify.Create {
					newEvent.ActionDescription = fmt.Sprintf("File created: %s", event.Name)
					newEvent.Action = enums.Create
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					newEvent.ActionDescription = fmt.Sprintf("File modified: %s", event.Name)
					newEvent.Action = enums.Write
				}
				if event.Op&fsnotify.Remove == fsnotify.Remove {
					newEvent.ActionDescription = fmt.Sprintf("File removed: %s", event.Name)
					newEvent.Action = enums.Delete
				}
				if event.Op&fsnotify.Rename == fsnotify.Rename {
					newEvent.ActionDescription = fmt.Sprintf("File renamed: %s", event.Name)
					newEvent.Action = enums.Rename
				}
				if event.Op&fsnotify.Chmod == fsnotify.Chmod {
					newEvent.ActionDescription = fmt.Sprintf("File permissions modified: %s", event.Name)
					newEvent.Action = enums.Chmod
				}
				go func(transferredEvent models.Event) {
					data, err := json.Marshal(transferredEvent)
					if err != nil {
						fmt.Println(err)
						return
					}
					messageChan <- string(data)
				}(*newEvent)

			case err := <-fileWatcher.Errors:
				log.Println("Error:", err)
			}
		}
	}()

	for _, watchStructure := range watchStructures {
		if watchStructure.IsFolder && watchStructure.WatchRecursively { // Watch Recursively
			if err := filepath.Walk(watchStructure.Path, watchDirRecursively); err != nil {
				fmt.Println("ERROR", err)
			}
		} else {

			err = watcher.Add(watchStructure.Path)
			if err != nil {
				fmt.Println("ERROR", err)
			}
		}
	}


	<-done
}

func watchDirRecursively(path string, fi os.FileInfo, err error) error {

	// since fsnotify can watch all the files in a directory, watchers only need
	// to be added to each nested directory
	if fi.Mode().IsDir() {
		return fileWatcher.Add(path)
	}

	return nil
}
