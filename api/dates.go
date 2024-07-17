package api

type Date struct {
	Date string `json:"date"`
}

// dates.go
func GetDates() ([]Date, error) {
	// ... (Fetch date data from API)
	return nil, nil

}
