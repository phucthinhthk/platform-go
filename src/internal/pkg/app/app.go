package app

import "os"

// IsLocal determines if the application is running in a local environment
func IsLocal() bool {
	env := os.Getenv("APP_ENV")
	return env == "" || env == "local"
}
