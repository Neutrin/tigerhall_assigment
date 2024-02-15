package repositiories

import "github.com/nitin/tigerhall/core/internal/model"

type TigerRepo interface {
	CreateTiger(model.Tiger, ...interface{}) (int, error)
	CreateTigerSighting(model.TigerSightings) (int, error)
	ListAllTigers(pagParams Pagination) (*Pagination, error)
}
