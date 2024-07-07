package env

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func Init() {
	godotenv.Load() //nolint
}

func InitWithMandatoryVariables(mandatoryVariables []string) {
	godotenv.Load() //nolint
	validateEnvVariables(mandatoryVariables)
}

func validateEnvVariables(envVars []string) {
	for _, value := range envVars {
		_, defined := os.LookupEnv(value)
		if !defined {
			panic(envErrorMessage(value))
		}
	}
}

func envErrorMessage(value string) string {
	return fmt.Sprintf("ENV Variable is not defined: %s", value)
}
