package env

import "github.com/sirupsen/logrus"

type Params struct {
	LogLevel logrus.Level

	DbHost     string
	DbPort     uint16
	DbName     string
	DbUser     string
	DbPassword string

	RedisUrl string
}
