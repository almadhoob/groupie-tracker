package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"
)

type Relation struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

type ProcessedRelation struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
	UniqueAllDates []string            `json:"uniqueAllDates"`
}

func GetRelations() (map[int]ProcessedRelation, error) {
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/relation")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var relationResponse struct {
		Index []Relation `json:"index"`
	}

	err = json.NewDecoder(resp.Body).Decode(&relationResponse)
	if err != nil {
		return nil, err
	}

	processedRelations := make(map[int]ProcessedRelation)
	for _, relation := range relationResponse.Index {
		processed := ProcessedRelation{
			ID:             relation.ID,
			DatesLocations: make(map[string][]string),
		}

		allDates := []string{}
		for location, dates := range relation.DatesLocations {
			sortedDates := sortDates(dates)
			processed.DatesLocations[location] = sortedDates
			allDates = append(allDates, sortedDates...)
		}

		processed.UniqueAllDates = removeDuplicatesAndSort(allDates)
		processedRelations[relation.ID] = processed
	}

	return processedRelations, nil
}

func GetRelationsForArtist(artistID int) (ProcessedRelation, error) {
	allRelations, err := GetRelations()
	if err != nil {
		return ProcessedRelation{}, err
	}
	if relation, ok := allRelations[artistID]; ok {
		return relation, nil
	}
	return ProcessedRelation{}, fmt.Errorf("no relations found for artist ID %d", artistID)
}

func removeDuplicatesAndSort(strSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range strSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return sortDates(list)
}

func sortDates(dates []string) []string {
	sort.Slice(dates, func(i, j int) bool {
		date1, _ := time.Parse("02-01-2006", dates[i])
		date2, _ := time.Parse("02-01-2006", dates[j])
		return date2.Before(date1)
	})
	return dates
}
