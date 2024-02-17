package model

import (
	"time"

	"gorm.io/gorm"
)

type Tiger struct {
	gorm.Model
	Name              string    `gorm:"unique"`
	DOB               time.Time `gorm:"column:date_of_birth;type:timestamp"`
	LastSeenTimeStamp int64     `gorm:"column:last_seen;type:BIGINT"`
	Lat               float64   `gorm:"column:latititude;type:Decimal(8,6)"`
	Long              float64   `gorm:"column:longitude;type:Decimal(9,6)"`
	CreatedBy         uint      `gorm:"column:created_by;type:varchar(100)"`
	User              User      `gorm:"foreignKey:CreatedBy"`
}

type TigerSightings struct {
	gorm.Model
	Tiger             Tiger   `gorm:"foreignKey:TigerId"`
	TigerId           int     `gorm:"column:tiger_id"`
	LastSeenTimeStamp int64   `gorm:"column:last_seen;type:BIGINT"`
	Lat               float64 `gorm:"column:latititude;type:Decimal(8,6)"`
	Long              float64 `gorm:"column:longitude;type:Decimal(9,6)"`
	ImagePath         string  `gorm:"column:image_path;type:varchar(200)"`
	CreatedBy         uint    `gorm:"column:created_by;type:varchar(100)"`
	User              User    `gorm:"foreignKey:CreatedBy"`
}
