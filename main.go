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

	collations := []string{"", "utf8_general_ci", "utf8_unicode_ci", "utf8mb4_unicode_ci", "utf8mb3_bin", "utf8mb4_bin"}
	for _, collation := range collations {
		row, err := ReadRow(db, id, collation)
		if err != nil {
			log.Printf("Error occurred while reading the entry using collation: %s, error :%s", collation, err)
		}
		log.Printf("Row object (with collation: %s), %v\n", collation, row)

		if row.ExternalUUIDs != nil {
			log.Printf("External UUID[0]: %s, collation: %s\n",
				uuid.Must(uuid.FromBytes(row.ExternalUUIDs[0])).String(), collation)
		}
	}

	log.Println("Exiting the app!")
}
