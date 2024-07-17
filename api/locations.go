package api

// locations.go
type Location struct {
	Name string `json:"name"`
}

func GetLocations() ([]Location, error) {
	// ... (Fetch location data from API)
	return nil, nil

}
