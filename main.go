package main

import (
	"goipban/banning_loop"
	"log"
)

func main() {
	if err := banning_loop.StartBanningLoop(); err != nil {
		log.Printf("banning loop error %s", err)
		return
	}
}
