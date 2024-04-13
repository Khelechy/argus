package utils

import (
	"errors"
	"path/filepath"
	"regexp"
	"strings"
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

func TreatAsWildcard(pathWithWildcard string) (bool, string, string) {
	folderPath, fileName := filepath.Split(pathWithWildcard)

	dotIndex := strings.LastIndex(fileName, ".")

	if dotIndex != -1 && dotIndex < len(fileName)-1 {
		folderPath = strings.TrimSuffix(folderPath, "/")

		extension := filepath.Ext(fileName)

		if !strings.ContainsRune(fileName, '*') {
			return false, folderPath, extension
		}

		extension = strings.TrimPrefix(extension, "*")

		return true, folderPath, extension
	}

	return false, folderPath, "" // Return empty file extension if its an absolute folder
}
