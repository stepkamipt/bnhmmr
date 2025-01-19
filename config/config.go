package config

import "time"

const (
	// ban duration
	BanDuration = 1 * time.Minute
	// auto-update ban list interval
	UpdateInterval = 5 * time.Second
	// XRay access.log file path
	XRayLogsFilePath = "/var/log/xray/access.log"
	// XRay outbound tag for IPs be banned
	XRayBlacklistOutbound = "blacklist"
	// banned database file
	BannedDatabaseFile = "data/banned.sqlitedb"

	// don't really ban anybody, just log info
	DebugRunStubCommands = false
)

/*// Process params
var Process = struct {
	UpdateInterval time.Duration
}{
	UpdateInterval: 5 * time.Second,
}

// XRay logs params
var XRayLogs = struct {
	FilePath string
	BlacklistOutbound string
}{
	FilePath: "/var/log/xray/access.log",
	BlacklistOutbound: "blacklist",
}

// Banned database params
var BannedDB = struct {
	FilePath string
}{
	FilePath: "banned.sqlitedb",
}

// Testing params
var Testing = struct {
	RunStubCommands bool
}{
	RunStubCommands: false, // don't really ban anybody, just log info
}*/
