package main

import (
	"database/sql"
	_ "github.com/jackc/pgx/v4/stdlib"
	"log"
	"os"
)

var db *sql.DB

func getDbInstance() *sql.DB {
	if db == nil {
		lock.Lock()
		defer lock.Unlock()
		if firebaseApp == nil {
			var err error
			db, err = sql.Open("pgx", os.Getenv("DATABASE_URL"))
			if err != nil {
				log.Fatalf("Could not setup database connection")
			}
		}
	}
	return db
}

func storeRequisitionId(requisitionId string, uid string) error {
	db := getDbInstance()
	_, err := db.Exec("INSERT INTO UserRequisitions (RequisitionId, UserId) VALUES (?, ?)", requisitionId, uid)
	if err != nil {
		return err
	}
	return nil
}
