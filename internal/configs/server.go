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

var B2KeyID = os.Getenv("B2_APPLICATION_KEY_ID")
var B2Key = os.Getenv("B2_APPLICATION_KEY")
var B2BucketName = os.Getenv("B2_BUCKET_NAME")
var B2Endpoint = os.Getenv("B2_ENDPOINT")
var B2Region = os.Getenv("B2_REGION")

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

	if B2KeyID == "" {
		panic("B2_APPLICATION_KEY_ID is required")
	}

	if B2Key == "" {
		panic("B2_APPLICATION_KEY is required")
	}

	if B2BucketName == "" {
		panic("B2_BUCKET_NAME is required")
	}

	if B2Endpoint == "" {
		panic("B2_ENDPOINT is required")
	}

	if B2Region == "" {
		panic("B2_REGION is required")
	}
}
