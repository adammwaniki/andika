package main

import (
	"log"
	"net/http"

	"github.com/adammwaniki/andika/views"
)

func main() {
    mux := http.NewServeMux()

    // Route for homepage
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Rendering Home Page")
		views.Index().Render(r.Context(), w)
    })
    
	// Route for projects page
    mux.HandleFunc("/notes/new", func(w http.ResponseWriter, r *http.Request) {
        log.Println("Rendering New Notes Page")
		views.NewNote().Render(r.Context(), w)
    })

    // Serve Tailwind CSS and any other assets e.g., static, public etc.
    mux.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
    mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.Handle("/favicon.ico", http.FileServer(http.Dir("public")))


    log.Println("Listening on http://localhost:8080")
    http.ListenAndServe(":8080", mux)
}