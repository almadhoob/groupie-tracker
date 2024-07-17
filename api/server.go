package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func HandleArtists(w http.ResponseWriter, r *http.Request) {
	artists, err := GetArtists()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(artists)
}

func HandleLocations(w http.ResponseWriter, r *http.Request) {
	// ... (Similar to HandleArtists)
	locations, err := GetLocations()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(locations)
}

func HandleDates(w http.ResponseWriter, r *http.Request) {
	// ... (Similar to HandleArtists)
	dates, err := GetDates()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(dates)
}

func HandleRelations(w http.ResponseWriter, r *http.Request) {
	// ... (Similar to HandleArtists)
	relations, err := GetRelations()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(relations)
}

func StartServer() {
	http.HandleFunc("/artists", HandleArtists)
	http.HandleFunc("/locations", HandleLocations)
	http.HandleFunc("/dates", HandleDates)
	http.HandleFunc("/relations", HandleRelations)
	fmt.Println("Server started on port 8080")
	http.ListenAndServe(":8080", nil)
}
