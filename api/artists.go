// artists.go
package api

type Artist struct {
	Name       string   `json:"name"`
	Image      string   `json:"image"`
	BeginYear  int      `json:"begin_year"`
	FirstAlbum string   `json:"first_album"`
	Members    []string `json:"members"`
}

func GetArtists() ([]Artist, error) {
	// ... (Fetch artist data from API)
	return nil, nil
}
