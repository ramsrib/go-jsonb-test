package main_test

import (
	"testing"

	"github.com/gofrs/uuid"

	app "go-jsonb-test"
)

// test with utf8 multi-byte 3 collation settings
func TestCreate_MB3(t *testing.T) {
	// setup the db
	config := app.GetDbConfig()

	// NOTE: `utf8mb3_unicode_ci` is same as `utf8_unicode_ci`
	config.Collation = "utf8_unicode_ci"
	db := app.GetDB(config)
	defer db.Close()

	// migrate the tables
	app.RunMigration(db, config.DBName)

	// setup the test data
	id := uuid.Must(uuid.NewV4()).Bytes()
	externalIdString := "09878a8d-41a3-4ae8-ba6b-f3f1b273317f"
	name := "test1"
	externalUUIDs := [][]byte{uuid.Must(uuid.FromString(externalIdString)).Bytes()}

	// run the test
	err := app.InsertRow(db, id, name, externalUUIDs)
	if err != nil {
		t.Errorf("Error occurred while inserting the entry with collation: %s, error: %v", config.Collation, err)
	}

	// read the value from db
	collation := ""
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
	t.Log()

	collation = config.Collation
	row1, err1 := app.ReadRow(db, id, collation)
	if err1 != nil {
		t.Errorf("Error occurred while reading the entry with collation: %s, error: %v", collation, err1)
	}
	t.Logf("Row value: %+v", row1)
	t.Log()

	collation = "utf8mb4_unicode_ci"
	row2, err2 := app.ReadRow(db, id, collation)
	if err2 != nil {
		t.Errorf("Error occurred while reading the entry with collation: %s, error: %v", collation, err2)
	}
	t.Logf("Row value: %+v", row2)
	if row2.ExternalUUIDs != nil {
		rowUUID := uuid.Must(uuid.FromBytes(row2.ExternalUUIDs[0])).String()
		t.Logf("Row External UUID[0] : %s", rowUUID)
		if rowUUID != externalIdString {
			t.Errorf("Read external UUID does not match with the inserted value, [inserted: %s, expected: %s]",
				externalIdString, rowUUID)
		}
	}

	collation = "binary"
	row3, err3 := app.ReadRow(db, id, collation)
	if err3 != nil {
		t.Errorf("Error occurred while reading the entry with collation: %s, error: %v", collation, err3)
	}
	t.Logf("Row value: %+v", row3)
	if row3.ExternalUUIDs != nil {
		rowUUID := uuid.Must(uuid.FromBytes(row3.ExternalUUIDs[0])).String()
		t.Logf("Row External UUID[0] : %s", rowUUID)
	}
	t.Log()

	// clean the db
	app.WipeDatabase(db)
}

// test with utf8 multi-byte 4 collation settings
func TestCreate_MB4(t *testing.T) {
	// setup the db
	config := app.GetDbConfig()

	config.Collation = "utf8mb4_unicode_ci"
	db := app.GetDB(config)
	defer db.Close()

	// migrate the tables
	app.RunMigration(db, config.DBName)

	// setup the test data
	id := uuid.Must(uuid.NewV4()).Bytes()
	externalIdString := "09878a8d-41a3-4ae8-ba6b-f3f1b273317f"
	name := "test1"
	externalUUIDs := [][]byte{uuid.Must(uuid.FromString(externalIdString)).Bytes()}

	// run the test
	err := app.InsertRow(db, id, name, externalUUIDs)
	if err != nil {
		t.Errorf("Error occurred while inserting the entry with collation: %s, error: %v", config.Collation, err)
	}

	// read the value from db
	collation := ""
	row, err := app.ReadRow(db, id, collation)
	if err != nil {
		t.Errorf("Error occurred while reading the entry without any collation: %s, error: %v", collation, err)
	}
	t.Logf("Row value: %+v", row)
	t.Log()

	collation = "utf8mb3_unicode_ci"
	row1, err1 := app.ReadRow(db, id, collation)
	if err1 != nil {
		t.Errorf("Error occurred while reading the entry with collation: %s, error: %v", collation, err1)
	}
	t.Logf("Row value: %+v", row1)
	if row1.ExternalUUIDs != nil {
		rowUUID := uuid.Must(uuid.FromBytes(row1.ExternalUUIDs[0])).String()
		t.Logf("Row External UUID[0] : %s", rowUUID)
		if rowUUID != externalIdString {
			t.Errorf("Read external UUID does not match with the inserted value, [inserted: %s, expected: %s]",
				externalIdString, rowUUID)
		}
	}

	// NOTE: this is the main scenario, i.e server and client matches the same collation,
	// but the number of bytes in the response are much higher.
	collation = config.Collation
	row2, err2 := app.ReadRow(db, id, collation)
	if err2 != nil {
		t.Errorf("Error occurred while reading the entry with collation: %s, error: %v", collation, err2)
	}
	t.Logf("Row value: %v", row2)

	collation = "binary"
	row3, err3 := app.ReadRow(db, id, collation)
	if err3 != nil {
		t.Errorf("Error occurred while reading the entry with collation: %s, error: %v", collation, err3)
	}
	t.Logf("Row value: %v", row3)

	// clean the db
	app.WipeDatabase(db)
}
