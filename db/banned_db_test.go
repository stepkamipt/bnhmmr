package db

import (
	"goipban/config"
	"goipban/models"
	"testing"
	"time"
)

func TestConnect(t *testing.T) {
	config.BannedDB.FilePath = "testing_db.sqlite"

	db, err := ConnectToBannedDB()
	if err != nil {
		t.Errorf("err connect database %v", err)
	}
	defer db.deleteTestDB()
}

func TestInsertBannedIPs(t *testing.T) {
	config.BannedDB.FilePath = "testing_db.sqlite"

	db, err := ConnectToBannedDB()
	if err != nil {
		t.Errorf("err connect database %v", err)
	}
	defer db.deleteTestDB()

	bannedIP := models.BannedIPEntry{
		IP:         "1.2.3.4",
		BannedFrom: time.Now(),
	}
	err = db.InsertBannedIP(bannedIP)
	if err != nil {
		t.Errorf("err inserting database %v", err)
		return
	}

	bannedIPs, err := db.GetIPsBannedBefore(time.Now())
	if err != nil {
		t.Errorf("err get values from database %v", err)
		return
	}
	if len(bannedIPs) != 1 {
		t.Errorf("get errors from banned before returns %d items, expected 1", len(bannedIPs))
	}

	err = db.RemoveBannedIP(bannedIP.IP)
	if err != nil {
		t.Errorf("err remove banned IP from database %v", err)
		return
	}
	bannedIPs, err = db.GetIPsBannedBefore(time.Now())
	if err != nil {
		t.Errorf("err get values from database %v", err)
		return
	}
	if len(bannedIPs) != 0 {
		t.Errorf("get errors from banned before returns %d items, expected 0", len(bannedIPs))
	}
}
