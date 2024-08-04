// relations.go
package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Relation struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

type RelationResponse struct {
	Index []Relation `json:"index"`
}

func GetRelations() (map[int]map[string][]string, error) {
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/relation")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var relationResponse struct {
		Index []struct {
			ID             int                 `json:"id"`
			DatesLocations map[string][]string `json:"datesLocations"`
		} `json:"index"`
	}

	err = json.NewDecoder(resp.Body).Decode(&relationResponse)
	if err != nil {
		return nil, err
	}

	relations := make(map[int]map[string][]string)
	for _, item := range relationResponse.Index {
		relations[item.ID] = item.DatesLocations
	}

	return relations, nil
}

func GetRelationsForArtist(artistID int) (map[string][]string, error) {
	allRelations, err := GetRelations()
	if err != nil {
		return nil, err
	}
	if relations, ok := allRelations[artistID]; ok {
		return relations, nil
	}
	return nil, fmt.Errorf("no relations found for artist ID %d", artistID)
}
