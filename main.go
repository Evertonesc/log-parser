package main

import (
	"log"
)

func main() {
	err := parseLog("qgames.log")
	if err != nil {
		log.Fatal(err)
	}
}
