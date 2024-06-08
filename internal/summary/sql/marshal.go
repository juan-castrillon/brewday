package sql

import (
	"database/sql"
	"fmt"
	"strings"
)

// addToMarshalledArray is a helper method that will append a new object to a json array stored as TEXT in SQLite
// It accepts the column name where the array is stored and the string representation of the new object to add
func (s *SummaryPersistentStore) addToMarshalledArray(id, columnName, stringRep string) error {
	var current sql.NullString
	query := fmt.Sprintf(`SELECT %s FROM summaries WHERE recipe_id == ?`, columnName)
	err := s.dbClient.QueryRow(query, id).Scan(&current)
	if err != nil {
		return err
	}
	var newValue string
	if current.Valid {
		newValue = strings.Replace(current.String, "]", ","+stringRep+"]", 1)
	} else { //If it is NULL, it means its empty
		newValue = "[" + stringRep + "]"
	}
	updateQuery := fmt.Sprintf(`UPDATE summaries SET %s = ? WHERE recipe_id == ?`, columnName)
	_, err = s.dbClient.Exec(updateQuery, newValue, id)
	return err
}

// valueFromNullString is a helper method that will retrieve the value of a string stored in SQLite if not NULL
// If the value is null in the db it will return an empty string
func (s *SummaryPersistentStore) valueFromNullString(value sql.NullString) string {
	if value.Valid {
		return value.String
	}
	return ""
}

// sliceFromNullString is a helper method that will retrieve the value of a slice stored as a json string in SQLite if not NULL
// If the value is null in the db it will return the json representation of an empty slice
func (s *SummaryPersistentStore) sliceFromNullString(value sql.NullString) string {
	if value.Valid {
		return value.String
	}
	return "[]"
}

// valueFromNullFloat is a helper method that will retrieve the value of a float stored in SQLite if not NULL
// If the value is null in the db it will return 0
func (s *SummaryPersistentStore) valueFromNullFloat(value sql.NullFloat64) float32 {
	if value.Valid {
		return float32(value.Float64)
	}
	return 0
}

// valueFromNullInt is a helper method that will retrieve the value of an int stored in SQLite if not NULL
// If the value is null in the db it will return 0
func (s *SummaryPersistentStore) valueFromNullInt(value sql.NullInt32) int {
	if value.Valid {
		return int(value.Int32)
	}
	return 0
}
