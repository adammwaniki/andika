// backend/main.go
package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/adammwaniki/andika/backend/handler"
)

// ParamHandler type for handlers that need path parameters
type ParamHandler func(w http.ResponseWriter, r *http.Request, params map[string]string)

// CORS middleware
func corsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Allow local dev and production
        origin := r.Header.Get("Origin")
        if origin == "http://localhost:8080" || origin == "https://andika.adamndegwa.workers.dev" {
            w.Header().Set("Access-Control-Allow-Origin", origin)
        }

        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        w.Header().Set("Access-Control-Max-Age", "86400")

        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }

        next.ServeHTTP(w, r)
    })
}

// wrapper to convert ParamHandler into http.HandlerFunc
func withParams(pattern string, h ParamHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, pattern)
		path = strings.Trim(path, "/")

		segments := []string{}
		if path != "" {
			segments = strings.Split(path, "/")
		}

		params := make(map[string]string)
		if len(segments) >= 1 {
			params["id"] = segments[0]
		}
		// Only set sid if we have enough segments and it's a snapshot route
		if len(segments) >= 3 && segments[1] == "snapshots" {
			params["sid"] = segments[2] // The actual hash
		}

		h(w, r, params)
	}
}

func main() {
	mux := http.NewServeMux()
	apiRouter := http.NewServeMux()

	// /notes/search route MUST come before /notes/ catch-all
	apiRouter.HandleFunc("/notes/search", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handler.SearchNotesHandler(w, r)
			return
		}
		http.NotFound(w, r)
	})

	// /notes collection
	apiRouter.HandleFunc("/notes", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handler.CreateNoteHandler(w, r)
		case http.MethodGet:
			handler.ListNotesHandler(w, r)
		default:
			http.NotFound(w, r)
		}
	})

	// /notes/{id} and nested snapshots
	apiRouter.HandleFunc("/notes/", withParams("/notes/", func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		segments := strings.Split(strings.Trim(strings.TrimPrefix(r.URL.Path, "/notes/"), "/"), "/")

		// Add bounds checking
		if len(segments) == 0 {
			http.NotFound(w, r)
			return
		}

		// /notes/{id}/snapshots routes
		if len(segments) >= 2 && segments[1] == "snapshots" {
			switch {
			case len(segments) == 2 && r.Method == http.MethodGet:
				// GET /notes/{id}/snapshots
				handler.ListSnapshotsHandler(w, r, params)
				return
			case len(segments) == 3 && r.Method == http.MethodGet:
				// GET /notes/{id}/snapshots/{sid}
				handler.ViewSnapshotHandler(w, r, params)
				return
			case len(segments) == 4 && len(segments) > 3 && segments[3] == "restore" && r.Method == http.MethodPost:
				// POST /notes/{id}/snapshots/{sid}/restore
				handler.RestoreSnapshotHandler(w, r, params)
				return
			default:
				http.NotFound(w, r)
				return
			}
		}

		// /notes/{id} routes (only if we have exactly 1 segment)
		if len(segments) == 1 {
			switch r.Method {
			case http.MethodGet:
				handler.ViewNoteHandler(w, r)
			case http.MethodPut:
				handler.EditNoteHandler(w, r)
			case http.MethodDelete:
				handler.DeleteNoteHandler(w, r)
			default:
				http.NotFound(w, r)
			}
		} else {
			http.NotFound(w, r)
		}
	}))

	// /help
	apiRouter.HandleFunc("/help", handler.HelpHandler)

	// Mount versioned API at /api/v1/ with CORS middleware
	mux.Handle("/api/v1/", corsMiddleware(http.StripPrefix("/api/v1", apiRouter)))

	// Redirect /api/v1 → /api/v1/
	mux.HandleFunc("/api/v1", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/api/v1/", http.StatusPermanentRedirect)
	})

	log.Println("Andika API running on http://localhost:8160")
	if err := http.ListenAndServe(":8160", mux); err != nil {
		log.Fatal(err)
	}
}
