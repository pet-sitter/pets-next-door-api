package main

import (
	"log"
	"net/http"

	"github.com/pet-sitter/pets-next-door-api/internal/configs"
	firebaseinfra "github.com/pet-sitter/pets-next-door-api/internal/infra/firebase"
	"github.com/pet-sitter/pets-next-door-api/internal/server"
)

func main() {
	app := firebaseinfra.NewFirebaseApp(configs.FirebaseCredentialsPath)
	r := server.NewRouter(app)

	log.Printf("Starting server on port %s", configs.Port)
	log.Fatal(http.ListenAndServe(":"+configs.Port, r))
}
