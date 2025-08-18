package utils

import (
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

const RootDir = "vcs_storage"

// Snapshot structure
type Snaps struct {
	Files    map[string][]byte
	FileList []string
}

// Ensure directories exist
func ensureDirs(noteName string) (string, string, error) {
	noteDir := filepath.Join(RootDir, noteName)
	vcsDir := filepath.Join(noteDir, ".vcs")

	if err := os.MkdirAll(vcsDir, 0755); err != nil {
		return "", "", err
	}
	return noteDir, vcsDir, nil
}

// SaveSnapshot stores file contents into .vcs
func SaveSnapshot(noteName, content string) (string, error) {
	noteDir, vcsDir, err := ensureDirs(noteName)
	if err != nil {
		return "", err
	}

	// save/update main note file
	filePath := filepath.Join(noteDir, noteName+".txt")
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return "", err
	}

	// hash
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(content+time.Now().String())))

	// snapshot struct
	snap := Snaps{
		Files:    map[string][]byte{noteName + ".txt": []byte(content)},
		FileList: []string{noteName + ".txt"},
	}

	// save gob
	f, err := os.Create(filepath.Join(vcsDir, hash+".gob"))
	if err != nil {
		return "", err
	}
	defer f.Close()

	enc := gob.NewEncoder(f)
	if err := enc.Encode(snap); err != nil {
		return "", err
	}

	return hash, nil
}

// GetLatestFileContent returns latest plain text // not necessarily the latest snapshot?
func GetLatestFileContent(noteName string) (string, error) {
	noteDir, _, err := ensureDirs(noteName)
	if err != nil {
		return "", err
	}

	filePath := filepath.Join(noteDir, noteName+".txt")
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("no snapshots found for note: %s", noteName)
	}
	return string(data), nil
}


// ListSnapshots for a note chronologically
func ListSnapshots(noteName string) ([]string, error) {
	_, vcsDir, err := ensureDirs(noteName)
	if err != nil {
		return nil, err
	}
	files, err := filepath.Glob(filepath.Join(vcsDir, "*.gob"))
	if err != nil {
		return nil, err
	}

	// sort by file ModTime for chronological order
	sort.Slice(files, func(i, j int) bool {
		fi, _ := os.Stat(files[i])
		fj, _ := os.Stat(files[j])
		return fi.ModTime().Before(fj.ModTime())
	})

	var snapshots []string
	for _, f := range files {
		snapshots = append(snapshots, filepath.Base(f)) // only filename
	}
	return snapshots, nil
}



// ListAllNotes in vcs_storage
func ListAllNotes() ([]string, error) {
	entries, err := os.ReadDir(RootDir)
	if err != nil {
		return nil, err
	}
	var notes []string
	for _, entry := range entries {
		if entry.IsDir() {
			notes = append(notes, entry.Name())
		}
	}
	sort.Strings(notes)
	return notes, nil
}

// RestoreSnapshot restores a note to a previous snapshot
func RestoreSnapshot(noteName, snapshotFile string) error {
	noteDir, vcsDir, err := ensureDirs(noteName)
	if err != nil {
		return err
	}

	snapshotPath := filepath.Join(vcsDir, snapshotFile)
	f, err := os.Open(snapshotPath)
	if err != nil {
		return err
	}
	defer f.Close()

	var snap Snaps
	if err := gob.NewDecoder(f).Decode(&snap); err != nil {
		return err
	}

	for fname, content := range snap.Files {
		path := filepath.Join(noteDir, fname)
		if err := os.WriteFile(path, content, 0644); err != nil {
			return err
		}
	}

	return nil
}

// GetSnapshotContent returns the content of a specific snapshot without restoring it
func GetSnapshotContent(noteName, snapshotFile string) (string, error) {
	_, vcsDir, err := ensureDirs(noteName)
	if err != nil {
		return "", err
	}

	snapshotPath := filepath.Join(vcsDir, snapshotFile)
	f, err := os.Open(snapshotPath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	var snap Snaps
	if err := gob.NewDecoder(f).Decode(&snap); err != nil {
		return "", err
	}

	data, ok := snap.Files[noteName+".txt"]
	if !ok {
		return "", fmt.Errorf("file not found in snapshot")
	}
	return string(data), nil
}

// FindNoteBySnapshot finds which note contains a snapshot by hash
func FindNoteBySnapshot(snapHash string) (string, error) {
	notes, err := ListAllNotes()
	if err != nil {
		return "", err
	}

	for _, note := range notes {
		_, vcsDir, _ := ensureDirs(note)
		if _, err := os.Stat(filepath.Join(vcsDir, snapHash+".gob")); err == nil {
			return note, nil
		}
	}
	return "", fmt.Errorf("snapshot not found")
}

