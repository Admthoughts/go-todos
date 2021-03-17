package todo

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type Todo struct {
	Done      bool `json:"done"`
	ID        int
	Text      string    `json:"text"`
	CreatedOn time.Time `json:"created_on"`
}

type Error struct {
	Msg  string `json:"msg"`
	Body string `json:",omitempty"`
}

type Response struct {
	Todos     []Todo    `json:"todos,omitempty"`
	Total     int       `json:"total_count"`
	RCreation time.Time `json:"request_time"`
	Status    string    `json:"status"`
	Error     `json:"error,omitempty"`
}

type SQLDB interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Ping() error
}

type SQLResult interface {
	LastInsertId() (int64, error)
	RowsAffected() (int64, error)
}

func writeResponse(res Response, w *http.ResponseWriter) {
	bytes, err := json.Marshal(res)
	if err != nil {
		logrus.Error(err)
		http.Error(*w, err.Error(), http.StatusInternalServerError)
	}

	fmt.Fprintf(*w, "%s\n", bytes)
}

func (t *Todo) checkDB(db SQLDB) error {
	if err := db.Ping(); err != nil {
		return err
	}
	return nil
}

func (t *Todo) newTodo(db SQLDB) error {
	err := db.QueryRow(
		"INSERT INTO todos(text, done, createdon, updatedon) VALUES($1, $2, current_timestamp, current_timestamp) RETURNING id",
		t.Text, t.Done).Scan(&t.ID)
	if err != nil {
		return err
	}
	return nil
}

func (t *Todo) getTodo(db SQLDB) error {
	return db.QueryRow("SELECT text, done, CreatedOn FROM todos WHERE id=$1", t.ID).Scan(&t.Text,
		&t.Done, &t.CreatedOn)
}

func getTodos(db SQLDB, start, count int) ([]Todo, error) {
	rows, err := db.Query(
		"SELECT id, text, done, CreatedOn FROM todos LIMIT $1 OFFSET $2",
		count, start)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	todos := []Todo{}

	for rows.Next() {
		var t Todo
		if err := rows.Scan(&t.ID, &t.Text, &t.Done, &t.CreatedOn); err != nil {
			return nil, err
		}
		todos = append(todos, t)
	}

	return todos, nil

}

func (t *Todo) updateTodo(db SQLDB) error {
	_, err := db.Exec(
		"UPDATE todos SET text=$1, done=$2, UpdatedOn=current_timestamp WHERE id=$3", t.Text, t.Done, t.ID)
	return err
}

func (t *Todo) deleteTodo(db SQLDB) error {
	_, err := db.Exec(
		"DELETE FROM todos WHERE id=$1", t.ID)
	return err
}

