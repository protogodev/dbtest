package sql

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/RussellLuo/dbtest/spec"
)

type DB struct {
	db *sql.DB
}

func New(db *sql.DB) *DB {
	return &DB{db: db}
}

func (db *DB) Close() error {
	return db.db.Close()
}

func (db *DB) Insert(data map[string]spec.Rows) error {
	for tableName, rows := range data {
		for _, row := range rows {
			// len(row): 3 => placeholders: "?, ?, ?"
			placeholders := strings.TrimSuffix(strings.Repeat("?, ", len(row)), ", ")

			var names []string
			var values []interface{}
			for k, v := range row {
				names = append(names, k)
				values = append(values, v)
			}

			query := fmt.Sprintf(
				"INSERT INTO %s (%s) VALUES (%s)",
				tableName,
				strings.Join(names, ", "),
				placeholders,
			)
			stmt, err := db.db.Prepare(query)
			if err != nil {
				return err
			}
			defer stmt.Close()

			if _, err := stmt.Exec(values...); err != nil {
				return err
			}
		}
	}
	return nil
}

func (db *DB) Delete(keys ...string) error {
	for _, tableName := range keys {
		if _, err := db.db.Exec("TRUNCATE TABLE " + tableName); err != nil {
			return err
		}
	}
	return nil
}

func (db *DB) Select(query string) (result spec.Rows, err error) {
	rows, err := db.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	colTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		// Create a slice of interface{}'s to represent each column,
		// and a second slice to contain pointers to each item in the columns slice.
		//
		// See https://kylewbanks.com/blog/query-result-to-map-in-golang.
		columns := make([]interface{}, len(colTypes))
		columnPointers := make([]interface{}, len(colTypes))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		// Scan the result into the column pointers.
		if err := rows.Scan(columnPointers...); err != nil {
			return nil, err
		}

		// Store columns into a map.
		m := make(map[string]interface{})
		for i, ct := range colTypes {
			ptr := columnPointers[i].(*interface{})
			val, err := underlyingValue(ct.DatabaseTypeName(), *ptr)
			if err != nil {
				return nil, err
			}
			m[ct.Name()] = val
		}

		result = append(result, m)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func underlyingValue(typeName string, value interface{}) (interface{}, error) {
	v := value.([]byte)
	s := string(v)

	switch typeName {
	case "CHAR", "VARCHAR", "TEXT":
		return s, nil
	case "TINYINT", "INT", "BIGINT":
		i, err := strconv.Atoi(s)
		if err != nil {
			return nil, err
		}
		return i, err
	case "DECIMAL":
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return f, err
		}
		return f, nil
	case "BOOL":
		b, err := strconv.ParseBool(s)
		if err != nil {
			return nil, err
		}
		return b, nil
	case "DATE", "DATETIME":
		// Treat DATE or DATETIME values as strings for testing purpose.
		return s, nil
	default:
		return nil, fmt.Errorf("unsupported database type: %s", typeName)
	}
}
