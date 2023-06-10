package main

import (
	"log"
	"net/http"

	"github.com/pet-sitter/pets-next-door-api/internal/configs"
	"github.com/pet-sitter/pets-next-door-api/internal/server"
)

func main() {
	r := server.NewRouter()

	log.Printf("Starting server on port %s", configs.Port)
	log.Fatal(http.ListenAndServe(":"+configs.Port, r))
}
