package banning_loop

import (
	"fmt"
	"goipban/config"
	"goipban/db"
	"goipban/ipban"
	"goipban/logparse"
	"goipban/models"
	"log"
	"time"
)

func StartBanningLoop() error {
	// connect to database
	bannedDB, err := db.ConnectToBannedDB()
	if err != nil {
		return fmt.Errorf("failed to connect banned db: %w", err)
	}
	defer bannedDB.Close()

	for {
		time.Sleep(config.Process.Autoupdate)

		if err = banningLoopIteration(*bannedDB); err != nil {
			log.Printf("banning failed: %s", err)
		}
	}

	return nil
}

func banningLoopIteration(bannedDB db.BannedDB) error {
	// parse logs
	logEntriesToBan, err := logparse.GetBlacklistedXRayLogEntries()
	if err != nil {
		return fmt.Errorf("can not parse logs %s", err)
	}

	// ban every non-banned IP in logs
	if err = banNonBannedIPs(bannedDB, logEntriesToBan); err != nil {
		return err
	}

	// unban IP's which unban time has come
	if err = unbanReleasingIPs(bannedDB); err != nil {
		return err
	}

	return nil
}

func banNonBannedIPs(bannedDB db.BannedDB, logEntriesToBan []models.XRayLogEntry) error {
	// ban every non-banned IP in logs
	var bannedIPCount int
	for i := range logEntriesToBan {
		// check if IP already banned
		isBannedIP, err := bannedDB.IsBannedIP(logEntriesToBan[i].FromIP)
		if err != nil {
			return fmt.Errorf("can not query database %s", err)
		}
		if isBannedIP {
			continue
		}

		// ban IP
		banningIP := models.BannedIPEntry{
			IP:         logEntriesToBan[i].FromIP,
			BannedFrom: logEntriesToBan[i].Time,
		}
		err = bannedDB.InsertBannedIP(banningIP)
		if err != nil {
			return fmt.Errorf("can not query database %s", err)
		}

		err = ipban.BanIP(banningIP.IP)
		if err != nil {
			log.Printf("can not ban ip %s", err)
		}
		bannedIPCount++
	}
	log.Printf("Banned %d IP", bannedIPCount)

	return nil
}

func unbanReleasingIPs(bannedDB db.BannedDB) error {
	// unban IP's which unban time has come
	ipsToUnban, err := bannedDB.GetIPsBannedBefore(time.Now().Add(-config.Ban.Duration))
	if err != nil {
		return fmt.Errorf("can not query database %s", err)
	}

	var unbannedIPCount int
	for i := range ipsToUnban {
		err = ipban.UnbanIP(ipsToUnban[i].IP)
		if err != nil {
			log.Printf("can not ban ip %s", err)
			continue
		}
		err = bannedDB.RemoveBannedIP(ipsToUnban[i].IP)
		unbannedIPCount++
		if err != nil {
			log.Println("can not access database %s", err)
			return err
		}
	}
	log.Printf("Unbanned %d IP", unbannedIPCount)

	return nil
}
