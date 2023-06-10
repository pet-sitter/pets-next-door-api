package firebaseinfra

import (
	"context"
	"log"

	firebase "firebase.google.com/go"
	_ "firebase.google.com/go/auth"

	"google.golang.org/api/option"
)

type FirebaseApp struct {
	*firebase.App
}

func NewFirebaseApp(pathToCredentialsFile string) *FirebaseApp {
	opt := option.WithCredentialsFile(pathToCredentialsFile)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	return &FirebaseApp{app}
}
