package database

import (
	"database/sql"
	"fmt"
	"reflect"
)

type Scanner interface {
	Scan(dest ...any) error
}

type Scanable interface {
	Scan() []any
}

func ToStruct[T Scanable](s Scanner) (*T, error) {
	var v T
	if err := s.Scan(v.Scan()...); err != nil {
		return nil, err
	}
	return &v, nil
}

func ScanRowToStruct(rows *sql.Rows, dest any) error {
	v := reflect.ValueOf(dest)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("dest must be pointer to struct")
	}

	v = v.Elem()
	t := v.Type()

	cols, err := rows.Columns()
	if err != nil {
		return err
	}

	fieldMap := make(map[string]int)
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		tag := f.Tag.Get("db")
		if tag != "" {
			fieldMap[tag] = i
		}
	}

	scanArgs := make([]any, len(cols))
	for i, col := range cols {
		if idx, ok := fieldMap[col]; ok {
			scanArgs[i] = v.Field(idx).Addr().Interface()
		} else {
			var discard any
			scanArgs[i] = &discard
		}
	}

	return rows.Scan(scanArgs...)
}
