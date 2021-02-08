package main

import (
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {
	log.Info("Configuring TODOList API server")
	a := App{}
	log.Info("Server initialisation starting")
	a.Initialize(
		os.Getenv("TODO_USER"),
		os.Getenv("TODO_PASS"),
		os.Getenv("TODO_DBNAME"),
	)
	log.Debug("Server initialised")
	log.Info("Starting server")
	a.Run(":8080")
}
