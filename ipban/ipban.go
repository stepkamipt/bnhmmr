package ipban

import (
	"fmt"
	"goipban/config"
	"log"
	"os/exec"
	"strings"
)

var commands = struct {
	Ban   string
	Unban string
	Apply string
}{
	Ban:   "ufw insert 2 deny from %s",
	Unban: "ufw delete deny from %s",
	Apply: "ufw reload",
}

// Ban the IP using UFW
func BanIP(ip string) error {
	log.Printf("Banning IP: %s\n", ip)
	if err := runCommand(fmt.Sprintf(commands.Ban, ip)); err != nil {
		log.Printf("Error banning IP %s: %v\n", ip, err)
		return err
	}
	if err := runCommand(commands.Apply); err != nil {
		log.Printf("Error applying ban: %v\n", err)
		return err
	}
	return nil
}

// Unan the IP using UFW
func UnbanIP(ip string) error {
	log.Printf("Unbanning IP: %s\n", ip)
	if err := runCommand(fmt.Sprintf(commands.Unban, ip)); err != nil {
		log.Printf("Error unbanning IP %s: %v\n", ip, err)
		return err
	}
	if err := runCommand(commands.Apply); err != nil {
		log.Printf("Error applying unban: %v\n", err)
		return err
	}
	return nil
}

// runCommand executes shell commands
func runCommand(cmd string) error {
	log.Printf("Executing command: %s\n", cmd)
	parts := strings.Fields(cmd)
	head := parts[0]
	args := parts[1:]

	if config.DebugRunStubCommands {
		return nil // don't ban, just print
	}

	out, err := exec.Command(head, args...).CombinedOutput()
	if err != nil {
		log.Printf("Command output: %s\n", string(out))
		return err
	}
	return nil
}
