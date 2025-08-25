package main

import (
	"log"
	"net/http"

	"github.com/adammwaniki/andika/vcs/handler"
)

func main() {
	mux := http.NewServeMux()

	// Notes
	mux.HandleFunc("/notes/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handler.CreateNoteHandler(w, r)
		case http.MethodGet:
			if r.URL.Path == "/notes/" || r.URL.Path == "/notes" {
				handler.ListNotesHandler(w, r)
			} else if len(r.URL.Path) > len("/notes/") && r.URL.Path[len(r.URL.Path)-10:] == "/snapshots" {
				handler.ListSnapshotsHandler(w, r)
			} else {
				handler.ViewNoteHandler(w, r)
			}
		case http.MethodPut:
			handler.EditNoteHandler(w, r)
		default:
			http.NotFound(w, r)
		}
	})

	// Snapshots
	mux.HandleFunc("/snapshots/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handler.ViewSnapshotHandler(w, r)
			return
		}
		if r.Method == http.MethodPost && len(r.URL.Path) > len("/snapshots/") && r.URL.Path[len(r.URL.Path)-8:] == "/restore" {
			handler.RestoreSnapshotHandler(w, r)
			return
		}
		http.NotFound(w, r)
	})

	// Help
	mux.HandleFunc("GET /help", handler.HelpHandler)

	log.Println("VCS Notes API running on localhost:8160")
	if err := http.ListenAndServe(":8160", mux); err != nil {
		log.Fatal(err)
	}
}
