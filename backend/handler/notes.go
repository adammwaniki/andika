//backend/handler/notes.go
package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/adammwaniki/andika/backend/utils"
	"github.com/adammwaniki/andika/backend/vcs"
	"github.com/gofrs/uuid/v5"
)

type response struct {
	Message string `json:"message,omitempty"`
	ID      string `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Content string `json:"content,omitempty"`
	Hash    string `json:"hash,omitempty"`
	Error   string `json:"error,omitempty"`
}

type Note struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Content string `json:"content,omitempty"`
}

// Helper: generate UUID or empty string
func generateUUID() string {
	id, err := uuid.NewV4()
	if err != nil {
		return ""
	}
	return id.String()
}

// POST /notes
func CreateNoteHandler(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name    string `json:"name"`
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, response{Error: "invalid JSON body"})
		return
	}

	noteID := generateUUID()
	if noteID == "" {
		utils.WriteJSON(w, http.StatusInternalServerError, response{Error: "failed to generate note ID"})
		return
	}

	hash, err := vcs.SaveSnapshot(noteID, body.Content)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, response{Error: err.Error()})
		return
	}

	// Save mapping noteID -> note Title
	err = vcs.SaveNoteMeta(noteID, body.Name)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, response{Error: err.Error()})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, response{
		Message: "note created", ID: noteID, Name: body.Name, Hash: hash,
	})
}

// GET /notes/{id}
func ViewNoteHandler(w http.ResponseWriter, r *http.Request) {
	noteID := strings.TrimPrefix(r.URL.Path, "/notes/")
	meta, err := vcs.GetNoteMeta(noteID)
	if err != nil {
		utils.WriteJSON(w, http.StatusNotFound, response{Error: "note not found"})
		return
	}

	content, err := vcs.GetLatestFileContent(noteID)
	if err != nil {
		utils.WriteJSON(w, http.StatusNotFound, response{Error: err.Error()})
		return
	}

	utils.WriteJSON(w, http.StatusOK, response{
		ID: noteID,
		Name: meta.Title,
		Content: content,
	})
}

// PUT /notes/{id}
func EditNoteHandler(w http.ResponseWriter, r *http.Request) {
	noteID := strings.TrimPrefix(r.URL.Path, "/notes/")
	var body struct {
		Mode    string `json:"mode"`    // append | overwrite
		Content string `json:"content"` // new content
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, response{Error: "invalid JSON body"})
		return
	}

	current := ""
	if body.Mode == "append" || body.Mode == "edit" {
		existing, err := vcs.GetLatestFileContent(noteID)
		if err == nil {
			current = existing
		}
	}

	finalContent := body.Content
	if body.Mode == "append" || body.Mode == "edit" {
		finalContent = current + body.Content
	}

	hash, err := vcs.SaveSnapshot(noteID, finalContent)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, response{Error: err.Error()})
		return
	}

	meta, _ := vcs.GetNoteMeta(noteID)
	utils.WriteJSON(w, http.StatusOK, response{
		Message: "note updated",
		ID: noteID,
		Name: meta.Title,
		Hash: hash,
	})
}

// GET /notes
func ListNotesHandler(w http.ResponseWriter, r *http.Request) {
	notes, err := vcs.ListAllNotesMeta()
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, response{Error: err.Error()})
		return
	}

	// Convert NoteMeta list into response-friendly struct
	var resp []Note
	for _, n := range notes {
		resp = append(resp, Note{
			ID: n.ID,
			Name: n.Title,
		})
	}

	utils.WriteJSON(w, http.StatusOK, resp)
}

// DELETE /notes/{id}
func DeleteNoteHandler(w http.ResponseWriter, r *http.Request) {
	noteID := strings.TrimPrefix(r.URL.Path, "/notes/")
	if noteID == "" {
		utils.WriteJSON(w, http.StatusBadRequest, response{Error: "note ID is required"})
		return
	}

	// Check if note exists
	meta, err := vcs.GetNoteMeta(noteID)
	if err != nil {
		utils.WriteJSON(w, http.StatusNotFound, response{Error: "note not found"})
		return
	}

	err = vcs.DeleteNote(noteID)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, response{Error: err.Error()})
		return
	}

	utils.WriteJSON(w, http.StatusOK, response{
		Message: "note deleted successfully",
		ID: noteID,
		Name: meta.Title,
	})
}

