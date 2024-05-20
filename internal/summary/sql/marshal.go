package sql

import (
	"database/sql"
	"fmt"
	"strings"
)

func (s *SummaryRecorderPersistentStore) addToMarshalledArray(id, columnName, stringRep string) error {
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

func (s *SummaryRecorderPersistentStore) ValueFromNullString(value sql.NullString) string {
	if value.Valid {
		return value.String
	}
	return ""
}

func (s *SummaryRecorderPersistentStore) ValueFromNullFloat(value sql.NullFloat64) float32 {
	if value.Valid {
		return float32(value.Float64)
	}
	return 0
}

func (s *SummaryRecorderPersistentStore) ValueFromNullInt(value sql.NullInt32) int {
	if value.Valid {
		return int(value.Int32)
	}
	return 0
}
