package repositiories

import (
	"math"

	"github.com/nitin/tigerhall/core/internal/repositiories"
	"gorm.io/gorm"
)

func Paginate(value interface{}, pagination *repositiories.Pagination, db *gorm.DB,
	whereCond map[string]interface{}) func(db *gorm.DB) *gorm.DB {
	var (
		totalRows int64
		tx        *gorm.DB
	)
	tx = db.Debug().Model(value)
	if len(whereCond) > 0 {
		tx = tx.Where(whereCond)
	}
	tx.Count(&totalRows)
	pagination.TotalRows = totalRows
	totalPages := int(math.Ceil(float64(totalRows) / float64(pagination.GetLimit())))
	pagination.TotalPages = totalPages

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit()).Order(pagination.GetSort())
	}
}
