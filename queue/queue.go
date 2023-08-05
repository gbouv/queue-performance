package queue

import (
	"context"

	"github.com/gbouv/queue-performance/queue/model"
)

type Queue interface {
	Insert(ctx context.Context, jobToBeQueued *model.QueuedJob) error

	FetchTentatively(ctx context.Context) (*model.QueuedJob, bool, error)

	Remove(ctx context.Context, jobId model.JobId) error

	Size(ctx context.Context) (int, error)
}
