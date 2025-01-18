package models

import "time"

type XRayLogEntry struct {
	Time     time.Time `json:"time"`
	FromIP   string    `json:"from_ip"`
	FromPort string    `json:"from_port"`
	To       string    `json:"to"`
	Inbound  string    `json:"inbound"`
	Outbound string    `json:"outbound"`
}

type BannedIPEntry struct {
	IP         string    `json:"ip"`
	BannedFrom time.Time `json:"banned_from"`
}
