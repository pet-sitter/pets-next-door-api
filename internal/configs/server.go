package configs

import (
	"os"
	"strings"

	// Load environment variables from .env file
	_ "github.com/joho/godotenv/autoload"
)

var Port = os.Getenv("PORT")

var (
	DatabaseURL   = os.Getenv("DATABASE_URL")
	MigrationPath = os.Getenv("MIGRATION_PATH")
)

var (
	KakaoRestAPIKey  = os.Getenv("KAKAO_REST_API_KEY")
	KakaoRedirectURI = os.Getenv("KAKAO_REDIRECT_URI")
)

var FirebaseCredentialsPath = os.Getenv("FIREBASE_CREDENTIALS_PATH")

type FirebaseCredentialsJSONType struct {
	Type                    string `json:"type"`
	ProjectID               string `json:"project_id"`
	PrivateKeyID            string `json:"private_key_id"`
	PrivateKey              string `json:"private_key"`
	ClientEmail             string `json:"client_email"`
	ClientID                string `json:"client_id"`
	AuthURI                 string `json:"auth_uri"`
	TokenURI                string `json:"token_uri"`
	AuthProviderX509CertURL string `json:"auth_provider_x509_cert_url"`
	ClientX509CertURL       string `json:"client_x509_cert_url"`
	UniverseDomain          string `json:"universe_domain"`
}

func GetFirebaseCredentialsJSON() FirebaseCredentialsJSONType {
	return FirebaseCredentialsJSONType{
		Type:                    os.Getenv("FIREBASE_CREDENTIALS_TYPE"),
		ProjectID:               os.Getenv("FIREBASE_CREDENTIALS_PROJECT_ID"),
		PrivateKeyID:            os.Getenv("FIREBASE_CREDENTIALS_PRIVATE_KEY_ID"),
		PrivateKey:              strings.ReplaceAll(os.Getenv("FIREBASE_CREDENTIALS_PRIVATE_KEY"), "\\n", "\n"),
		ClientEmail:             os.Getenv("FIREBASE_CREDENTIALS_CLIENT_EMAIL"),
		ClientID:                os.Getenv("FIREBASE_CREDENTIALS_CLIENT_ID"),
		AuthURI:                 os.Getenv("FIREBASE_CREDENTIALS_AUTH_URI"),
		TokenURI:                os.Getenv("FIREBASE_CREDENTIALS_TOKEN_URI"),
		AuthProviderX509CertURL: os.Getenv("FIREBASE_CREDENTIALS_AUTH_PROVIDER_X509_CERT_URL"),
		ClientX509CertURL:       os.Getenv("FIREBASE_CREDENTIALS_CLIENT_X509_CERT_URL"),
		UniverseDomain:          os.Getenv("FIREBASE_CREDENTIALS_UNIVERSE_DOMAIN"),
	}
}

var (
	B2KeyID      = os.Getenv("B2_APPLICATION_KEY_ID")
	B2Key        = os.Getenv("B2_APPLICATION_KEY")
	B2BucketName = os.Getenv("B2_BUCKET_NAME")
	B2Endpoint   = os.Getenv("B2_ENDPOINT")
	B2Region     = os.Getenv("B2_REGION")
)

var (
	GoogleSheetsAPIKey   = os.Getenv("GOOGLE_SHEETS_API_KEY")
	BreedsGoogleSheetsID = os.Getenv("BREEDS_GOOGLE_SHEETS_ID")
)

//nolint:gochecknoinits
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

	if FirebaseCredentialsPath == "" && GetFirebaseCredentialsJSON() == (FirebaseCredentialsJSONType{}) {
		panic("FIREBASE_CREDENTIALS_PATH or FIREBASE_CREDENTIALS_JSON is required")
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
