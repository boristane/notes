package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"

	_ "github.com/go-sql-driver/mysql"
)

var db *gorm.DB

var (
	dbName = os.Getenv("DB_NAME")
	dbPass = os.Getenv("DB_PASS")
	dbHost = os.Getenv("DB_HOST")
	dbPort = os.Getenv("DB_PORT")
)

func connectToDb() {
	log.Println("Connecting to the databse")
	dbSource := fmt.Sprintf("root:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true", dbPass, dbHost, dbPort, dbName)
	log.Println("Database source:", dbSource)

	var err error
	db, err = gorm.Open("mysql", dbSource)

	if err != nil {
		panic("Failed to connect to the database " + err.Error())
	}
	log.Println("Connection to the database established")
}

func migrateDb() {
	log.Println("Migrating the database to match model")
	db.AutoMigrate(&Note{}).AddUniqueIndex("idx_note_title_user", "title", "user_id")
}

func initialiseDb() {
	connectToDb()
	migrateDb()
	db.Debug()
}

func main() {
	initialiseDb()
	startServer()
}
