package converters

import "database/sql"

func NullStringToString(nullString sql.NullString) string {
	if nullString.Valid {
		return nullString.String
	}
	return ""
}

func StringToNullString(str string) sql.NullString {
	return sql.NullString{String: str, Valid: str != ""}
}
