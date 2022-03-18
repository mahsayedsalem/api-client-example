package form3_api_client

import "os"

var BaseUrlEnvVariable string = getEnv("Form3_API_URL", "http://localhost:8080/")

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}
