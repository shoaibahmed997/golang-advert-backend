package database

import (
	"database/sql"
	"fmt"
	"log"
)

var err error
var DB *sql.DB

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = 1234
	dbname   = "advert_db"
)

func postTable() {
	stmt, stmtErr := DB.Prepare("CREATE TABLE IF NOT EXISTS posts (ID TEXT PRIMARY KEY,Title TEXT NOT NULL,Description TEXT NOT NULL,Price INT NOT NULL,Category TEXT NOT NULL,Location TEXT NOT NULL,lattitude NUMERIC,longitude NUMERIC, userEmail TEXT NOT NULL, by TEXT , CreatedAt TIMESTAMP DEFAULT NOW(), FOREIGN KEY (userEmail) REFERENCES users (email)  ON DELETE CASCADE );")
	if stmtErr != nil {
		log.Fatal(stmtErr)
	}
	_, resErr := stmt.Exec()
	if stmtErr != nil {
		log.Fatal(resErr)
	}

}

func imagesTable() {
	stmt, stmtErr := DB.Prepare("CREATE TABLE IF NOT EXISTS images (ID TEXT PRIMARY KEY, imgpath TEXT NOT NULL, PostID TEXT NOT NULL, FOREIGN KEY (PostID) REFERENCES posts (ID) ON DELETE CASCADE );")
	if stmtErr != nil {
		log.Fatal(stmtErr)
	}
	_, resErr := stmt.Exec()
	if stmtErr != nil {
		log.Fatal(resErr)
	}
}

func DBInit() {
	psqlConURI := fmt.Sprintf("host=%s port=%d user=%s password=%d dbname=%s sslmode=disable", host, port, user, password, dbname)
	DB, err = sql.Open("postgres", psqlConURI)
	if err != nil {
		log.Println(err)
	}

	// _, err1 := DB.Exec("PRAGMA foreign_keys = ON", nil)
	// if err1 != nil {
	// 	log.Println(err1)
	// }

	stmt, stmtError := DB.Prepare("CREATE TABLE IF NOT EXISTS users (ID TEXT PRIMARY KEY, Firstname TEXT NOT NULL,Lastname TEXT, Email TEXT NOT NULL UNIQUE, Password TEXT NOT NULL, CreatedAt TIMESTAMP DEFAULT NOW());")
	if stmtError != nil {
		log.Println("error stmt", stmtError)
	}
	// stmt.Exec()
	_, resErr := stmt.Exec()
	if resErr != nil {
		log.Println("res error", err)
	}
	postTable()
	imagesTable()
	log.Println("Connected to database successfully ")
}
