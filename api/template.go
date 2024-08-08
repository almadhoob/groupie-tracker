// templates.go
package api

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var templates *template.Template

func init() {
	templates = setupTemplates()
}

func setupTemplates() *template.Template {
	funcMap := template.FuncMap{
		"join":      strings.Join,
		"split":     strings.Split,
		"trimSpace": strings.TrimSpace,
	}

	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current working directory: %v", err)
	}

	// Construct the full path to the templates directory
	templatesDir := filepath.Join(cwd, "templates")

	// Parse all templates and add the function map
	t, err := template.New("").Funcs(funcMap).ParseGlob(filepath.Join(templatesDir, "*.html"))
	if err != nil {
		log.Fatalf("Error parsing templates: %v", err)
	}

	return t
}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	err := templates.ExecuteTemplate(w, tmpl, data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func headersWritten(w http.ResponseWriter) bool {
	if rw, ok := w.(interface{ Written() bool }); ok {
		return rw.Written()
	}
	return false
}
