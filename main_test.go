package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"
)

var (
	a App
)

const (
	tableName = "todos"
	tableCreationQuery = `CREATE TABLE IF NOT EXISTS todos
(
    id SERIAL,
    text TEXT NOT NULL,
    done BOOLEAN NOT NULL DEFAULT 'f',
    CreatedOn timestamp NOT NULL,
    UpdatedOn timestamp NOT NULL,
    CONSTRAINT todos_pkey PRIMARY KEY (id)
)`
)

func TestMain(m *testing.M) {
	a.Initialize(
		os.Getenv("TODO_USER"),
		os.Getenv("TODO_PASS"),
		tableName)
	ensureTableExists()
	code := m.Run()
	clearTable()
	os.Exit(code)
}

func ensureTableExists() {
	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	a.DB.Exec("DELETE FROM todos")
	a.DB.Exec("ALTER SEQUENCE todos_id_seq RESTART WITH 1")
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d, got %d\n", expected, actual)
	}
}

func TestEmptyTable(t *testing.T) {
	clearTable()
	req, _ := http.NewRequest("GET", "/todos", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func TestNonExistentTodo(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/todos/11", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	if err := json.Unmarshal(response.Body.Bytes(), &m); err != nil {
		t.Errorf("error unmarshalling response to map, response: %v", response.Body.String())
	}
	if len(m["error"]) == 0 {
		t.Errorf("Expected an error key in response, none found. Response body: %v", response.Body.String())
	}
	if m["error"] != "Todo not found" {
		t.Errorf("Expected the 'error' key to be set to 'Todo not found'. Got %v", m["error"])
	}

}

func TestCreateTodo(t *testing.T) {
	clearTable()

	var jsonStr = []byte(`{"text":"test todo", "done": false}`)
	req, _ := http.NewRequest("POST", "/todos", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}

	json.Unmarshal(response.Body.Bytes(), &m)

	if m["text"] != "test todo" {
		t.Errorf("Expected text to be 'test todo', Got: %v", m["text"])
	}
	if m["ID"] != 1.0 {
		t.Errorf("Got 0 for ID")
	}

}

func randomBool() bool {
	rand.Seed(time.Now().UnixNano())
	return rand.Float32() < 0.5
}

func TestGetTodo(t *testing.T) {
	clearTable()
	addTodos(1)

	req, _ := http.NewRequest("GET", "/todos/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
}

func addTodos(count int) {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		a.DB.Exec(
			fmt.Sprintf(
				"INSERT INTO %s(text, done, createdon, updatedon) VALUES($1, $2, current_timestamp, current_timestamp)",
				tableName),
			"Todo "+strconv.Itoa(i), randomBool())
	}
}

func TestUpdateTodo(t *testing.T) {
	clearTable()
	addTodos(1)

	req, _ := http.NewRequest("GET", "/todos/1", nil)
	response := executeRequest(req)

	var originalTodo map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalTodo)

	var jsonStr = []byte(`{"text": "test todo - updated text", "done": true}`)
	req, _ = http.NewRequest("PUT", "/todos/1", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
	if response.Code != 200 {
		t.Errorf("Bad response code: %v response body: %v", response.Code, response.Body.String())
	}

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["id"] != originalTodo["id"] {
		t.Errorf("Expect the id to remain the same (%v). Got %v", originalTodo["id"], m["id"])
	}
	if m["text"] == originalTodo["name"] {
		t.Errorf("expected the name to change")
	}

}

func TestDeleteProduct(t *testing.T) {
	clearTable()
	addTodos(1)

	req, _ := http.NewRequest("GET", "/todos/1", nil)

	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/todos/1", nil)
	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/todos/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}