//backend/utils/utils.go
package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int , data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)

}

func ReadJSON(r *http.Request, data any) error {
	if r.Body == nil {
		return fmt.Errorf("missing request body")
	} 
	return json.NewDecoder(r.Body).Decode(data)
}

// Provide standard output for http error messages
func WriteError(w http.ResponseWriter, status int, errorMessage error) {
	WriteJSON(w, status, map[string]string{"error": errorMessage.Error()})
}