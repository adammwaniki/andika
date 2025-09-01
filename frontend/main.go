package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/adammwaniki/andika/frontend/views"
)

// Note struct to match the backend response (now includes content)
type Note struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Content string `json:"content,omitempty"`
}

// fetchNotesFromAPI calls the backend API to get all notes (with content)
func fetchNotesFromAPI() ([]views.Note, error) {
	resp, err := http.Get("http://localhost:8160/api/v1/notes")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch notes: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	var apiNotes []Note
	if err := json.Unmarshal(body, &apiNotes); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	// Convert API notes to view notes (with content)
	var viewNotes []views.Note
	for _, note := range apiNotes {
		viewNotes = append(viewNotes, views.Note{
			ID:      note.ID,
			Title:   note.Name,
			Content: note.Content,
		})
	}

	return viewNotes, nil
}

// createNoteViaAPI creates a new note via the backend API and returns created note (including content)
func createNoteViaAPI(title, content string) (Note, error) {
	// Prepare the request payload
	noteData := map[string]string{
		"name":    title,
		"content": content,
	}

	jsonData, err := json.Marshal(noteData)
	if err != nil {
		return Note{}, fmt.Errorf("failed to marshal note data: %v", err)
	}

	log.Printf("Sending note data to backend: %s", string(jsonData))

	// Make POST request to backend
	resp, err := http.Post("http://localhost:8160/api/v1/notes", "application/json", strings.NewReader(string(jsonData)))
	if err != nil {
		return Note{}, fmt.Errorf("failed to create note: %v", err)
	}
	defer resp.Body.Close()

	log.Printf("Backend response status: %d", resp.StatusCode)

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return Note{}, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response to get the created note (id + content)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Note{}, fmt.Errorf("failed to read response: %v", err)
	}

	log.Printf("Backend response body: %s", string(body))

	var createdNote Note
	if err := json.Unmarshal(body, &createdNote); err != nil {
		// Fallback: server may have returned just id/name/hash; in that case populate with supplied content
		var fallback struct {
			ID      string `json:"id"`
			Name    string `json:"name"`
			Content string `json:"content"`
		}
		if err2 := json.Unmarshal(body, &fallback); err2 != nil {
			return Note{}, fmt.Errorf("failed to parse response: %v", err)
		}
		createdNote.ID = fallback.ID
		createdNote.Name = fallback.Name
		createdNote.Content = fallback.Content
	}

	return createdNote, nil
}

// fetchIndividualNote fetches a single note with content (still used by /api/notes/{id} if needed)
func fetchIndividualNote(noteID string) (string, error) {
	url := fmt.Sprintf("http://localhost:8160/api/v1/notes/%s", noteID)
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch note: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API returned status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %v", err)
	}

	var fullNote Note
	if err := json.Unmarshal(body, &fullNote); err != nil {
		return "", fmt.Errorf("failed to parse JSON: %v", err)
	}

	return fullNote.Content, nil
}

func main() {
	mux := http.NewServeMux()

	// Route for homepage
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Rendering Home Page")
		views.Index().Render(r.Context(), w)
	})

	// Route for notes page with API integration
	mux.HandleFunc("/notes", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Rendering Notes Page")

		// Fetch notes from the backend API (with content)
		notes, err := fetchNotesFromAPI()
		if err != nil {
			log.Printf("Error fetching notes: %v", err)
			notes = []views.Note{}
		}

		log.Printf("Fetched %d notes for listing", len(notes))

		views.IndexNotes(notes).Render(r.Context(), w)
	})

	// HTMX routes for creating notes and searching (no per-note content loading anymore)
	mux.HandleFunc("/api/notes/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("HTMX API Request: %s %s", r.Method, r.URL.Path)
		log.Printf("Request headers: %v", r.Header)

		// Extract note path tail
		path := strings.TrimPrefix(r.URL.Path, "/api/notes/")

		// Handle note creation (HTMX)
		if path == "create" && r.Method == "POST" {
			log.Println("Creating new note via HTMX")

			// Parse incoming form data (htmx posts form-encoded)
			if err := r.ParseForm(); err != nil {
				log.Printf("Error parsing form: %v", err)
				http.Error(w, "Invalid form data", http.StatusBadRequest)
				return
			}

			title := r.FormValue("title")
			content := r.FormValue("content")

			log.Printf("Form data - Title: %s, Content length: %d", title, len(content))

			if title == "" {
				log.Println("Title is empty")
				http.Error(w, "Title is required", http.StatusBadRequest)
				return
			}

			createdNote, err := createNoteViaAPI(title, content)
			if err != nil {
				log.Printf("Error creating note: %v", err)
				http.Error(w, "Failed to create note", http.StatusInternalServerError)
				return
			}

			// Render the new note card including content (so HTMX can prepend it)
			w.Header().Set("Content-Type", "text/html")
			views.NoteCard(createdNote.ID, createdNote.Name, createdNote.Content).Render(r.Context(), w)
			return
		}

		// Handle search (HTMX): /api/notes/search?q=...
		if path == "search" && r.Method == "GET" {
			query := strings.TrimSpace(r.URL.Query().Get("q"))
			qLower := strings.ToLower(query)

			// Get all notes (with content)
			allNotes, err := fetchNotesFromAPI()
			if err != nil {
				log.Printf("Search: error fetching notes: %v", err)
				allNotes = []views.Note{}
			}

			// If empty query, render all cards
			w.Header().Set("Content-Type", "text/html")
			if qLower == "" {
				for _, n := range allNotes {
					views.NoteCard(n.ID, n.Title, n.Content).Render(r.Context(), w)
				}
				return
			}

			// Filter by title or content
			var filtered []views.Note
			for _, n := range allNotes {
				if strings.Contains(strings.ToLower(n.Title), qLower) || strings.Contains(strings.ToLower(n.Content), qLower) {
					filtered = append(filtered, n)
				}
			}

			if len(filtered) == 0 {
				io.WriteString(w, `<div class="mx-6 my-6 text-gray-500 italic">No notes matched your search.</div>`)
				return
			}
			for _, n := range filtered {
				views.NoteCard(n.ID, n.Title, n.Content).Render(r.Context(), w)
			}
			return
		}

		log.Printf("No matching route for path: %s", path)
		http.NotFound(w, r)
	})

	// API proxy endpoints to handle CORS and forward requests to backend
	mux.HandleFunc("/api/v1/", func(w http.ResponseWriter, r *http.Request) {
		// Enable CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Forward request to backend
		backendURL := "http://localhost:8160" + r.URL.Path

		client := &http.Client{Timeout: 30 * time.Second}

		var req *http.Request
		var err error

		if r.Body != nil {
			req, err = http.NewRequest(r.Method, backendURL, r.Body)
		} else {
			req, err = http.NewRequest(r.Method, backendURL, nil)
		}

		if err != nil {
			http.Error(w, "Failed to create request", http.StatusInternalServerError)
			return
		}

		// Copy headers
		for key, values := range r.Header {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}

		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, "Failed to forward request", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		// Copy response headers
		for key, values := range resp.Header {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}

		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	})

	// Serve static files
	mux.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
	mux.Handle("/favicon.ico", http.FileServer(http.Dir("public")))

	log.Println("Listening on http://localhost:8080")
	if err := http.ListenAndServe("localhost:8080", mux); err != nil {
		log.Printf("error listening: %v", err)
	}
}
