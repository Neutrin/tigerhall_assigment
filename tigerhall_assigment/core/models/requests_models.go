package models

type TigerSightingReq struct {
	TigerId           string `json:"tiger_id" validate:"number"`
	LastSeenTimeStamp string `json:"last_seen" validate:"datetime=02-01-2006 15:04:05"`
	Lat               string `json:"latitude" validate:"latitude"`
	Long              string `json:"longitude" validate:"longitude"`
}

type CreateTigerReq struct {
	Name              string `json:"name" validate:"alphanum"`
	LastSeenTimeStamp string `json:"last_seen" validate:"datetime=02-01-2006 15:04:05"`
	Lat               string `json:"latitude" validate:"latitude"`
	Long              string `json:"longitude" validate:"longitude"`
	DOB               string `json:"date_of_birth" validate:"datetime=02-01-2006"`
}
