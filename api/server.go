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
	artists, err := GetArtists()
	if err != nil {
		log.Printf("Error getting artists: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
			artists[i].ConcertLocations = relData
		}
	}

	renderTemplate(w, "artists.html", artists)
}

func extractDates(relations map[string][]string) []string {
	var dates []string
	for _, dateList := range relations {
		dates = append(dates, dateList...)
	}
	return removeDuplicates(dates)
}

func removeDuplicates(strSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range strSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// func HandleLocations(w http.ResponseWriter, r *http.Request) {
// 	locations, err := GetLocations()
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	json.NewEncoder(w).Encode(locations)
// }

// func HandleDates(w http.ResponseWriter, r *http.Request) {
// 	dates, err := GetDates()
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	json.NewEncoder(w).Encode(dates)
// }

func HandleRelations(w http.ResponseWriter, r *http.Request) {
	relations, err := GetRelations()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(relations)
}

func (a *Artist) GetConcertDatesSlice() []string {
	return strings.Split(a.ConcertDates, ",")
}

func HandleArtistDetail(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/artist/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid artist ID", http.StatusBadRequest)
		return
	}

	artist, err := GetArtistByID(id)
	if err != nil {
		http.Error(w, "Artist not found", http.StatusNotFound)
		return
	}

	relations, err := GetRelations()
	if err != nil {
		log.Printf("Error getting relations: %v", err)
		http.Error(w, "Error fetching artist data", http.StatusInternalServerError)
		return
	}

	if relData, ok := relations[artist.ID]; ok {
		artist.ConcertLocations = relData
	}

	// Directly call renderTemplate without attempting to capture an error
	renderTemplate(w, "artist.html", artist)
}

func GetArtistByID(id int) (Artist, error) {
	artists, err := GetArtists()
	if err != nil {
		return Artist{}, fmt.Errorf("error fetching artists: %v", err)
	}

	for _, artist := range artists {
		if artist.ID == id {
			return artist, nil
		}
	}

	return Artist{}, fmt.Errorf("artist with ID %d not found", id)
}

func SetupRoutes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", HandleArtists)
	mux.HandleFunc("/artists", HandleArtists)
	mux.HandleFunc("/artist/", HandleArtistDetail) // New route for artist detail
	return mux
}
