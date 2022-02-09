package env

import (
	"os"
)

func GetDotEnvVariable(key string) string {
	return os.Getenv(key)
}
