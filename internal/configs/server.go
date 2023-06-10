package configs

import "os"

var Port = os.Getenv("PORT")

func init() {
	if Port == "" {
		Port = "8080"
	}
}
