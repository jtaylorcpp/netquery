package netquery

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func newPostgres(host, port, user, password, dbname string) *sql.DB {
	pqinfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", pqinfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	log.Println("database connected")

	return db
}
