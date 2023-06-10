package configs

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

var Port = os.Getenv("PORT")

var KakaoRestAPIKey = os.Getenv("KAKAO_REST_API_KEY")
var KakaoRedirectURI = os.Getenv("KAKAO_REDIRECT_URI")

var FirebaseCredentialsPath = os.Getenv("FIREBASE_CREDENTIALS_PATH")

func init() {
	if Port == "" {
		Port = "8080"
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
