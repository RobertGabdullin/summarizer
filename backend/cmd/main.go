package main

import (
	"log"
	"net/http"

	"github.com/RobertGabdullin/summarizer/internal/handlers"
)

func main() {
	http.Handle("/api/process", &handlers.Process{})
	http.Handle("/", handlers.NewStatic())

	log.Println("Listening on http://localhost:1717")
	http.ListenAndServe(":1717", nil)
}
