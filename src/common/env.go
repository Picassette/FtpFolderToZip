package common

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

/*
* Get Environement variable as String
 */
func GetEnvVarAsStr(name string, required bool) (*string, error) {
	envVarValue := os.Getenv(name)
	if envVarValue == "" {
		if required {
			return nil, fmt.Errorf("required value %s is not set", name)
		}
		return new(string), nil
	} else {
		return &envVarValue, nil
	}
}

/*
* Get Environement variable as Number
 */
func GetEnvVarAsNumber(name string, required bool) (*int64, error) {
	val, err := GetEnvVarAsStr(name, required)
	if err != nil {
		return nil, err
	}
	number, err := strconv.ParseInt(*val, 10, 16)
	if err != nil {
		return nil, fmt.Errorf("error during env var number conversion for var %s : %s", name, err.Error())
	}
	return &number, nil
}

/*
* Get Environement variable as String
* Return the value with compatible filesystem path
 */
func GetEnvVarAsStrPath(name string, required bool) (*string, error) {
	val, err := GetEnvVarAsStr(name, required)
	if err != nil {
		return nil, err
	}
	*val = filepath.FromSlash(*val)
	return val, nil
}
