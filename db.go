package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-sql-driver/mysql"
)

func GetDbConfig() mysql.Config {
	config := mysql.Config{
		Addr:                 getEnv("DB_HOST", "127.0.0.1") + ":" + getEnv("DB_PORT", "3309"),
		User:                 getEnv("DB_USERNAME", "go-jsonb-test"),
		Passwd:               getEnv("DB_PASSWORD", "go-jsonb-test"),
		DBName:               getEnv("DB_NAME", "go-jsonb-test"),
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
		//Collation:            "utf8_general_ci", // default value in golang sql driver
		//Collation:            "utf8_unicode_ci", // also works (uses utf8mb3 charset)
		Collation: "utf8mb4_unicode_ci", // doesn't work, but server's config
	}
	return config
}

func GetDB(config mysql.Config) *sql.DB {
	db, err := sql.Open("mysql", config.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
	// ensure the connection really works
	if err := db.Ping(); err != nil {
		log.Fatal("Got an error while pinging the database:", err)
	}
	log.Println("Successfully connected to the database!")
	return db
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func InsertRow(db *sql.DB, id []byte, name string, externalUUIDs [][]byte) error {
	log.Println("Trying to insert the data into t1...")

	insertQuery := "INSERT INTO t1 (id, name, external_uuids) VALUES(?, ?, JSON_ARRAY(%s))"

	args := make([]interface{}, 0, len(externalUUIDs)+2)
	args = append(args, id, name)

	externalUUIDSParams := make([]string, 0, len(externalUUIDs))
	for _, uid := range externalUUIDs {
		args = append(args, uid)
		externalUUIDSParams = append(externalUUIDSParams, "?")
	}
	jsonArrayParams := strings.Join(externalUUIDSParams, ",")

	query := fmt.Sprintf(insertQuery, jsonArrayParams)
	log.Println("Executing query:", query)

	result, err := db.Exec(query, args...)
	if err != nil {
		return err
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected != 1 {
		return errors.New("invalid rows affected: " + string(rowsAffected))
	}
	log.Println("Inserted successfully.")
	return nil
}

func ReadRow(db *sql.DB, id []byte, collation string) (T1Row, error) {
	log.Printf("Trying to read the row with '%s' collation...", collation)
	var row T1Row
	var jsonArray JSONArray

	selectQuery := "SELECT id, name, external_uuids FROM t1 WHERE id = ?"
	if collation != "" {
		selectQuery = fmt.Sprintf("SELECT id, name, external_uuids COLLATE `%s` FROM t1 WHERE id = ?", collation)
	}

	if err := db.QueryRow(selectQuery, id).Scan(&row.Id, &row.Name, &jsonArray); err != nil {
		return T1Row{}, err
	}
	row.ExternalUUIDs = jsonArray.ByteValues
	log.Println("Read the row successfully...")
	return row, nil
}

// Row object from table t1
type T1Row struct {
	Id            []byte
	Name          string
	ExternalUUIDs [][]byte
}
