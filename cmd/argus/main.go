package main

import (
	//"fmt"

	"github.com/khelechy/argus"
)

func main() {
	go func (){
		argus.Watch("C:/Users/PFY-102.PFY-102/source/repos/Mine/argus/file.txt")
	}()
	
	argus.SetupTCP()

}

