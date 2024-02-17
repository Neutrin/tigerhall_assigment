package models

type TigerResp struct {
	Name     string `json:"name"`
	LastSeen string `json:"last_seen_at"`
	Lat      string `json:"latitude"`
	Long     string `json:"longitute"`
}

type TigerSightingResp struct {
	TigerId   int        `json:"tiger_id"`
	Name      string     `json:"name"`
	Sightings []Sighting `json:"sightings"`
}

type Sighting struct {
	Lat      string `json:"latitude"`
	Long     string `json:"longitute"`
	LastSeen string `json:"last_seen_at"`
	Image    string `json:"image"`
}
