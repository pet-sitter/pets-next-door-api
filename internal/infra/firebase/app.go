package firebaseinfra

import (
	"context"
	"encoding/json"
	"log"

	firebase "firebase.google.com/go"
	_ "firebase.google.com/go/auth"
	"github.com/pet-sitter/pets-next-door-api/internal/configs"
	"golang.org/x/oauth2/google"

	"google.golang.org/api/option"
)

type FirebaseApp struct {
	*firebase.App
}

func NewFirebaseAppFromCredentialsPath(firebaseCredentialsPath string) *FirebaseApp {
	opt := option.WithCredentialsFile(firebaseCredentialsPath)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	return &FirebaseApp{app}
}

func NewFirebaseAppFromCredentialsJSON(firebaseCredentialsJSON configs.FirebaseCredentialsJSONType) *FirebaseApp {
	firebaseCredentialsJSONBytes, err := json.Marshal(firebaseCredentialsJSON)
	if err != nil {
		log.Fatalf("error marshalling firebase credentials json: %v\n", err)
	}

	opt := option.WithCredentials(
		&google.Credentials{
			ProjectID: firebaseCredentialsJSON.ProjectID,
			JSON:      firebaseCredentialsJSONBytes,
		},
	)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	return &FirebaseApp{app}
}
