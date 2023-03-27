package utils

import (
	"os"
)

func GetEnvDefault(key, defVal string) string {
	val, exist := os.LookupEnv(key)
	if !exist {
		return defVal
	}
	return val
}
