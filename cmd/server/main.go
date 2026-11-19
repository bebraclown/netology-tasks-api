package main

import (
	"log"
	"net/http"

	"tasks-api/internal/handlers"
	"tasks-api/internal/storage"
)

func main() {
	const address = "127.0.0.1:8080"
	log.Printf("сервер слушает http://%s", address)
	log.Fatal(http.ListenAndServe(address, handlers.New(storage.NewMemory())))
}
