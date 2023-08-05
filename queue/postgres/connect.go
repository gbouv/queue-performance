package postgres

import (
	"fmt"

	"github.com/gbouv/queue-performance/queue/model"
	"github.com/palantir/stacktrace"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(hostame string, port uint16, databaseName string, username string, password string) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d",
		hostame,
		username,
		password,
		databaseName,
		port,
	)
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return nil, stacktrace.Propagate(err, "Error opening postgres database")
	}

	initSchema(database)
	return database, nil
}

func initSchema(database *gorm.DB) {
	database.AutoMigrate(&model.FinishedJob{})
	database.AutoMigrate(&model.QueuedJob{})
}
