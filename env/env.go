package env

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

func Init() {
	godotenv.Load() //nolint
}

func InitWithMandatoryVariables(mandatoryVariables []string) {
	godotenv.Load() //nolint
	validateEnvVariables(mandatoryVariables)
}

func Env(v string) string {
	return os.Getenv(v)
}

func EnvInt(v string) int {
	val, err := strconv.Atoi(os.Getenv(v))
	if err != nil {
		panic(err)
	}
	return val
}

func EnvTimeDuration(v string) time.Duration {
	return time.Duration(EnvInt(v))
}

func EnvBool(v string) bool {
	val, err := strconv.ParseBool(os.Getenv(v))
	if err != nil {
		panic(err)
	}
	return val
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
