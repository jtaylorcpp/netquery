package queries

import (
	"github.com/jinzhu/gorm"
)

type query struct {
	db *gorm.DB
}

func NewQuery(db *gorm.DB) *query {
	return &query{db: db.New()}
}
