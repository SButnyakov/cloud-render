package converters_test

import (
	"cloud-render/internal/lib/converters"
	"database/sql"
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestNullStringToString(t *testing.T) {
	var tests = []struct {
		name  string
		input sql.NullString
		want  string
	}{
		{"valid null string", sql.NullString{String: "test", Valid: true}, "test"},
		{"invalid null string", sql.NullString{String: "test", Valid: false}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := converters.NullStringToString(tt.input)
			assert.Equal(t, tt.want, res)
		})
	}
}

func TestStringToNullString(t *testing.T) {
	var tests = []struct {
		name  string
		input string
		want  sql.NullString
	}{
		{"empty string", "test", sql.NullString{String: "test", Valid: true}},
		{"not empty string", "", sql.NullString{String: "", Valid: false}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := converters.StringToNullString(tt.input)
			assert.Equal(t, tt.want, res)
		})
	}
}
