package repositiories

import "github.com/nitin/tigerhall/core/internal/model"

type TigerRepo interface {
	CreateTiger(model.Tiger, ...interface{}) (int, error)
	CreateTigerSighting(model.TigerSightings) (int, error)
	ListAllTigers(Pagination) (*Pagination, error)
	TigerById(int) (model.Tiger, error)
	ListSightings(int, Pagination) (*Pagination, error)
	InitTigerCreateMap() map[int]map[int]struct{}
}
