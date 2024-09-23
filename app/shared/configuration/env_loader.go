package configuration

import (
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"strings"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/joho/godotenv"
)

type EnvLoader struct {
}

func init() {
	ioc.Registry(NewEnvLoader)
}
func NewEnvLoader() EnvLoader {
	if err := godotenv.Load(); err != nil {
		slog.Warn(".env not found, loading environment from system.")
	}
	return EnvLoader{}
}

func (env EnvLoader) Get(key string) string {
	return os.Getenv(key)
}

func validateConfig[T any](conf T) (T, error) {
	var validationErrors []error

	val := reflect.ValueOf(conf)
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		value := val.Field(i).String()

		requiredTag := field.Tag.Get("required")
		if requiredTag == "true" && value == "" {
			validationErrors = append(validationErrors, fmt.Errorf("%s is required but not set", field.Name))
		}
	}
	if len(validationErrors) > 0 {
		// Convert errors to strings
		var errorStrings []string
		for _, err := range validationErrors {
			errorStrings = append(errorStrings, err.Error())
		}
		return conf, fmt.Errorf("configuration errors:\n%s", strings.Join(errorStrings, "\n"))
	}
	return conf, nil
}
