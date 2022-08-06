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

func storeRequisitionId(requisitionId string, requisitionReference, uid string) error {
	db := getDbInstance()
	statement := `INSERT INTO UserRequisitions (RequisitionId, RequisitionReference, UserId, Approved) VALUES ($1, $2, $3, false)`
	_, err := db.Exec(statement, requisitionId, requisitionReference, uid)
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

func userHasRequisition(requisitionId string, uid string) bool {
	db := getDbInstance()
	statement := "SELECT RequsitionId FROM UserRequisitions WHERE  UserId=$1 AND RequisitionId=$2"
	row := db.QueryRow(statement, uid, requisitionId)
	var result string
	if err := row.Scan(&result); err != nil {
		return false
	}
	return result == requisitionId
}

func getRequisitionByReference(requisitionReference string) (string, error) {
	db := getDbInstance()
	statement := "SELECT RequisitionId FROM UserRequisitions WHERE ReferenceId=$1"
	row := db.QueryRow(statement, requisitionReference)
	var requisitionId string
	if err := row.Scan(&requisitionId); err != nil {
		return "", err
	}
	return requisitionId, nil
}

func updateRequisitionState(requisitionReference string) error {
	db := getDbInstance()
	statement := "UPDATE UserRequisitions SET Approved=true WHERE RequisitionReference=$1"
	_, err := db.Exec(statement, requisitionReference)
	return err
}
