package helper

import "os"

/*
The GetEnv function returns the value of an OS environment variable.
If the environment variable does not exist, the default value is returned.

Input
  - variableName: The OS environment variable name.
  - defaultValue: A value to return if the variableName does not exist in the environment.

Output
  - The actual or defaulted value of the OS environment variable.
*/
func GetEnv(variableName string, defaultValue string) string {
	osenvValue := os.Getenv(variableName)
	if len(osenvValue) > 0 {
		return osenvValue
	}
	return defaultValue
}
