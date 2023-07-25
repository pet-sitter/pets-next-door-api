package configs

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

var Port = os.Getenv("PORT")

var DatabaseURL = os.Getenv("DATABASE_URL")
var MigrationPath = os.Getenv("MIGRATION_PATH")

var KakaoRestAPIKey = os.Getenv("KAKAO_REST_API_KEY")
var KakaoRedirectURI = os.Getenv("KAKAO_REDIRECT_URI")

var FirebaseCredentialsPath = os.Getenv("FIREBASE_CREDENTIALS_PATH")

func init() {
	if Port == "" {
		Port = "8080"
	}

	if DatabaseURL == "" {
		panic("DATABASE_URL is required")
	}

	if MigrationPath == "" {
		MigrationPath = "db/migrations"
	}

	if KakaoRestAPIKey == "" {
		panic("KAKAO_REST_API_KEY is required")
	}

	if KakaoRedirectURI == "" {
		panic("KAKAO_REDIRECT_URI is required")
	}

	if FirebaseCredentialsPath == "" {
		FirebaseCredentialsPath = "firebase-credentials.json"
	}
}
