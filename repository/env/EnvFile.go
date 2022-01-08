package env

import (
	"os"
)

func GetDotEnvVariable(key string) string {
	//err := godotenv.Load(".env")

	//if err != nil {
	//	log.Fatalf("Error loading .env file")
	//}
	return os.Getenv(key)
}