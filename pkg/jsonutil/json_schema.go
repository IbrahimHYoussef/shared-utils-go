package jsonutil

import (
	"os"
	"strings"
)

// LoadSchema reads a JSON schema file and returns its contents as a string.
func LoadSchema(filePath string) (string, error) {
	logger.Info("loading", "file_path", filePath)
	date, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(date), nil
}

// LoadSchemas reads all .json files in dirPath and returns their contents keyed
// by file name.
func LoadSchemas(dirPath string) (map[string]string, error) {
	if dirPath[len(dirPath)-1] != '/' {
		dirPath = dirPath + "/"
	}
	logger.Info("loading json schemas", "path", dirPath)
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}
	result := make(map[string]string)
	for _, file := range files {
		sliced := strings.Split(file.Name(), ".")
		if sliced[len(sliced)-1] == "json" {
			fcontent, err := LoadSchema(dirPath + file.Name())
			if err != nil {
				return nil, err
			}
			result[file.Name()] = fcontent
		}
	}
	return result, nil
}
