package env

import (
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func LoadEnv(mode string) error {
	var fileName string
	switch mode {
	case "dev":
		fileName = ".env.dev"
	case "prod":
		fileName = ".env"
	case "test":
		fileName = ".env.test"
	default:
		fileName = ".env.dev"
	}

	for _, filePath := range []string{fileName, filepath.Join("..", fileName)} {
		if _, err := os.Stat(filePath); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return err
		}

		if err := godotenv.Load(filePath); err != nil {
			return err
		}

		log.Printf("Loaded env from %s", filePath)
		return nil
	}

	wd, _ := os.Getwd()
	log.Printf("No %s file found in current dir or parent. Working directory: %s. Continuing with existing environment.", fileName, wd)
	return nil
}

