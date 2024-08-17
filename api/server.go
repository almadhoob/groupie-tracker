package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func HandleArtists(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" && r.URL.Path != "/artists" {
		http.Error(w, "404 - Page not found", http.StatusNotFound)
		return
	}

	artists, err := GetArtists()
	if err != nil {
		log.Printf("Error getting artists: %v", err)
		http.Error(w, "500 - Failed to fetch artists data", http.StatusInternalServerError)
		return
	}

	relations, err := GetRelations()
	if err != nil {
		log.Printf("Error getting relations: %v", err)
		http.Error(w, "500 - Failed to get relations", http.StatusInternalServerError)
		return
	}

	for i, artist := range artists {
		if relData, ok := relations[artist.ID]; ok {
			artists[i].ConcertLocations = relData.DatesLocations
			artists[i].UniqueDates = relData.UniqueAllDates
		}
	}

	err = renderTemplate(w, "index.html", artists)
	if err != nil {
		if strings.HasPrefix(err.Error(), "404") {
			log.Printf("Error rendering template: %v", err)
			http.Error(w, "404 - Page not found", http.StatusNotFound)
		} else if strings.HasPrefix(err.Error(), "500") {
			log.Printf("Error rendering template: %v", err)
			http.Error(w, "500 - Failed to render index page", http.StatusInternalServerError)
		}
		return
	}
}

func HandleRelations(w http.ResponseWriter, r *http.Request) {
	relations, err := GetRelations()
	if err != nil {
		log.Printf("Error getting relations: %v", err)
		http.Error(w, "500 - Failed to get relations", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(relations)
}

func HandleArtistDetail(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/artist/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Error rendering template: %v", err)
		http.Error(w, "404 - Page not found", http.StatusNotFound)
		return
	}

	artist, err := GetArtistByID(id)
	if err != nil {
		log.Printf("Error rendering template: %v", err)
		http.Error(w, "404 - Page not found", http.StatusNotFound)
		return
	}

	relation, err := GetRelationsForArtist(artist.ID)
	if err != nil {
		log.Printf("Error getting relations for artist %d: %v", artist.ID, err)
		http.Error(w, "500 - Failed to fetch artist relations", http.StatusInternalServerError)
		return
	}

	artist.ConcertLocations = relation.DatesLocations
	artist.UniqueDates = relation.UniqueAllDates

	err = renderTemplate(w, "artist.html", artist)
	if err != nil {
		if strings.HasPrefix(err.Error(), "404") {
			log.Printf("Error rendering template: %v", err)
			http.Error(w, "404 - Page not found", http.StatusNotFound)
		} else if strings.HasPrefix(err.Error(), "500") {
			log.Printf("Error rendering template: %v", err)
			http.Error(w, "500 - Failed to render artist page", http.StatusInternalServerError)
		}
		return
	}
}

func GetArtistByID(id int) (Artist, error) {
	artists, err := GetArtists()
	if err != nil {
		return Artist{}, fmt.Errorf("500 - Failed to fetch artists data: %v", err)
	}

	for _, artist := range artists {
		if artist.ID == id {
			return artist, nil
		}
	}
	return Artist{}, fmt.Errorf("404 - artist with ID %d not found", id)
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/favico.ico")
}

func SetupRoutes() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	mux.HandleFunc("/favicon.ico", faviconHandler)
	mux.HandleFunc("/", HandleArtists)
	mux.HandleFunc("/artists", HandleArtists)
	mux.HandleFunc("/artist/", HandleArtistDetail)
	return mux
}
