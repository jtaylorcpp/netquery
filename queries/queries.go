package queries

import (
	"github.com/jinzhu/gorm"
)

type Query struct {
	db *gorm.DB
}

func NewQuery(db *gorm.DB) *Query {
	return &Query{db: db.New()}
}
