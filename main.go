package main

import (
	"github.com/gelleson/gomond/gomond/cmd"
	"log"
)

func main() {
	err := cmd.Execute()

	if err != nil {
		log.Fatal(err)
	}

}
