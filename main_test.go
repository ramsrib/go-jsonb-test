package main_test

import (
	"database/sql"
	"testing"

	"github.com/gofrs/uuid"

	app "go-jsonb-test"
)

func setupTestCase(collation string) (*sql.DB, func()) {
	config := app.GetDbConfig()
	config.Collation = collation
	db := app.GetDB(config)

	// migrate the tables
	app.RunMigration(db, config.DBName)

	// tear down function
	return db, func() {
		// clean the db
		defer db.Close()
		// wipe the db at the end of the test case
		app.WipeDatabase(db)
	}
}

// test with utf8 multi-byte 3 collation config
func TestReadRow(t *testing.T) {

	// setup the test data
	id := uuid.Must(uuid.NewV4()).Bytes()
	externalIdString := "09878a8d-41a3-4ae8-ba6b-f3f1b273317f"
	name := "test1"
	externalUUIDs := [][]byte{uuid.Must(uuid.FromString(externalIdString)).Bytes()}

	// NOTE: `utf8mb3_unicode_ci` is same as `utf8_unicode_ci`
	serverCollations := []string{"utf8_unicode_ci", "utf8mb4_unicode_ci"}

	for _, sCollations := range serverCollations {
		db, tearDownTestCase := setupTestCase(sCollations)
		defer tearDownTestCase()

		t.Run(sCollations, func(t *testing.T) {
			// insert the test data
			err := app.InsertRow(db, id, name, externalUUIDs)
			if err != nil {
				t.Errorf("Error occurred while inserting the entry with collation: %s, error: %v", sCollations, err)
			}

			// run the test with different collations
			collations := []string{"", "utf8_unicode_ci", "utf8mb4_unicode_ci", "utf8mb3_bin", "utf8mb4_bin", "binary"}
			for _, collation := range collations {
				t.Run(collation, func(t *testing.T) {
					row, err := app.ReadRow(db, id, collation)
					if err != nil {
						t.Errorf("Error occurred while reading the entry without any collation: %s, error: %v", collation, err)
					}
					t.Logf("Row value: %+v", row)
					if row.ExternalUUIDs != nil {
						rowUUID := uuid.Must(uuid.FromBytes(row.ExternalUUIDs[0])).String()
						t.Logf("Row External UUID[0] : %s", rowUUID)
						if rowUUID != externalIdString {
							t.Errorf("Read external UUID does not match with the inserted value, [inserted: %s, expected: %s]",
								externalIdString, rowUUID)
						}
					}
				})
			}
		})
	}
}
