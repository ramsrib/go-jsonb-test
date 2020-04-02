package main

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gofrs/uuid"
)

// This JSONArray implementation is not driver-specific
type JSONArray struct {
	Values     []string
	ByteValues [][]byte
}

// Scan implements the Scanner interface.
// The value type must be []byte,
// otherwise Scan fails.
func (nt *JSONArray) Scan(value interface{}) (err error) {
	if value == nil {
		nt.Values = nil
		return
	}

	switch v := value.(type) {
	case []byte:
		var arr []string
		if err := json.Unmarshal(v, &arr); err != nil {
			log.Printf("Error unmarshalling: %v", err)
			return err
		}
		uuidBytes, err := parseUUIDs(arr)
		if err != nil {
			return err
		}
		nt.Values = arr
		nt.ByteValues = uuidBytes
		return nil
	}

	return fmt.Errorf("can't convert %T to JSONArray", value)
}

// Value implements the driver Valuer interface.
func (nt JSONArray) Value() (driver.Value, error) {
	return nt.Values, nil
}

func parseUUIDs(jsonBytes []string) ([][]byte, error) {
	uuidBytes := make([][]byte, 0, len(jsonBytes))
	for _, id := range jsonBytes {
		uid, err := uuid.FromBytes([]byte(id))
		if err != nil {
			return nil, fmt.Errorf("unable to convert the uuid from string to bytes, value={%v}, %w", id, err)
		}
		uuidBytes = append(uuidBytes, uid.Bytes())
	}
	return uuidBytes, nil
}
