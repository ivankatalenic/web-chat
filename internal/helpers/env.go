package helpers

import "os"

// GetEnvVarOrDefault returns a value of the environment variable with the provided name.
// If the environment variable is not set, then it returns a default string.
func GetEnvVarOrDefault(envName, def string) string {
	val := os.Getenv(envName)
	if len(val) == 0 {
		return def
	}
	return val
}
