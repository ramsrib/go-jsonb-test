package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"

	_ "github.com/golang-migrate/migrate/source/file" // required
)

func RunMigration(db *sql.DB, dbName string) {
	var migrationDir = "./migrations"

	// Run migrations
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		log.Fatalf("could not start sql migration... %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationDir), dbName, driver)

	if err != nil {
		log.Fatalf("migration failed... %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("An error occurred while syncing the database.. %v", err)
	}
	log.Println("Database migrated successfully!")
}

func WipeDatabase(db *sql.DB) {
	log.Println("Try to wipe the database...")
	// drop all the known tables
	db.Exec("DROP TABLE t1")
	db.Exec("DROP TABLE schema_migrations")
	log.Println("Wiped the database!")
}
