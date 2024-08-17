package api

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var templates *template.Template

func init() {
	var err error
	templates, err = setupTemplates()
	if err != nil {
		log.Fatalf("Failed to set up templates: %v", err)
	}
}

func setupTemplates() (*template.Template, error) {
	funcMap := template.FuncMap{
		"join":      strings.Join,
		"split":     strings.Split,
		"trimSpace": strings.TrimSpace,
	}

	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("error getting current working directory: %v", err)
	}

	templatesDir := filepath.Join(cwd, "templates")

	t, err := template.New("").Funcs(funcMap).ParseGlob(filepath.Join(templatesDir, "*.html"))
	if err != nil {
		return nil, fmt.Errorf("error parsing templates: %v", err)
	}

	return t, nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) error {
	err := templates.ExecuteTemplate(w, tmpl, data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		return fmt.Errorf("500 - failed to execute template: %v", err)
	}
	return nil
}
