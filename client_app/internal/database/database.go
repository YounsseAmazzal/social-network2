package database

import (
	"database/sql"
	"log"
	"social/internal/models"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
		_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

func NeewDb(dbPath string) *models.DB {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal("Errror from Databassse opening", err.Error())
	}
	if err := db.Ping(); err != nil {
		log.Fatal("Error from Databassse pinging  ", err.Error())
	}
	RuunThemigration(db)
	///daaabase conected
	return &models.DB{
		Database: db,
	}
}
func RuunThemigration(dbinstant *sql.DB) {
	driver, err := sqlite3.WithInstance(dbinstant, &sqlite3.Config{})
	if err != nil {
		log.Fatal("migration driver error ", err.Error())
	}
	mig, err := migrate.NewWithDatabaseInstance("file://sqlinteraction/migrations",
		"sqlite3",
		driver,
	)
	if err != nil {
		log.Fatal("migration init func not work", err.Error())
	}
	if err := mig.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("migration upp error !", err.Error())
	}
	log.Println("everrything is good  lett's works ")
}
