package main

import (
	"github.com/joho/godotenv"
	"log"
)

func init() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Println(err)
		log.Fatal("Error loading .env file")
	}
}

func main() {
	getFirebaseInstance()
	getDbInstance()
	// defer db.Close()

	initializeWebServer()

}
