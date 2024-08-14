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
		Render404(w)
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for i, artist := range artists {
		if relData, ok := relations[artist.ID]; ok {
			artists[i].ConcertLocations = relData.DatesLocations
			artists[i].UniqueDates = relData.UniqueAllDates
		}
	}
	err = renderTemplate(w, "artists.html", artists)
	if err != nil {
		if strings.HasPrefix(err.Error(), "404") {
			Render404(w)
		} else {
			log.Printf("Error rendering template: %v", err)
			http.Error(w, "500 - Failed to render artists page", http.StatusInternalServerError)
		}
		return
	}
}

func HandleRelations(w http.ResponseWriter, r *http.Request) {
	relations, err := GetRelations()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(relations)
}

func HandleArtistDetail(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/artist/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		Render404(w)
		return
	}

	artist, err := GetArtistByID(id)
	if err != nil {
		http.Error(w, "404 - Artist not found", http.StatusNotFound)
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
			Render404(w)
		} else {
			log.Printf("Error rendering template: %v", err)
			http.Error(w, "500 - Failed to render artist detail page", http.StatusInternalServerError)
		}
		return
	}
}

func GetArtistByID(id int) (Artist, error) {
	artists, err := GetArtists()
	if err != nil {
		return Artist{}, fmt.Errorf("500 - error fetching artists: %v", err)
	}

	for _, artist := range artists {
		if artist.ID == id {
			return artist, nil
		}
	}

	return Artist{}, fmt.Errorf("404 - artist with ID %d not found", id)
}

func SetupRoutes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", HandleArtists)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	mux.HandleFunc("/artists", HandleArtists)
	mux.HandleFunc("/artist/", HandleArtistDetail)
	return mux
}
