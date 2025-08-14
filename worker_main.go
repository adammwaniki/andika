package main

import (
    "context"
    "log"
    "os"
    "path/filepath"

    "github.com/a-h/templ"
    "github.com/adammwaniki/andika/views"
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

func main() {
    // Home page
    renderPage("dist/index.html", views.Index())

    // Projects page
    renderPage("dist/projects/index.html", views.AllProjects())

    // NavPanel
    renderPage("dist/header/closed/index.html", views.Header()) // when closed
    renderPage("dist/header/open/index.html", views.Header()) // when open

    log.Println("Static export completed!")
}