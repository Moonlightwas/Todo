package models

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

type Todo struct {
	ID          int       `sql:"type=SERIAL PRIMARY KEY"`
	Priority    string    `sql:"type=CHAR(1), NOT NULL"`
	CreatedAt   time.Time `sql:"type=TIMESTAMP, NOT NULL"`
	Description string    `sql:"type=TEXT, NOT NULL"`
	Project     string    `sql:"type=VARCHAR(256)"`
	Context     string    `sql:"type=VARCHAR(256)"`
	Deadline    time.Time `sql:"type=TIMESTAMP"`
}

func GenerateParams(model interface{}) []string {
	t := reflect.TypeOf(model)
	params := make([]string, t.NumField())

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		sqlType := GoToSQL(field)
		if sqlType != "" {
			params[i] = fmt.Sprintf("%s %s", strings.ToLower(string(t.Field(i).Name)), sqlType)
		}
	}

	return params
}

func GoToSQL(f reflect.StructField) string {
	sqlTag := f.Tag.Get("sql")
	if sqlTag != "" { //checking sql tag in model struct
		sqlTag = strings.TrimLeft(sqlTag, "type=")
		sqlTag = strings.ToUpper(sqlTag)
		tagParts := strings.Split(sqlTag, ", ")

		return strings.Join(tagParts, " ")
	}

	switch f.Type.Kind() {
	case reflect.Int, reflect.Int32, reflect.Int64:
		if f.Name == "Priority" {
			return "CHAR(1)"
		}
		if f.Name == "ID" {
			return "ID SERIAL PRIMARY KEY"
		}
		return "INTEGER"
	case reflect.String:
		return "VARCHAR(256)"
	case reflect.Bool:
		return "BOOLEAN"
	case reflect.Struct:
		if f.Type == reflect.TypeOf(time.Time{}) {
			return "TIMESTAMP"
		}
	}
	return ""
}
