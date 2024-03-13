package utils

import (
	"errors"
	"regexp"
)

func ExtractAuthData(connectionString string) (string, string, error) {
	//Expects an <ArgusAuth>Username:Password</ArgusAuth> string

	pattern := `<ArgusAuth>(.+):(.+)</ArgusAuth>`

	regex := regexp.MustCompile(pattern)
	matches := regex.FindStringSubmatch(connectionString)

	if len(matches) == 3 {
		username := matches[1]
		password := matches[2]

		return username, password, nil
	} else {
		return "", "", errors.New("could not extract auth data")
	}
}
