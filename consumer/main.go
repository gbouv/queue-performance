package main

import (
	"context"
	"time"

	"github.com/gbouv/queue-performance/queue"
	"github.com/gbouv/queue-performance/queue/env"
	"github.com/gbouv/queue-performance/queue/hybrid"
	"github.com/gbouv/queue-performance/queue/model"
	"github.com/gbouv/queue-performance/queue/postgres"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

const (
	queueFullyConsumedMaxRetry      = 1000
	queueFullyConsumerRetryInterval = 10 * time.Millisecond
)

func main() {
	ctx := context.Background()
	params, err := env.GetParamsFromEnv()
	if err != nil {
		exitError(err)
	}
	logrus.SetLevel(params.LogLevel)

	logrus.Infof("Starting worker")

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
	queueFullyConsumedRetry := 0
	timeStart := time.Now()
	nJobsExecuted := 0
	nJobsExecutedLastPrint := time.Now()
	for {
		queuedJob, found, err := queue.FetchTentatively(ctx)
		if err != nil {
			exitError(err)
		}
		if !found {
			queueFullyConsumedRetry += 1
			if queueFullyConsumedRetry == queueFullyConsumedMaxRetry {
				logrus.Infof("No more job to fetch from the queue after %d retries. Exiting", queueFullyConsumedMaxRetry)
				return
			} else {
				logrus.Infof("No more job to fetch from the queue. Will retry in %v (%d/%d)",
					queueFullyConsumerRetryInterval, queueFullyConsumedRetry, queueFullyConsumedMaxRetry)
			}
			time.Sleep(queueFullyConsumerRetryInterval)
			continue
		}
		queueFullyConsumedRetry = 0 // reset this number in case it has retried a few times

		logrus.Debugf("Fetched job with ID %s (difficulty %d)", queuedJob.JobId, queuedJob.Difficulty)
		if err := executeJob(queuedJob); err != nil {
			exitError(err)
		}

		logrus.Debugf("Removing job %s from queue", queuedJob.JobId)
		if err := queue.Remove(ctx, queuedJob.JobId); err != nil {
			exitError(err)
		}
		logrus.Debugf("Successfully executed job with ID %s", queuedJob.JobId)

		nJobsExecuted += 1
		if time.Since(nJobsExecutedLastPrint) > 1*time.Second {
			logrus.Infof("Jobs Per Second:|%6d|", nJobsExecuted/int(time.Since(timeStart).Seconds()))
			nJobsExecutedLastPrint = time.Now()
		}
	}
}

func exitError(err error) {
	logrus.Error(err.Error())
	logrus.Exit(1)
}

func executeJob(job *model.QueuedJob) error {
	// To test performance job is a no-op, difficulty is not taken into account
	// For now difficulty is just the duration of the job in millisecond
	// time.Sleep(time.Duration(job.Difficulty * 1_000))
	return nil
}
