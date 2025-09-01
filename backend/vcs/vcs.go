package vcs

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

// Ensure note directory and .vcs exist
func ensureDirs(noteID string) (string, string, error) {
	noteDir := filepath.Join(RootDir, noteID)
	vcsDir := filepath.Join(noteDir, ".vcs")

	if err := os.MkdirAll(vcsDir, 0755); err != nil {
		return "", "", err
	}
	return noteDir, vcsDir, nil
}

func SaveSnapshot(noteID, content string) (string, error) {
    noteDir := filepath.Join(RootDir, noteID)
    vcsDir := filepath.Join(noteDir, ".vcs")
    if err := os.MkdirAll(vcsDir, 0755); err != nil {
        return "", err
    }

    // main note file
    filePath := filepath.Join(noteDir, noteID+".txt")
    if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
        return "", err
    }

    hash := fmt.Sprintf("%x", sha256.Sum256([]byte(content+time.Now().String())))
    snap := Snaps{
        Files:    map[string][]byte{noteID + ".txt": []byte(content)},
        FileList: []string{noteID + ".txt"},
    }

    gobPath := filepath.Join(vcsDir, hash+".gob")
    f, err := os.Create(gobPath)
    if err != nil {
        return "", err
    }
    defer f.Close()
    if err := gob.NewEncoder(f).Encode(snap); err != nil {
        return "", err
    }

    return hash, nil
}

func GetLatestFileContent(noteID string) (string, error) {
    filePath := filepath.Join(RootDir, noteID, noteID+".txt")
    data, err := os.ReadFile(filePath)
    if err != nil {
        return "", fmt.Errorf("note not found")
    }
    return string(data), nil
}

func ListSnapshots(noteID string) ([]string, error) {
    vcsDir := filepath.Join(RootDir, noteID, ".vcs")
    files, err := filepath.Glob(filepath.Join(vcsDir, "*.gob"))
    if err != nil {
        return nil, err
    }

    sort.Slice(files, func(i, j int) bool {
        fi, _ := os.Stat(files[i])
        fj, _ := os.Stat(files[j])
        return fi.ModTime().Before(fj.ModTime())
    })

    var snaps []string
    for _, f := range files {
        snaps = append(snaps, filepath.Base(f))
    }
    return snaps, nil
}

func GetSnapshotContent(noteID, snapshotFile string) (string, error) {
    f, err := os.Open(filepath.Join(RootDir, noteID, ".vcs", snapshotFile))
    if err != nil {
        return "", err
    }
    defer f.Close()

    var snap Snaps
    if err := gob.NewDecoder(f).Decode(&snap); err != nil {
        return "", err
    }

    data, ok := snap.Files[noteID+".txt"]
    if !ok {
        return "", fmt.Errorf("file not found in snapshot")
    }
    return string(data), nil
}

func RestoreSnapshot(noteID, snapshotFile string) error {
    f, err := os.Open(filepath.Join(RootDir, noteID, ".vcs", snapshotFile))
    if err != nil {
        return err
    }
    defer f.Close()

    var snap Snaps
    if err := gob.NewDecoder(f).Decode(&snap); err != nil {
        return err
    }

    for fname, content := range snap.Files {
        path := filepath.Join(RootDir, noteID, fname)
        if err := os.WriteFile(path, content, 0644); err != nil {
            return err
        }
    }
    return nil
}


// DeleteNote removes the note directory and all snapshots
func DeleteNote(noteID string) error {
	noteDir := filepath.Join(RootDir, noteID)
	if _, err := os.Stat(noteDir); os.IsNotExist(err) {
		return fmt.Errorf("note '%s' does not exist", noteID)
	}

	if err := os.RemoveAll(noteDir); err != nil {
		return fmt.Errorf("failed to delete note '%s': %v", noteID, err)
	}

	return nil
}

// ListAllNotes returns all note IDs (directories) in storage
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
