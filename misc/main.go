package misc

import "os"

func GetEnv(varName, defaultValue string) string {
	value, exists := os.LookupEnv(varName)
	if exists {
		return value
	}
	return defaultValue
}
