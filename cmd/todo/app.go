package todo

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

type GeneralErr interface {
	Error() string
}

type Err struct {
	Code int
	Msg string
}

type Health struct {
	Msg string
}

type AppError struct {
	Err
}

type DBError struct {
	Err
}

func newAppErr(msg string) *AppError {
	return &AppError{
		Err{
		Code: 1,
		Msg:  msg,
	}}
}

func newDBError(msg string) *DBError {
	return &DBError{
		Err{
			Code: 0,
			Msg:  msg,
		},
	}
}

func (ae *AppError) Error() string {
	return ae.Msg
}

func (de *DBError) Error() string {
	return de.Msg
}

func (a *App) Initialize(user, pass, host, dbname string) {
	connectionString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		host, user, pass, dbname)
	var err error
	a.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) initializeRoutes(){
	a.Router.HandleFunc("/healthz", a.Healthz).Methods("GET")
	a.Router.HandleFunc("/todos", a.getTodos).Methods("GET")
	a.Router.HandleFunc("/todos", a.createTodo).Methods("POST")
	a.Router.HandleFunc("/todos/{id:[0-9]+}", a.getTodo).Methods("GET")
	a.Router.HandleFunc("/todos/{id:[0-9]+}", a.updateTodo).Methods("PUT")
	a.Router.HandleFunc("/todos/{id:[0-9]+}", a.deleteProduct).Methods("DELETE")
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(":8080", a.Router))
}

func (a *App) Healthz(w http.ResponseWriter, r *http.Request) {
	log.Info("API application Health is OK")
	log.Info("checking DB health")
	if err := a.CheckDB(); err != nil {
		log.Errorf("Error checking database: %v", err)
		respondWithError(w, http.StatusInternalServerError, newDBError(err.Error()))
	}
	h := Health{Msg: "everything working ok"}
	respondWithJSON(w, http.StatusOK, h)
}

func (a *App) CheckDB() error {
	t := Todo{}
	if err := t.checkDB(a.DB); err != nil {
		return err
	}
	return nil
}

func (a *App) getTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, newAppErr("Invalid Todo ID"))
		return
	}

	t := Todo{ID: id}
	if err := t.getTodo(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, newDBError("Todo not found"))
		default:
			respondWithError(w, http.StatusInternalServerError, newAppErr(err.Error()))
		}
		return
	}

	respondWithJSON(w, http.StatusOK, t)
}

func (a *App) getTodos(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}

	todos, err := getTodos(a.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, newAppErr(err.Error()))
	}
	respondWithJSON(w, http.StatusOK, todos)
}

func (a *App) createTodo(w http.ResponseWriter, r *http.Request) {
	var t Todo
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&t); err != nil {
		respondWithError(w, http.StatusBadRequest, newAppErr("Invalid request payload"))
		return
	}
	defer r.Body.Close()

	if err := t.newTodo(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, newAppErr(err.Error()))
	}

	respondWithJSON(w, http.StatusCreated, t)
}

func (a *App) updateTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, newAppErr("Invalid Todo ID"))
	}

	var t Todo
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&t); err != nil {
		msg := fmt.Sprintf("Invalid todo payload json error: %v", err)
		respondWithError(w, http.StatusBadRequest, newAppErr(msg))
		return
	}

	defer r.Body.Close()

	t.ID = id

	if err := t.updateTodo(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, newAppErr(err.Error()))
		return
	}

	respondWithJSON(w, http.StatusOK, t)
}

func (a *App) deleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err :=strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, newAppErr(err.Error()))
		return
	}

	t := Todo{ID: id}
	if err := t.deleteTodo(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, newDBError(err.Error()))
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})

}

func respondWithError(w http.ResponseWriter, code int, err GeneralErr) {
	respondWithJSON(w, code, err)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}){
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

