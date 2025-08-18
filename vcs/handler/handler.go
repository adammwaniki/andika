package handler

import (
	"bufio"
	"fmt"
	"os"

	"github.com/adammwaniki/andika/vcs/utils"
)

// Handles all commands
func HandleCommand(command string, args []string) {
	switch command {
	case "create":
		if len(args) < 1 {
			fmt.Println("Usage: create <file>")
			return
		}
		noteName := args[0]
		fmt.Println("Enter note content (type 'xsave' on a new line to save):")
		content := readMultiLineInput()
		hash, err := utils.SaveSnapshot(noteName, content)
		if err != nil {
			fmt.Println("Error creating note:", err)
		} else {
			fmt.Println("Note created with snapshot:", hash)
		}

	case "view":
		if len(args) < 1 {
			fmt.Println("Usage: view <file>")
			return
		}
		noteName := args[0]
		content, err := utils.GetLatestFileContent(noteName)
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("\n---", noteName, "---")
			fmt.Println(content)
		}

	case "edit", "append", "overwrite":
		if len(args) < 1 {
			fmt.Printf("Usage: %s <file>\n", command)
			return
		}
		noteName := args[0]
		var current string
		var err error

		if command == "append" || command == "edit" {
			current, err = utils.GetLatestFileContent(noteName)
			if err != nil {
				current = "" // new file or no content
			}
		}

		if command == "overwrite" || command == "edit" {
			fmt.Println("Existing content (edit mode):")
			fmt.Println(current)
		}

		fmt.Println("Type note content (type 'xsave' on a new line to save):")
		newContent := readMultiLineInput()
		finalContent := newContent
		if command == "append" || command == "edit" {
			finalContent = current + newContent
		}

		hash, err := utils.SaveSnapshot(noteName, finalContent)
		if err != nil {
			fmt.Println("Error saving snapshot:", err)
		} else {
			fmt.Println("Snapshot saved:", hash)
		}

	case "list":
		notes, err := utils.ListAllNotes()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		if len(notes) == 0 {
			fmt.Println("No notes found.")
			return
		}
		fmt.Println("Available notes:")
		for _, n := range notes {
			fmt.Println(" -", n)
		}

	case "list_snaps":
		if len(args) < 1 {
			fmt.Println("Usage: list_snaps <file>")
			return
		}
		noteName := args[0]
		snaps, err := utils.ListSnapshots(noteName)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		if len(snaps) == 0 {
			fmt.Println("No snapshots for", noteName)
			return
		}
		fmt.Println("Snapshots for", noteName, ":")
		for i, s := range snaps {
			fmt.Printf("%d. %s\n", i+1, s)
		}

	case "snap":
		if len(args) < 2 {
			fmt.Println("Usage: snap <view|restore> <snapshot hash>")
			return
		}
		action := args[0]
		snapHash := args[1]

		noteName, err := utils.FindNoteBySnapshot(snapHash)
		if err != nil {
			fmt.Println("Snapshot not found:", err)
			return
		}

		switch action {
		case "view":
			content, err := utils.GetSnapshotContent(noteName, snapHash+".gob")
			if err != nil {
				fmt.Println("Error viewing snapshot:", err)
			} else {
				fmt.Println("\n--- Snapshot content ---")
				fmt.Println(content)
			}
		case "restore":
			if err := utils.RestoreSnapshot(noteName, snapHash+".gob"); err != nil {
				fmt.Println("Error restoring snapshot:", err)
			} else {
				fmt.Println("Restored", noteName, "to snapshot", snapHash)
			}
		default:
			fmt.Println("Invalid snap action:", action)
		}

	case "help":
		PrintHelp()

	default:
		fmt.Println("Unknown command:", command)
	}
}


// Reads multiline input until 'xsave'
func readMultiLineInput() string {
	scanner := bufio.NewScanner(os.Stdin)
	var input string
	for scanner.Scan() {
		line := scanner.Text()
		if line == "xsave" {
			break
		}
		input += line + "\n"
	}
	return input
}

// Prints help message
func PrintHelp() {
	fmt.Println("  VCS Notes Manager - Commands:")
	fmt.Println("  Usage: <command> [arguments]")
	fmt.Println()
	fmt.Println("  create <file>             Create a new note file and open for editing")
	fmt.Println("  view <file>               View the latest content of a note")
	fmt.Println("  edit <file>               Edit the note (overwrite or append)")
	fmt.Println("  append <file>             Append content to an existing note")
	fmt.Println("  overwrite <file>          Overwrite an existing note")
	fmt.Println("  list                      List all note files")
	fmt.Println("  list_snaps <file>         List all snapshots for a note (chronological order)")
	fmt.Println("  snap view <snapshot>      View the content of a specific snapshot")
	fmt.Println("  snap restore <snapshot>   Restore a note to a specific snapshot")
	fmt.Println("  help                      Show this help message")
	fmt.Println("  exit / quit               Exit the session")
	fmt.Println()
	fmt.Println("Typing content for notes ends when you enter 'xsave' on a new line.")
}