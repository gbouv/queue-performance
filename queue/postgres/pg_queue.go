package postgres

import (
	"context"
	"time"

	"github.com/gbouv/queue-performance/queue"
	"github.com/gbouv/queue-performance/queue/model"
	"github.com/palantir/stacktrace"
	"gorm.io/gorm"
)

type PostgresQueue struct {
	database *gorm.DB
}

func NewPostgresQueue(database *gorm.DB) queue.Queue {
	return &PostgresQueue{
		database: database,
	}
}

func (queue *PostgresQueue) Insert(ctx context.Context, jobToBeQueued *model.QueuedJob) error {
	err := queue.database.Transaction(func(tx *gorm.DB) error {
		result := tx.Exec(
			"INSERT INTO queued_jobs (job_id, created_time, difficulty) VALUES (?, ?, ?)",
			jobToBeQueued.JobId,
			jobToBeQueued.CreatedTime,
			jobToBeQueued.Difficulty,
		)
		if result.Error != nil {
			return stacktrace.Propagate(result.Error, "An error occurred inserting queued job with ID %s", jobToBeQueued.JobId)
		}
		return nil
	})
	if err != nil {
		return stacktrace.Propagate(err, "An error occurred inserting queued job with ID %s", jobToBeQueued.JobId)
	}
	return nil
}

func (queue *PostgresQueue) FetchTentatively(ctx context.Context) (*model.QueuedJob, bool, error) {
	jobFound := false
	queuedJob := new(model.QueuedJob)
	err := queue.database.Transaction(func(tx *gorm.DB) error {
		result := tx.Raw(
			`UPDATE queued_jobs SET started_time = ? WHERE job_id = (
                SELECT job_id 
                FROM queued_jobs inner_queued_jobs
                WHERE inner_queued_jobs.started_time IS NULL
                ORDER BY inner_queued_jobs.created_time ASC 
                LIMIT 1
                FOR UPDATE SKIP LOCKED
            )
			RETURNING job_id, created_time, started_time, difficulty`,
			time.Now(),
		).Scan(queuedJob)
		if result.Error != nil {
			return stacktrace.Propagate(result.Error, "An error occurred updating queued job record")
		}
		if result.RowsAffected > 0 {
			jobFound = true
		}
		return nil
	})
	if err != nil {
		return nil, false, stacktrace.Propagate(err, "An error occurred fetching a queued job")
	}
	if jobFound {
		return queuedJob, true, nil
	}
	return nil, false, nil
}

func (queue *PostgresQueue) Remove(ctx context.Context, jobId model.JobId) error {
	err := queue.database.Transaction(func(tx *gorm.DB) error {
		// First we deleted the job from the queue
		deletedJob := new(model.QueuedJob)
		deleteResult := tx.Raw(
			"DELETE FROM queued_jobs WHERE job_id = ? RETURNING job_id, created_time, started_time, difficulty",
			jobId,
		).Scan(deletedJob)
		if deleteResult.Error != nil {
			return stacktrace.Propagate(deleteResult.Error, "An error occurred updating queued job record")
		}
		if deleteResult.RowsAffected == 0 {
			return stacktrace.NewError("An error occurred removing job from queued_job table")
		}

		// Then we insert it into the finished_job table
		insertResult := tx.Exec(
			"INSERT INTO finished_jobs (job_id, started_time, duration_ms) VALUES (?, ?, ?)",
			deletedJob.JobId,
			deletedJob.StartedTime,
			time.Since(*deletedJob.StartedTime),
		)
		if insertResult.Error != nil {
			return stacktrace.Propagate(insertResult.Error, "An error occurred updating queued job record")
		}
		if insertResult.RowsAffected == 0 {
			return stacktrace.NewError("An error occurred persisting job to finished_job table")
		}
		return nil
	})
	if err != nil {
		return stacktrace.Propagate(err, "An error occurred persisting job as finished")
	}
	return nil
}

func (queue *PostgresQueue) Size(ctx context.Context) (int, error) {
	var numberOfQueuedJobs int
	err := queue.database.Transaction(func(tx *gorm.DB) error {
		result := tx.Raw(
			"SELECT COUNT(*) FROM queued_jobs",
		).Scan(&numberOfQueuedJobs)
		if result.Error != nil {
			return stacktrace.Propagate(result.Error, "An error occurred counting the number of queued jobs")
		}
		return nil
	})
	if err != nil {
		return 0, stacktrace.Propagate(err, "An error occurred counting the number of queued jobs")
	}
	return numberOfQueuedJobs, nil
}
