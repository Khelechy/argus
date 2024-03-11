package main

import (
	"encoding/json"
	"fmt"
	"os"
	"io"

	"github.com/khelechy/argus/core"
	"github.com/khelechy/argus/models"
)

type Config struct {
	Server struct {
		Host string `json:"host"`
		Port string    `json:"port"`
	} `json:"server"`
	Watch []models.WatchStructure `json:"watch"`
}

func main() {

	//strip config and validate
	configFile, err := os.Open("C:/Users/PFY-102.PFY-102/source/repos/Mine/argus/config.json")
	if err != nil {
		fmt.Println(err)
	}

	defer  configFile.Close()
	byteResult, _ := io.ReadAll(configFile)

	var config Config
	json.Unmarshal([]byte(byteResult), &config)

	go func(watch []models.WatchStructure) {
		core.Watch(watch)
	}(config.Watch)

	core.SetupTCP(config.Server.Host, config.Server.Port)

}
