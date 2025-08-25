package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/adammwaniki/andika/vcs/utils"
)

type response struct {
	Message  string `json:"message,omitempty"`
	NoteName string `json:"noteName,omitempty"`
	Content  string `json:"content,omitempty"`
	Hash     string `json:"hash,omitempty"`
	Error    string `json:"error,omitempty"`
}

// Helper: write JSON response
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// POST /notes/{noteName}
func CreateNoteHandler(w http.ResponseWriter, r *http.Request) {
	noteName := strings.TrimPrefix(r.URL.Path, "/notes/")
	var body struct {
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, response{Error: "invalid JSON body"})
		return
	}

	hash, err := utils.SaveSnapshot(noteName, body.Content)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, response{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, response{Message: "note created", Hash: hash, NoteName: noteName})
}

// GET /notes/{noteName}
func ViewNoteHandler(w http.ResponseWriter, r *http.Request) {
	noteName := strings.TrimPrefix(r.URL.Path, "/notes/")
	content, err := utils.GetLatestFileContent(noteName)
	if err != nil {
		writeJSON(w, http.StatusNotFound, response{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, response{NoteName: noteName, Content: content})
}

// PUT /notes/{noteName}
func EditNoteHandler(w http.ResponseWriter, r *http.Request) {
	noteName := strings.TrimPrefix(r.URL.Path, "/notes/")
	var body struct {
		Mode    string `json:"mode"`    // append | overwrite
		Content string `json:"content"` // new content
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, response{Error: "invalid JSON body"})
		return
	}

	current := ""
	if body.Mode == "append" || body.Mode == "edit" {
		existing, err := utils.GetLatestFileContent(noteName)
		if err == nil {
			current = existing
		}
	}
	finalContent := body.Content
	if body.Mode == "append" || body.Mode == "edit" {
		finalContent = current + body.Content
	}

	hash, err := utils.SaveSnapshot(noteName, finalContent)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, response{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, response{Message: "note updated", Hash: hash, NoteName: noteName})
}

// GET /notes
func ListNotesHandler(w http.ResponseWriter, r *http.Request) {
	notes, err := utils.ListAllNotes()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, response{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, notes)
}

// GET /notes/{noteName}/snapshots
func ListSnapshotsHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		writeJSON(w, http.StatusBadRequest, response{Error: "invalid path"})
		return
	}
	noteName := parts[2]
	snaps, err := utils.ListSnapshots(noteName)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, response{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, snaps)
}

// GET /snapshots/{hash}
func ViewSnapshotHandler(w http.ResponseWriter, r *http.Request) {
	hash := strings.TrimPrefix(r.URL.Path, "/snapshots/")
	noteName, err := utils.FindNoteBySnapshot(hash)
	if err != nil {
		writeJSON(w, http.StatusNotFound, response{Error: err.Error()})
		return
	}
	content, err := utils.GetSnapshotContent(noteName, hash+".gob")
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, response{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, response{NoteName: noteName, Content: content, Hash: hash})
}

// POST /snapshots/{hash}/restore
func RestoreSnapshotHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/snapshots/"), "/")
	if len(parts) < 2 || parts[1] != "restore" {
		writeJSON(w, http.StatusBadRequest, response{Error: "invalid restore path"})
		return
	}
	hash := parts[0]

	noteName, err := utils.FindNoteBySnapshot(hash)
	if err != nil {
		writeJSON(w, http.StatusNotFound, response{Error: err.Error()})
		return
	}
	if err := utils.RestoreSnapshot(noteName, hash+".gob"); err != nil {
		writeJSON(w, http.StatusInternalServerError, response{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, response{Message: "note restored", NoteName: noteName, Hash: hash})
}

// GET /help
func HelpHandler(w http.ResponseWriter, r *http.Request) {
	commands := []string{
		"POST   /notes/{noteName}            -> create a note",
		"GET    /notes/{noteName}            -> view latest note content",
		"PUT    /notes/{noteName}            -> edit/append/overwrite note",
		"GET    /notes                       -> list all notes",
		"GET    /notes/{noteName}/snapshots  -> list all snapshots of a note",
		"GET    /snapshots/{hash}            -> view snapshot content",
		"POST   /snapshots/{hash}/restore    -> restore note to snapshot",
	}
	writeJSON(w, http.StatusOK, map[string]any{"commands": commands})
}
