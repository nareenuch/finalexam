package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

func Conn() *sql.DB {
	if db != nil {
		fmt.Println("db exist")
		return db
	}

	var err error
	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("connect to database error", err)
	}
	return db
}


