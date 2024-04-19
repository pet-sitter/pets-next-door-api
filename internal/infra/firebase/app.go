package firebaseinfra

import (
	"context"
	"encoding/json"
	"fmt"

	firebase "firebase.google.com/go"

	// Firebase Auth initialization
	_ "firebase.google.com/go/auth"
	"github.com/pet-sitter/pets-next-door-api/internal/configs"
	"golang.org/x/oauth2/google"

	"google.golang.org/api/option"
)

type FirebaseApp struct {
	*firebase.App
}

func NewFirebaseAppFromCredentialsPath(firebaseCredentialsPath string) (*FirebaseApp, error) {
	opt := option.WithCredentialsFile(firebaseCredentialsPath)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, fmt.Errorf("error initializing app: %v", err)
	}

	return &FirebaseApp{app}, nil
}

func NewFirebaseAppFromCredentialsJSON(firebaseCredentialsJSON configs.FirebaseCredentialsJSONType) (*FirebaseApp, error) {
	firebaseCredentialsJSONBytes, err := json.Marshal(firebaseCredentialsJSON)
	if err != nil {
		return nil, fmt.Errorf("error marshalling firebase credentials json: %v", err)
	}

	opt := option.WithCredentials(
		&google.Credentials{
			ProjectID: firebaseCredentialsJSON.ProjectID,
			JSON:      firebaseCredentialsJSONBytes,
		},
	)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, fmt.Errorf("error initializing app: %v", err)
	}

	return &FirebaseApp{app}, nil
}
