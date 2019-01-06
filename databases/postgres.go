package databases

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	bp "github.com/jtaylorcpp/broparser"
)

func NewPostgresDB(host, port, user, password, dbname string) *gorm.DB {
	pqinfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := gorm.Open("postgres", pqinfo)
	if err != nil {
		panic(err)
	}

	log.Println("database connected")

	return db
}

func PSQLInit(db *gorm.DB) {
	db.AutoMigrate(&bp.ConnRecord{})
	db.AutoMigrate(&bp.DNSAnswer{})
	db.AutoMigrate(&bp.DNSRecord{})
}

func PSQLClear(db *gorm.DB) {
	db.DropTableIfExists(&bp.ConnRecord{})
	db.DropTableIfExists(&bp.DNSAnswer{})
	db.DropTableIfExists(&bp.DNSRecord{})
	log.Println("all data cleared")
}
