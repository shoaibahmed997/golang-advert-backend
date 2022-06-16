package database

import (
	"database/sql"
	"log"
)

var err error
var DB *sql.DB

func postTable() {
	stmt, stmtErr := DB.Prepare("CREATE TABLE IF NOT EXISTS posts (ID TEXT PRIMARY KEY,Title TEXT NOT NULL,Desc TEXT NOT NULL,Price INTEGER NOT NULL,Category TEXT NOT NULL,Location TEXT NOT NULL,lattitude INTEGER,longitude INTEGER, userEmail TEXT NOT NULL,by TEXT,CreatedAt DEFAULT CURRENT_TIMESTAMP, FOREIGN KEY (userEmail) REFERENCES users (email)  ON DELETE CASCADE );")
	if stmtErr != nil {
		log.Fatal(stmtErr)
	}
	_, resErr := stmt.Exec()
	if stmtErr != nil {
		log.Fatal(resErr)
	}

}

func imagesTable() {
	stmt, stmtErr := DB.Prepare("CREATE TABLE IF NOT EXISTS images (ID TEXT PRIMARY KEY, imgpath TEXT, PostID INTEGER, FOREIGN KEY (PostID) REFERENCES posts (ID) ON DELETE CASCADE );")
	if stmtErr != nil {
		log.Fatal(stmtErr)
	}
	_, resErr := stmt.Exec()
	if stmtErr != nil {
		log.Fatal(resErr)
	}
}

func DBInit() {
	DB, err = sql.Open("sqlite3", "./data.db")
	if err != nil {
		log.Println(err)
	}

	_, err1 := DB.Exec("PRAGMA foreign_keys = ON", nil)
	if err1 != nil {
		log.Println(err1)
	}

	stmt, stmtError := DB.Prepare("CREATE TABLE IF NOT EXISTS users (ID TEXT PRIMARY KEY, Firstname TEXT NOT NULL,Lastname TEXT, Email TEXT NOT NULL UNIQUE, Password TEXT NOT NULL, CreatedAt DEFAULT CURRENT_TIMESTAMP);")
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
