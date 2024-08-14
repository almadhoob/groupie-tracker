package api

import (
	"encoding/json"
	"log"
	"net/http"
)

type Artist struct {
	ID               int                 `json:"id"`
	Name             string              `json:"name"`
	Image            string              `json:"image"`
	Members          []string            `json:"members"`
	CreationDate     int                 `json:"creationDate"`
	FirstAlbum       string              `json:"firstAlbum"`
	ConcertDates     string              `json:"concertDates"`
	Relations        string              `json:"relations"` // Change this to string
	ConcertLocations map[string][]string `json:"-"`         // This will be populated later
	UniqueDates      []string            `json:"-"`
}

func GetArtists() ([]Artist, error) {
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		log.Printf("Error getting artists: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	var artists []Artist
	err = json.NewDecoder(resp.Body).Decode(&artists)
	if err != nil {
		return nil, err
	}

	return artists, nil
}
