package config

import (
	"errors"
	"io/ioutil"
)

//LoadToken loads the bot token from given file path
func LoadToken(path string) (string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", errors.New("Token file not found")
	}
	return string(data), nil
}
