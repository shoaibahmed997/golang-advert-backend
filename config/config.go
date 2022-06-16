package config

import (
	"os"
)

func Config(key string) string {
	// err := godotenv.Load(".env")
	// if err != nil {
	// 	log.Println("env load error", err)
	// }
	return os.Getenv(key)

}
