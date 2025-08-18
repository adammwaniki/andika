package vcs

import (
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
)

// Directory to hold the snapshot data including a subdirectory to hold the file contents. Python equivalent of a dict with a single key "files" whose value is another empty dict
// Adding gob tags for addition of fields later
type Snapshot struct {
    Files map[string][]byte	`gob:"files"`
	FileList []string		`gob:"file_list"`
}

// Directory creation function to hold our snapshots
// It needs execute permissions to open and/or traverse it
func initVCS() {
	if err := os.Mkdir("vcs_storage", 0750); err != nil && !os.IsExist(err) { // mode 0750 is rwxr-x--- for permissions: owner full, group read+execute, others none
		log.Fatal(err)
	}
}

// Snapshot creation function
func snapshot(directory string) (string, error) {
	// Create a SHA256 hasher
	snapshotHash := sha256.New()

	// Snapshot data
	snapshotData := Snapshot{
		Files: make(map[string][]byte),
	}	// Alternatively we could have used a map of maps e.g., snapshotData := map[string]map[string]any{"files": {},}
	
	// Now we walk through the directory capturing the directory tree and files
	err := filepath.WalkDir(directory, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip vcs_storage files
		if d.IsDir() && d.Name() == "vcs_storage" {
			return filepath.SkipDir
		}

		// Skip directories
		if d.IsDir() {
			return nil
		}

		// Read file contents
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		// Update and Store snapshot hash with file path + content
		// Using relative path for hashing and storing so that the paths aren't OS specific i.e., \ on Windows and / on Linux
		relPath, err := filepath.Rel(directory, path)
		if err != nil {
			return err
		}

		snapshotHash.Write([]byte(relPath))
		snapshotHash.Write(content)

		// Store file content in snapshot
		snapshotData.Files[relPath] = content
		return nil
	})

	if err != nil {
		return "", err
	}

	// Now we finalize the hash calculation for the snapshot
	hashDigest := fmt.Sprintf("%x", snapshotHash.Sum(nil))
	
	// Then we sort and save the list of files that we collected in the snapshot for later reference. 
	for filePath := range snapshotData.Files {
		snapshotData.FileList = append(snapshotData.FileList, filePath)
	}
	sort.Strings(snapshotData.FileList)

	// Ensure vcs_storage exists
	if err := os.MkdirAll("vcs_storage", 0750); err != nil {
		return "", err
	}

	// We now serialise and save the snapshot's data to a file named after the snapshot's hash using gob
	file, err := os.Create(fmt.Sprintf("vcs_storage/%s", hashDigest))
	if err != nil {
		return "", err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	if err := encoder.Encode(snapshotData); err != nil {
		return "", err
	}

	fmt.Println("Snapshot saved:", hashDigest)
	return hashDigest, nil
}

// loadSnapshot loads a snapshot's data from vcs_storage by its hash
func loadSnapshot(hash string) (*Snapshot, error) {
	// Check if snapshot exists
	snapshotPath := fmt.Sprintf("vcs_storage/%s", hash)
	if _, err := os.Stat(snapshotPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("snapshot %s does not exist", hash)
	}

	// Open the snapshot file
	file, err := os.Open(fmt.Sprintf("vcs_storage/%s", hash)) // os.Open() also handles non-existent files, I just wanted an explicit check
	if err != nil {
		return nil, fmt.Errorf("failed to open snapshot: %w", err)
	}
	defer file.Close()

	// Decode the gob data into a Snapshot struct
	var snapshotData Snapshot
	decoder := gob.NewDecoder(file)
	if err := decoder.Decode(&snapshotData); err != nil {
		return nil, fmt.Errorf("failed to decode snapshot: %w", err)
	}

	return &snapshotData, nil
}