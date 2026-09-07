// Package env loads environment variables from mode-specific .env files.
package env

import (
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// LoadEnv loads a mode-specific dotenv file into the current process
// environment.
//
// The supported modes are "dev" for .env.dev, "prod" for .env, and "test" for
// .env.test. Unknown modes default to .env.dev. LoadEnv checks the current
// directory first and then its parent. Missing files are not treated as errors.
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
