package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/gbouv/queue-performance/queue"
	"github.com/gbouv/queue-performance/queue/env"
	"github.com/gbouv/queue-performance/queue/hybrid"
	"github.com/gbouv/queue-performance/queue/model"
	"github.com/gbouv/queue-performance/queue/postgres"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

const (
	maxQueueSize               = 1_000_000             // when the queue contains this number of rows, stop filling it
	pauseDurationWhenQueueFull = 10 * time.Millisecond // when the queue is full, duration to wait for before checking again
	insertBatchSize            = 100                   // the size of the queue will be checked after each batch is inserted
)

func main() {
	ctx := context.Background()
	params, err := env.GetParamsFromEnv()
	if err != nil {
		exitError(err)
	}
	logrus.SetLevel(params.LogLevel)

	logrus.Info("Starting orchestrator")

	database, err := postgres.Connect(params.DbHost, params.DbPort, params.DbName, params.DbUser, params.DbPassword)
	if err != nil {
		exitError(err)
	}

	var queue queue.Queue
	if params.RedisUrl == "" {
		queue = postgres.NewPostgresQueue(database)
	} else {
		parsedRedisUrl, err := redis.ParseURL(params.RedisUrl)
		if err != nil {
			exitError(err)
		}
		redisClient := redis.NewClient(parsedRedisUrl)
		queue = hybrid.NewHybridQueue(database, redisClient)
	}

	runIndefinitely(ctx, queue)
}

func runIndefinitely(ctx context.Context, queue queue.Queue) {
	insertionErrors := 0
	insertedJobs := 0
	defer func() {
		if insertionErrors > 0 {
			exitError(fmt.Errorf("inserted %d jobs but %d insertion errors occurred", insertedJobs, insertionErrors))
		}
		logrus.Infof("Successfully inserted %d jobs", insertedJobs)
	}()

	lastQueueSizePrint := time.Now()
	for {
		numberOfQueuedJobs, err := queue.Size(ctx)
		if err != nil {
			exitError(err)
		}

		if time.Since(lastQueueSizePrint) > 1*time.Second {
			logrus.Infof("Queue size:|%6d|", numberOfQueuedJobs)
			lastQueueSizePrint = time.Now()
		}

		if numberOfQueuedJobs >= maxQueueSize {
			logrus.Debugf("Too many queued jobs. Pausing for %v", pauseDurationWhenQueueFull)
			time.Sleep(pauseDurationWhenQueueFull)
			continue
		}

		for i := 0; i <= insertBatchSize; i++ {
			randomJob := generateRandomJob()
			logrus.Debugf("Inserting job with ID %s", randomJob.JobId)
			if err := queue.Insert(ctx, randomJob); err != nil {
				logrus.Errorf("An error occurred inserting job with ID %s\n%s", randomJob.JobId, err.Error())
				insertionErrors += 1
			} else {
				insertedJobs += 1
			}
		}
	}
}

func exitError(err error) {
	logrus.Error(err.Error())
	logrus.Exit(1)
}

func generateRandomJob() *model.QueuedJob {
	jobId := uuid.New()
	jobDifficulty := rand.Intn(10_000)
	return &model.QueuedJob{
		JobId:       model.JobId(jobId.String()),
		CreatedTime: time.Now(),
		StartedTime: nil, // Will be filled by the queue when the job is fetched
		Difficulty:  jobDifficulty,
	}
}
