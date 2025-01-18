package ipban

import (
	"fmt"
	"goipban/config"
	"log"
	"os/exec"
	"strings"
)

func BanIP(ip string) error {
	// Ban the IP using UFW
	log.Printf("Banning IP: %s\n", ip)
	err := runCommand(fmt.Sprintf(config.Ban.BanCommand, ip))
	if err != nil {
		log.Printf("Error banning IP %s: %v\n", ip, err)
	}
	return err
}

func UnbanIP(ip string) error {
	// Unan the IP using UFW
	log.Printf("Unbanning IP: %s\n", ip)
	err := runCommand(fmt.Sprintf(config.Ban.UnbanCommand, ip))
	if err != nil {
		log.Printf("Error unbanning IP %s: %v\n", ip, err)
	}
	return err
}

// runCommand executes shell commands
func runCommand(cmd string) error {
	log.Printf("Executing command: %s\n", cmd)
	parts := strings.Fields(cmd)
	head := parts[0]
	args := parts[1:]

	if config.Testing.RunStubCommands {
		return nil // don't ban, just print
	}

	out, err := exec.Command(head, args...).CombinedOutput()
	if err != nil {
		log.Printf("Command output: %s\n", string(out))
		return err
	}
	return nil
}
