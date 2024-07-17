package api

// relations.go
type Relation struct {
	ArtistID   int    `json:"artist_id"`
	LocationID int    `json:"location_id"`
	DateID     int    `json:"date_id"`
	Type       string `json:"type"` // "upcoming" or "past"
}

// relations.go
func GetRelations() ([]Relation, error) {
	// ... (Fetch relation data from API)
	return nil, nil

}
