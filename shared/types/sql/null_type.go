package sql

import (
	"database/sql"
	"time"
)

type SQL struct{}

func (s *SQL) NullStringToPointer(nullString sql.NullString) *string {
	if nullString.Valid {
		return &nullString.String
	}
	return nil
}

func (s *SQL) NullTimeToPointer(nullTime sql.NullTime) *time.Time {
	if nullTime.Valid {
		return &nullTime.Time
	}
	return nil
}

func (s *SQL) StringToNullString(string *string) sql.NullString {
	if s != nil {
		return sql.NullString{String: *string, Valid: true}
	}
	return sql.NullString{Valid: false}
}

func (s *SQL) TimeToNullTime(time *time.Time) sql.NullTime {
	if time != nil {
		return sql.NullTime{Time: *time, Valid: true}
	}
	return sql.NullTime{Valid: false}
}
