//backend/handler/snapshots.go
package handler

import (
	"net/http"
	"strings"

	"github.com/adammwaniki/andika/backend/utils"
	"github.com/adammwaniki/andika/backend/vcs"
)

// ListSnapshotsHandler handles GET /notes/{id}/snapshots
func ListSnapshotsHandler(w http.ResponseWriter, r *http.Request, params map[string]string) {
	noteID := params["id"]

	meta, err := vcs.GetNoteMeta(noteID)
	if err != nil {
		utils.WriteJSON(w, http.StatusNotFound, response{Error: "note not found"})
		return
	}

	snaps, err := vcs.ListSnapshots(noteID)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, response{Error: err.Error()})
		return
	}

	// Strip .gob extension for cleaner API response
	for i := range snaps {
		snaps[i] = strings.TrimSuffix(snaps[i], ".gob")
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"id":     noteID,
		"name":   meta.Title,
		"hashes": snaps,
	})
}

// ViewSnapshotHandler handles GET /notes/{id}/snapshots/{sid}
func ViewSnapshotHandler(w http.ResponseWriter, r *http.Request, params map[string]string) {
	noteID := params["id"]
	hash := params["sid"] // hash without .gob

	meta, err := vcs.GetNoteMeta(noteID)
	if err != nil {
		utils.WriteJSON(w, http.StatusNotFound, response{Error: "note not found"})
		return
	}

	content, err := vcs.GetSnapshotContent(noteID, hash+".gob")
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, response{Error: err.Error()})
		return
	}

	utils.WriteJSON(w, http.StatusOK, response{
		ID:      noteID,
		Name:    meta.Title,
		Content: content,
		Hash:    hash,
	})
}

// RestoreSnapshotHandler handles POST /notes/{id}/snapshots/{sid}/restore
func RestoreSnapshotHandler(w http.ResponseWriter, r *http.Request, params map[string]string) {
	noteID := params["id"]
	hash := params["sid"] // hash without .gob

	meta, err := vcs.GetNoteMeta(noteID)
	if err != nil {
		utils.WriteJSON(w, http.StatusNotFound, response{Error: "note not found"})
		return
	}

	if err := vcs.RestoreSnapshot(noteID, hash+".gob"); err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, response{Error: err.Error()})
		return
	}

	utils.WriteJSON(w, http.StatusOK, response{
		Message: "note restored",
		ID:      noteID,
		Name:    meta.Title,
		Hash:    hash,
	})
}
