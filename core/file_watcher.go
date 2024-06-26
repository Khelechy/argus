package core

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
	"strings"

	"path/filepath"

	"github.com/khelechy/argus/enums"
	"github.com/khelechy/argus/models"
	"github.com/khelechy/argus/utils"

	"github.com/fsnotify/fsnotify"
)

var messageChan chan string

func Watch(watchStructures []models.WatchStructure) {

	messageChan = make(chan string)

	// setup watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
		return
	}

	defer watcher.Close()

	done := make(chan bool)

	// use goroutine to start the watcher
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				newEvent := &models.Event{}
				newEvent.Timestamp = time.Now()
				newEvent.Name = event.Name

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
						log.Fatal(err)
						return
					}
					messageChan <- string(data)
				}(*newEvent)

			case err := <-watcher.Errors:
				log.Println("Error:", err)
			}
		}
	}()

	for _, watchStructure := range watchStructures {

		isWildcard, folderPath, extension := utils.TreatAsWildcard(watchStructure.Path)

		if isWildcard { // Handle as wildcard

			if watchStructure.WatchRecursively{
				if err := filepath.Walk(folderPath, func(path string, fi os.FileInfo, err error) error {

					if filepath.Ext(path) == extension {
						if err := watcher.Add(path); err != nil {
							log.Println("ERROR", err)
						} 
					}

					return nil
				}); err != nil {
					log.Println("ERROR", err)
				}
			}else{

				dir, err := os.Open(folderPath)
				if err != nil {
					log.Println("ERROR", err)
				}
				defer dir.Close()

				// Read the contents of the folder
				files, err := dir.Readdir(0)
				if err != nil {
					log.Println("ERROR", err)
				}

				// Iterate over the files in the folder
				for _, file := range files {

					// Check if the file has the supplied extension
					if strings.HasSuffix(file.Name(), extension) {

						newPath := fmt.Sprintf("%s/%s", folderPath, file.Name())
						err = watcher.Add(newPath)
						if err != nil {
							log.Println("ERROR", err)
						}
					}
				}

			}

		}else if !isWildcard && len(extension) > 0 { // Handle as a file
			
			err = watcher.Add(watchStructure.Path)
			if err != nil {
				log.Println("ERROR", err)
			}

		}else if !isWildcard && len(extension) == 0 { // Handle as a folder
			if watchStructure.WatchRecursively {
				if err := filepath.Walk(watchStructure.Path, func(path string, fi os.FileInfo, err error) error {
					if _, err := os.Stat(path); err != nil {
						return err
					}
	
					// since fsnotify can watch all the files in a directory, watchers only need
					// to be added to each nested directory
					if fi.Mode().IsDir() {
						return watcher.Add(path)
					}
	
					return nil
	
				}); err != nil {
					log.Println("ERROR", err)
				}
			}else{
				watcher.Add(watchStructure.Path)
			}
		}
		
	}

	<-done
}
