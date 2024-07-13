package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/khelechy/argus/core"
	"github.com/khelechy/argus/models"
)

type Config struct {
	Server struct {
		Host     string `json:"host"`
		Port     string `json:"port"`
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"server"`
	Watch []models.WatchStructure `json:"watch"`
}

var configFileName = flag.String("config", "config.json", "Location of the config file")

func main() {

	fmt.Println("  _    _    _    _   _  ")
	fmt.Println(" / \\  / \\  / \\  / \\ / \\ ")
	fmt.Println("( A )( R )( G )( U )( S )")
	fmt.Println(" \\_/  \\_/  \\_/  \\_/ \\_/ ")

	flag.Parse()

	//strip config and validate
	configFile, err := os.Open(*configFileName)
	if err != nil {
		log.Fatal(err)
		return
	}

	defer configFile.Close()
	byteResult, err := io.ReadAll(configFile)

	if err != nil {
		log.Fatalln(err)
		return
	}

	var config Config
	err = json.Unmarshal([]byte(byteResult), &config)

	if err != nil {
		log.Fatalln(err)
		return
	}

	if len(config.Server.Host) == 0 || len(config.Server.Port) == 0 {
		err = errors.New("empty server host or port")
		log.Fatalln(err)
		return
	}

	if len(config.Watch) == 0 {
		err = errors.New("no configured file or folder to watch")
		log.Fatalln(err)
		return
	}

	core.ConnUsername = config.Server.Username
	core.ConnPassword = config.Server.Password

	go func(watch []models.WatchStructure) {
		core.Watch(watch)
	}(config.Watch)

	core.SetupTCP(config.Server.Host, config.Server.Port)

}
