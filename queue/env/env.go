package env

import (
	"math"
	"os"
	"strconv"

	"github.com/palantir/stacktrace"
	"github.com/sirupsen/logrus"
)

func GetParamsFromEnv() (*Params, error) {
	logLevelStr := getEnvOrDefault("LOG_LEVEL", "INFO")
	logLevel, err := logrus.ParseLevel(logLevelStr)
	if err != nil {
		return nil, stacktrace.Propagate(err, "Error parsing log level '%s'", logLevelStr)
	}

	dbHost := getEnvOrDefault("DB_HOST", "localhost")

	dbPortStr := getEnvOrDefault("DB_PORT", "5432")
	dbPort, err := toUint16(dbPortStr)
	if err != nil {
		return nil, stacktrace.Propagate(err, "Error parsing DB port parameter")
	}

	dbName := getEnvOrDefault("DB_NAME", "postgres")

	dbUser := getEnvOrDefault("DB_USER", "postgres")

	dbPassword := getEnvOrDefault("DB_PASSWORD", "password")

	redisUrl := getEnvOrDefault("REDIS_URL", "")

	return &Params{
		LogLevel: logLevel,

		DbHost:     dbHost,
		DbPort:     dbPort,
		DbName:     dbName,
		DbUser:     dbUser,
		DbPassword: dbPassword,

		RedisUrl: redisUrl,
	}, nil
}

func getEnvOrDefault(envName string, defaultValue string) string {
	if value, found := os.LookupEnv(envName); found {
		return value
	}
	return defaultValue
}

func toUint16(valueStr string) (uint16, error) {
	valueUint64, err := strconv.ParseUint(valueStr, 10, 16)
	if err != nil {
		return 0, stacktrace.Propagate(err, "Unable to convert value '%s' to uint16", valueStr)
	}
	if valueUint64 > math.MaxUint16 || valueUint64 < 0 {
		return 0, stacktrace.Propagate(err, "Unable to convert value '%s' to uint16, it is out of range", valueStr)
	}
	return uint16(valueUint64), nil
}
