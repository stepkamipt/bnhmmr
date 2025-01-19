package main

import (
	"goipban/banning_loop"
	"goipban/config"
	"log"
)

func main() {
	log.Printf("Tool to ban XRay users running.")
	log.Printf("Update Ban list interval: %s", config.UpdateInterval)
	log.Printf("Ban time: %s", config.BanDuration)
	log.Printf("Banning xray outbound: %s", config.XRayBlacklistOutbound)

	if err := banning_loop.StartBanningLoop(); err != nil {
		log.Printf("banning loop error %s", err)
		return
	}
}
