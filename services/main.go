package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/adammwaniki/andika/services/vcs"
)

func main() {
	// Initialize VCS storage
	vcs.InitVCS()
	fmt.Println("Initialized vcs_storage")

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\nChoose an action:")
		fmt.Println("1. Create snapshot")
		fmt.Println("2. List snapshots")
		fmt.Println("3. Revert to snapshot")
		fmt.Println("4. Exit")
		fmt.Print("> ")

		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			fmt.Print("Enter directory to snapshot: ")
			dir, _ := reader.ReadString('\n')
			dir = strings.TrimSpace(dir)

			// Convert to absolute path
			absDir, err := filepath.Abs(dir)
			if err != nil {
				log.Println("Failed to get absolute path:", err)
				continue
			}

			// Create the directory if it doesn't exist
			if _, err := os.Stat(absDir); os.IsNotExist(err) {
				fmt.Println("Directory does not exist, creating it...")
				if err := os.MkdirAll(absDir, 0755); err != nil {
					log.Println("Failed to create directory:", err)
					continue
				}
			}

			hash, err := vcs.Snapshot(absDir)
			if err != nil {
				log.Println("Failed to create snapshot:", err)
				continue
			}
			fmt.Println("Snapshot created with hash:", hash)

		case "2":
			fmt.Println("Available snapshots in vcs_storage:")
			files, err := os.ReadDir("vcs_storage")
			if err != nil {
				log.Println("Failed to list snapshots:", err)
				continue
			}
			for _, file := range files {
				fmt.Println("-", file.Name())
			}

		case "3":
			fmt.Print("Enter snapshot hash to revert to: ")
			hash, _ := reader.ReadString('\n')
			hash = strings.TrimSpace(hash)

			if err := vcs.RevertToSnapshot(hash); err != nil {
				log.Println("Failed to revert:", err)
				continue
			}

		case "4":
			fmt.Println("Exiting...")
			return

		default:
			fmt.Println("Invalid choice, try again")
		}
	}
}
