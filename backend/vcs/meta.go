//backend/vcs/meta.go
package vcs

import (
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// NoteMeta stores metadata for a note
type NoteMeta struct {
	ID        string
	Title     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// SaveNoteMeta creates or updates metadata for a note
func SaveNoteMeta(noteID, title string) error {
	noteDir, _, err := ensureDirs(noteID)
	if err != nil {
		return err
	}

	metaPath := filepath.Join(noteDir, "meta.gob")

	var meta NoteMeta
	if _, err := os.Stat(metaPath); err == nil {
		// Load existing metadata
		f, err := os.Open(metaPath)
		if err == nil {
			defer f.Close()
			_ = gob.NewDecoder(f).Decode(&meta)
		}
	}

	// Update fields
	if meta.ID == "" {
		meta.ID = noteID
		meta.CreatedAt = time.Now()
	}
	meta.Title = title
	meta.UpdatedAt = time.Now()

	// Save metadata
	f, err := os.Create(metaPath)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := gob.NewEncoder(f).Encode(meta); err != nil {
		return err
	}

	return nil
}

// GetNoteMeta retrieves metadata for a specific note
func GetNoteMeta(noteID string) (*NoteMeta, error) {
	noteDir, _, err := ensureDirs(noteID)
	if err != nil {
		return nil, err
	}

	metaPath := filepath.Join(noteDir, "meta.gob")
	f, err := os.Open(metaPath)
	if err != nil {
		return nil, fmt.Errorf("no metadata found for note '%s'", noteID)
	}
	defer f.Close()

	var meta NoteMeta
	if err := gob.NewDecoder(f).Decode(&meta); err != nil {
		return nil, err
	}

	return &meta, nil
}

// ListAllNotesMeta retrieves metadata for all notes
func ListAllNotesMeta() ([]NoteMeta, error) {
	notes, err := ListAllNotes()
	if err != nil {
		return nil, err
	}

	var allMeta []NoteMeta
	for _, noteID := range notes {
		meta, err := GetNoteMeta(noteID)
		if err == nil {
			allMeta = append(allMeta, *meta)
		}
	}

	// Sort by CreatedAt ascending
	sort.Slice(allMeta, func(i, j int) bool {
		return allMeta[i].CreatedAt.Before(allMeta[j].CreatedAt)
	})

	return allMeta, nil
}
