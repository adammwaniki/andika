package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/adammwaniki/andika/vcs/utils"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\nChoose option: ")
		fmt.Println("1. Create new note")
		fmt.Println("2. View note")
		fmt.Println("3. Append to note")
		fmt.Println("4. Overwrite note")
		fmt.Println("5. List all notes")
		fmt.Println("6. List snapshots for a note")
		fmt.Println("7. Restore note to snapshot")
		fmt.Println("8. View contents of a snapshot")
		fmt.Println("9. Exit")
		fmt.Print("Enter choice: ")

		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1": // Create new note
			fmt.Print("Enter note name: ")
			noteName := readLine(reader)
			content := readMultiLineInput()
			hash, err := utils.SaveSnapshot(noteName, content)
			if err != nil {
				fmt.Println("Error creating note:", err)
			} else {
				fmt.Println("Note created with snapshot:", hash)
			}

		case "2": // View note
			fmt.Print("Enter note name: ")
			noteName := readLine(reader)
			content, err := utils.GetLatestFileContent(noteName)
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println("\n---", noteName, "---")
				fmt.Println(content)
			}

		case "3": // Append
			fmt.Print("Enter note name: ")
			noteName := readLine(reader)
			latestContent, err := utils.GetLatestFileContent(noteName)
			if err != nil {
				latestContent = "" // no previous content
			}
			fmt.Println("Type notes to append (type 'vcs_save' to save and exit):")
			newContent := readMultiLineInput()
			finalContent := latestContent + newContent
			hash, err := utils.SaveSnapshot(noteName, finalContent)
			if err != nil {
				fmt.Println("Error saving snapshot:", err)
			} else {
				fmt.Println("Snapshot saved:", hash)
			}

		case "4": // Overwrite
			fmt.Print("Enter note name: ")
			noteName := readLine(reader)
			fmt.Println("Type new content (type 'vcs_save' to save and exit):")
			newContent := readMultiLineInput()
			hash, err := utils.SaveSnapshot(noteName, newContent)
			if err != nil {
				fmt.Println("Error saving snapshot:", err)
			} else {
				fmt.Println("Snapshot saved:", hash)
			}

		case "5": // List all notes
			notes, err := utils.ListAllNotes()
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}
			if len(notes) == 0 {
				fmt.Println("No notes found.")
			} else {
				fmt.Println("Available notes:")
				for _, n := range notes {
					fmt.Println(" -", n)
				}
			}

		case "6": // List snapshots oldest to newest
			fmt.Print("Enter note name: ")
			noteName := readLine(reader)
			snapshots, err := utils.ListSnapshots(noteName)
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}
			if len(snapshots) == 0 {
				fmt.Println("No snapshots found for", noteName)
			} else {
				fmt.Println("Snapshots for", noteName, ":")
				for i, s := range snapshots {
					fmt.Printf(" %d. %s\n", i+1, s)
				}
			}

		case "7": // Restore
			fmt.Print("Enter note name: ")
			noteName := readLine(reader)
			snapshots, err := utils.ListSnapshots(noteName)
			if err != nil || len(snapshots) == 0 {
				fmt.Println("No snapshots found for note:", noteName)
				continue
			}

			fmt.Println("Select snapshot number to restore:")
			for i, s := range snapshots {
				fmt.Printf(" %d. %s\n", i+1, s)
			}

			fmt.Print("Enter number: ")
			numInput := readLine(reader)
			num, err := strconv.Atoi(numInput)
			if err != nil || num < 1 || num > len(snapshots) {
				fmt.Println("Invalid selection")
				continue
			}

			selectedSnapshot := snapshots[num-1]
			if err := utils.RestoreSnapshot(noteName, selectedSnapshot); err != nil {
				fmt.Println("Error restoring snapshot:", err)
			} else {
				fmt.Println("Restored", noteName, "to snapshot", selectedSnapshot)
			}
		
		case "8": // View snapshot content
			fmt.Print("Enter note name: ")
			noteName := readLine(reader)

			// List snapshots
			snapshots, err := utils.ListSnapshots(noteName)
			if err != nil || len(snapshots) == 0 {
				fmt.Println("No snapshots found for", noteName)
				break
			}

			// Print numbered list
			fmt.Println("Snapshots:")
			for i, s := range snapshots {
				fmt.Printf("%d. %s\n", i+1, s)
			}

			fmt.Print("Enter snapshot number to view: ")
			numStr := readLine(reader)
			var num int
			fmt.Sscan(numStr, &num)
			if num < 1 || num > len(snapshots) {
				fmt.Println("Invalid number")
				break
			}

			content, err := utils.GetSnapshotContent(noteName, snapshots[num-1])
			if err != nil {
				fmt.Println("Error reading snapshot:", err)
			} else {
				fmt.Println("\n--- Snapshot content ---")
				fmt.Println(content)
			}

		case "9": // Exit
			fmt.Println("Goodbye!")
			return

		default:
			fmt.Println("Invalid option")
		}
	}
}

// Helper to read a single line from stdin
func readLine(reader *bufio.Reader) string {
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

// Helper to read multiline input until vcs_save
func readMultiLineInput() string {
	scanner := bufio.NewScanner(os.Stdin)
	var input string
	for scanner.Scan() {
		line := scanner.Text()
		if line == "vcs_save" {
			break
		}
		input += line + "\n"
	}
	return input
}
