package orm

import (
	"gorm.io/gorm"
)

// Paginate gorm分页，默认page=1，size=20
func Paginate(page int, limit int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page < 1 {
			page = 1
		}
		switch {
		case limit > 100:
			limit = 100
		case limit < 1:
			limit = 20
		}
		offset := (page - 1) * limit
		return db.Offset(offset).Limit(limit)
	}
}
