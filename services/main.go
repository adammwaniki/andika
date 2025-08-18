package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/adammwaniki/andika/services/vcs"
)

func main() {
	vcs.InitVCS()
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\nMain Menu:")
		fmt.Println("1. Create new subdirectory")
		fmt.Println("2. Access existing subdirectory")
		fmt.Println("3. Exit")
		fmt.Print("> ")

		choice := readInput(reader)

		switch choice {
		case "1":
			createSubdirectory(reader)
		case "2":
			accessSubdirectory(reader)
		case "3":
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid choice")
		}
	}
}

func readInput(reader *bufio.Reader) string {
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func createSubdirectory(reader *bufio.Reader) {
	fmt.Print("Enter new subdirectory name: ")
	name := readInput(reader)
	path := filepath.Join("vcs_storage", name)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.Mkdir(path, 0750); err != nil {
			fmt.Println("Failed to create subdirectory:", err)
			return
		}
		fmt.Println("Subdirectory created:", name)
	} else {
		fmt.Println("Subdirectory already exists")
	}

	subdirectoryMenu(reader, path)
}

func accessSubdirectory(reader *bufio.Reader) {
	fmt.Println("Access subdirectory by:")
	fmt.Println("1. List all subdirectories")
	fmt.Println("2. Input subdirectory name")
	fmt.Println("3. Back")
	fmt.Print("> ")

	choice := readInput(reader)
	switch choice {
	case "1":
		subdirs, err := vcs.ListSubdirs()
		if err != nil {
			fmt.Println("Error listing subdirs:", err)
			return
		}
		fmt.Println("Subdirectories:")
		for _, d := range subdirs {
			fmt.Println("-", d)
		}
		fmt.Print("Enter subdirectory name to access: ")
		name := readInput(reader)
		path := filepath.Join("vcs_storage", name)
		subdirectoryMenu(reader, path)
	case "2":
		fmt.Print("Enter subdirectory name: ")
		name := readInput(reader)
		path := filepath.Join("vcs_storage", name)
		subdirectoryMenu(reader, path)
	case "3":
		return
	default:
		fmt.Println("Invalid choice")
	}
}

func subdirectoryMenu(reader *bufio.Reader, path string) {
	for {
		fmt.Printf("\nSubdirectory: %s\n", path)
		fmt.Println("1. Create snapshot")
		fmt.Println("2. List snapshots")
		fmt.Println("3. View snapshot by hash")
		fmt.Println("4. Revert to snapshot by hash")
		fmt.Println("5. Back")
		fmt.Println("6. Exit")
		fmt.Print("> ")

		choice := readInput(reader)
		switch choice {
		case "1":
			fmt.Print("Enter directory path to snapshot (will create if missing): ")
			dir := readInput(reader)
			absDir, _ := filepath.Abs(dir)
			if _, err := os.Stat(absDir); os.IsNotExist(err) {
				if err := os.MkdirAll(absDir, 0755); err != nil {
					fmt.Println("Failed to create directory:", err)
					continue
				}
			}
			hash, err := vcs.Snapshot(absDir)
			if err != nil {
				fmt.Println("Snapshot failed:", err)
			} else {
				fmt.Println("Snapshot created:", hash)
			}
		case "2":
			snapshots, err := vcs.ListSnapshots(path)
			if err != nil {
				fmt.Println("Failed to list snapshots:", err)
				continue
			}
			fmt.Println("Snapshots:")
			for _, s := range snapshots {
				fmt.Println("-", s)
			}
		case "3":
			fmt.Print("Enter snapshot hash to view: ")
			hash := readInput(reader)
			snap, err := vcs.LoadSnapshot(path, hash)
			if err != nil {
				fmt.Println("Failed to load snapshot:", err)
				continue
			}
			fmt.Println("Files in snapshot:")
			for _, f := range snap.FileList {
				fmt.Println("-", f)
			}
		case "4":
			fmt.Print("Enter snapshot hash to revert: ")
			hash := readInput(reader)
			if err := vcs.RevertToSnapshot(path, hash); err != nil {
				fmt.Println("Revert failed:", err)
			}
		case "5":
			return
		case "6":
			fmt.Println("Exiting...")
			os.Exit(0)
		default:
			fmt.Println("Invalid choice")
		}
	}
}
