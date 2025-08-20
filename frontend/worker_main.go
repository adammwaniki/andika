package main

import (
    "context"
    "log"
    "os"
    "path/filepath"

    "github.com/a-h/templ"
    "github.com/adammwaniki/andika/frontend/views"
)

func renderPage(path string, page templ.Component) {
    // Ensure the directory exists
    err := os.MkdirAll(filepath.Dir(path), 0755)
    if err != nil {
        log.Fatal(err)
    }

    f, err := os.Create(path)
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()

    if err := page.Render(context.Background(), f); err != nil {
        log.Fatal(err)
    }
}

func wmain() {
    // Home page
    renderPage("dist/index.html", views.Index())

    // New Notes page
    renderPage("dist/notes/new/index.html", views.NewNote())

    log.Println("Static export completed!")
}