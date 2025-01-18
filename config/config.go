package config

import "time"

// Ban params
var Ban = struct {
	Duration     time.Duration
	BanCommand   string
	UnbanCommand string
}{
	Duration:     5 * time.Minute,
	BanCommand:   "ufw insert 2 deny from %s",
	UnbanCommand: "ufw delete deny from %s",
}

// Process params
var Process = struct {
	UpdateInterval time.Duration
}{
	UpdateInterval: 15 * time.Second,
}

// XRay logs params
var XRayLogs = struct {
	FilePath            string
	BlacklistLineSample string
	BlacklistOutbound   string
}{
	FilePath:            "/var/log/xray/access.log",
	BlacklistLineSample: "2025/01/13 12:34:56 from 12.34.56.78:9999 accepted tcp:https://ya.ru:443 [inbound -> blacklist] email: user_mail",
	BlacklistOutbound:   "blacklist",
}

// Banned database params
var BannedDB = struct {
	FilePath      string
	Table         string
	IPCol         string
	BannedTillCol string
}{
	FilePath:      "banned.sqlitedb",
	Table:         "banned_users",
	IPCol:         "ip",
	BannedTillCol: "banned_till",
}

// Testing params
var Testing = struct {
	RunStubCommands bool
}{
	RunStubCommands: false, // don't really ban anybody, just log info
}
