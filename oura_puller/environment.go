package main

import (
	"errors"
	"fmt"
	"os"
)

func GetEnvString(key string) (string, error) {
	returnValue := os.Getenv(key)

	if returnValue == "" {
		return "", errors.New(fmt.Sprintf("Environment value '%s' does not exist", key))
	}

	return returnValue, nil
}
