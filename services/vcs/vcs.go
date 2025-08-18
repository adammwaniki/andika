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
type Snaps struct {
    Files map[string][]byte	`gob:"files"`
	FileList []string		`gob:"file_list"`
}

// Directory creation function to hold our snapshots
// It needs execute permissions to open and/or traverse it
func InitVCS() {
	if err := os.Mkdir("vcs_storage", 0750); err != nil && !os.IsExist(err) { // mode 0750 is rwxr-x--- for permissions: owner full, group read+execute, others none
		log.Fatal(err)
	}
}

// Snapshot creation function
func Snapshot(directory string) (string, error) {
	// Create a SHA256 hasher
	snapshotHash := sha256.New()

	// Snapshot data
	snapshotData := Snaps{
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
func LoadSnapshot(hash string) (*Snaps, error) {
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
	var snapshotData Snaps
	decoder := gob.NewDecoder(file)
	if err := decoder.Decode(&snapshotData); err != nil {
		return nil, fmt.Errorf("failed to decode snapshot: %w", err)
	}

	return &snapshotData, nil
}

// RevertToSnapshot restores the filesystem to the state of the given snapshot.
func RevertToSnapshot(hash string) error {
	// Load snapshot
	snapshotData, err := LoadSnapshot(hash)
	if err != nil {
		return err
	}

	// Restore files from snapshot deterministically
	var restoreFiles []string
	for f := range snapshotData.Files {
		restoreFiles = append(restoreFiles, f)
	}
	sort.Strings(restoreFiles)

	for _, filePath := range restoreFiles {
		content := snapshotData.Files[filePath]
		if err := os.MkdirAll(filepath.Dir(filePath), 0750); err != nil {
			return fmt.Errorf("failed to create directory for %s: %w", filePath, err)
		}
		if err := os.WriteFile(filePath, content, 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", filePath, err)
		}
	}

	// Collect current files on disk
	currentFiles := make(map[string]struct{})
	err = filepath.WalkDir(".", func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and vcs_storage to avoid touching snapshot data
		if d.IsDir() {
			if d.Name() == "vcs_storage" {
				return filepath.SkipDir
			}
			return nil
		}

		// Using relative paths for consistency
		relPath, err := filepath.Rel(".", path)
		if err != nil {
			return err
		}
		currentFiles[relPath] = struct{}{}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to walk current directory: %w", err)
	}

	// Delete files not in snapshot
	snapshotFiles := make(map[string]struct{})
	for _, f := range snapshotData.FileList {
		snapshotFiles[f] = struct{}{}
	}

	var toDelete []string
	for filePath := range currentFiles {
		if _, exists := snapshotFiles[filePath]; !exists {
			toDelete = append(toDelete, filePath)
		}
	}
	sort.Strings(toDelete)
	for _, filePath := range toDelete {
		if err := os.Remove(filePath); err != nil {
			fmt.Printf("Failed to remove %s: %v\n", filePath, err)
			continue
		}
		fmt.Println("Removed", filePath)
	}

	fmt.Println("Reverted to snapshot", hash)
	return nil
}