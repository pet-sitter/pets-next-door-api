package main

import (
	"log"
	"net/http"
	"time"

	"github.com/rs/zerolog"

	"github.com/pet-sitter/pets-next-door-api/internal/configs"
	firebaseinfra "github.com/pet-sitter/pets-next-door-api/internal/infra/firebase"
	_ "github.com/pet-sitter/pets-next-door-api/pkg/docs"
)

// @title 이웃집멍냥 API 문서
// @version 0.12.0
// @description 이웃집멍냥 백엔드 API 문서입니다.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url
// @contact.email petsnextdoordev@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host
// @BasePath /api

// @securityDefinitions.apiKey FirebaseAuth
// @in header
// @name Authorization
func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	var app *firebaseinfra.FirebaseApp
	var err error
	if configs.GetFirebaseCredentialsJSON() != (configs.FirebaseCredentialsJSONType{}) {
		app, err = firebaseinfra.NewFirebaseAppFromCredentialsJSON(
			configs.GetFirebaseCredentialsJSON(),
		)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		app, err = firebaseinfra.NewFirebaseAppFromCredentialsPath(configs.FirebaseCredentialsPath)
		if err != nil {
			log.Fatal(err)
		}
	}

	r, err := NewRouter(app)
	if err != nil {
		log.Fatal(err)
	}

	server := &http.Server{
		Addr:              ":" + configs.Port,
		Handler:           r,
		ReadHeaderTimeout: 3 * time.Second,
	}

	log.Printf("Starting server on port %s", configs.Port)
	log.Fatal(server.ListenAndServe())
}
