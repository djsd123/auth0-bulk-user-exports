package default_env_value

import "os"

func SetDefaultEnvValue(key, defaultVal string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = defaultVal
	}
	return value
}
