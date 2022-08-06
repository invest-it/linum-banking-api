package main

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"log"
	"os"
)

var db *sql.DB

func getDbInstance() *sql.DB {
	fmt.Println("Get instance")
	if db == nil {
		lock.Lock()
		defer lock.Unlock()
		if db == nil {
			var err error
			fmt.Println(os.Getenv("DATABASE_URL"))
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
	statement := `INSERT INTO UserRequisitions (RequisitionId, UserId) VALUES ($1, $2)`
	_, err := db.Exec(statement, requisitionId, uid)
	if err != nil {
		return err
	}
	return nil
}

func getRequisitionsForUser(uid string) ([]string, error) {
	db := getDbInstance()
	statement := `SELECT RequisitionId FROM UserRequisitions WHERE UserId=$1`
	rows, err := db.Query(statement, uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	requisitionIds := make([]string, 3) // Most sensible size for requisitions
	for rows.Next() {
		var requisitionId string
		err := rows.Scan(&requisitionId)
		if err != nil {
			return nil, err
		}
		requisitionIds = append(requisitionIds, requisitionId)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return requisitionIds, nil
}

func getRequisitionForUser(requisitionId string, uid string) ([]string, error) {
	db := getDbInstance()
	rows, err := db.Query("SELECT RequsitionId FROM UserRequisitions WHERE  UserId=? AND RequisitionId=?", uid, requisitionId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	requisitionIds := make([]string, 3) // Most sensible size for requisitions
	for rows.Next() {
		var requisitionId string
		err := rows.Scan(&requisitionId)
		if err != nil {
			return nil, err
		}
		requisitionIds = append(requisitionIds, requisitionId)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return requisitionIds, nil
}
