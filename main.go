package main

import (
	"github.com/Admthoughts/go-todo/todo"
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {
	log.Info("Configuring TODOList API server")
	a := todo.App{}
	log.Info("Server initialisation starting")
	a.Initialize(
		os.Getenv("TODO_USER"),
		os.Getenv("TODO_PASS"),
		os.Getenv("TODO_DB_HOST"),
		os.Getenv("TODO_DBNAME"),
	)
	if err := a.CheckDB(); err != nil {
		log.Fatalf("error checking database after initialisation: %v", err)
	}
	log.Debug("Server initialised")
	log.Info("Starting server")
	a.Run(":8080")
}
