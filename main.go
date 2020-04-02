package main

import (
	"log"

	"github.com/gofrs/uuid"
)

func main() {
	log.Println("Starting the app!...")

	config := GetDbConfig()
	db := GetDB(config)
	defer db.Close()

	// run the migration before the apo
	RunMigration(db, config.DBName)

	// actual app logic begins here

	// setup  the test data
	id := uuid.Must(uuid.NewV4()).Bytes()
	externalIdString := "09878a8d-41a3-4ae8-ba6b-f3f1b273317f"
	name := "test1"
	externalUUIDs := [][]byte{uuid.Must(uuid.FromString(externalIdString)).Bytes()}

	// insert the externalId and read it back and compare.
	err := InsertRow(db, id, name, externalUUIDs)
	if err != nil {
		log.Println("Error occurred while inserting the entry", err)
	}

	// without any collation in the select query
	collation := ""
	row, err := ReadRow(db, id, collation)
	if err != nil {
		log.Printf("Error occurred while reading the entry wihtout using collation: %s, error :%s", collation, err)
	}
	log.Printf("Row object (with collation: %s), %v\n", collation, row)

	// our main scenario
	collation = "utf8mb4_unicode_ci"
	row1, err1 := ReadRow(db, id, collation)
	if err1 != nil {
		log.Printf("Error occurred while reading the entry using collation: %s, error :%s", collation, err1)
	}
	log.Printf("Row object (with collation: %s), %v\n", collation, row1)

	if row1.ExternalUUIDs != nil {
		log.Printf("External UUID[0]: %s, collation: %s\n",
			uuid.Must(uuid.FromBytes(row1.ExternalUUIDs[0])).String(), collation)
	}

	// another scenario
	collation = "utf8_unicode_ci"
	row2, err2 := ReadRow(db, id, collation)
	if err2 != nil {
		log.Printf("Error occurred while reading the entry using collation: %s, error :%s", collation, err2)
	}
	log.Printf("Row object (with collation: %s), %v\n", collation, row2)

	log.Println("Exiting the app!")
}
