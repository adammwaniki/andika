//backend/handler/help.go
package handler

import (
	"net/http"

	"github.com/adammwaniki/andika/backend/utils"
)

// HelpHandler returns a list of available API endpoints
func HelpHandler(w http.ResponseWriter, r *http.Request) {
	commands := []string{
		"POST   /notes                              -> create a note (body: {name, content})",
		"GET    /notes                              -> list all notes",
		"GET    /notes/{id}                         -> view latest content of a note",
		"PUT    /notes/{id}                         -> edit mode append/overwrite note (body: {mode, content})",
		"DELETE /notes/{id}                         -> delete a note",
		"GET    /notes/{id}/snapshots               -> list all snapshots of a note",
		"GET    /notes/{id}/snapshots/{sid}         -> view a specific snapshot",
		"POST   /notes/{id}/snapshots/{sid}/restore -> restore note to a snapshot",
		"GET    /help                               -> display this help message",
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"message":  "Andika API v1 Help Menu",
		"commands": commands,
	})
}
