package helper

import "os"

func GetEnv(variableName string, defaultValue string) string {
	osenvValue := os.Getenv(variableName)
	if len(osenvValue) > 0 {
		return osenvValue
	}
	return defaultValue
}
