package todo

import (
	"database/sql"
)

type MockDB struct {
	callParams []interface{}
}

func (mdb *MockDB) Exec(query string, args ...interface{}) (*sql.Result, error) {
	return nil, nil
}

func (mdb *MockDB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return nil, nil
}

func (mdb *MockDB) QueryRow(query string, args ...interface{}) *sql.Row {
	return nil
}

type MockResult struct {
	InsertId int64
	Rows int64
}

func (mr *MockResult) LastInsertId() (int64, error) {
	return mr.InsertId, nil
}

func (mr *MockResult) RowsAffected() (int64, error) {
	return mr.Rows, nil
}
