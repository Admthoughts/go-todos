package main

import (
	"os"
)

func main() {
	a := App{}
	a.Initialize(
		os.Getenv("TODO_USER"),
		os.Getenv("TODO_PASS"),
		os.Getenv("TODO_DBNAME"),
	)

	a.Run(":8080")

}
