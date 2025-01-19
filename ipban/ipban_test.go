package ipban

import "testing"

func TestBanIP(t *testing.T) {
	ipBanner := CreateIPBanner(false)
	if err := ipBanner.BanIP("11.22.33.44"); err != nil {
		t.Errorf("Banning error: %v", err)
	} else {
		t.Logf("Banned successfully!")
	}
}

func TestUnbanIP(t *testing.T) {
	ipBanner := CreateIPBanner(false)
	if err := ipBanner.UnbanIP("11.22.33.44"); err != nil {
		t.Errorf("Unbanning error: %v", err)
	} else {
		t.Logf("Unbanned successfully!")
	}
}
