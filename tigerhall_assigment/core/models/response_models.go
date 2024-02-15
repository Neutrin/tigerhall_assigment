package models

type TigerResp struct {
	Name     string `json:"name"`
	LastSeen string `json:"last_seen_at"`
	Lat      string `json:"latitude"`
	Long     string `json:"longitute"`
}
