package workers

import (
	"context"
	"fmt"

	redisdb "github.com/cblokkeel/newspaper/internal/db/redis"
	"github.com/robfig/cron"
	"go.uber.org/fx"
)

type Worker interface {
	Start(context.Context)
	Periodicity() string
}

type WorkerManager struct {
	Workers []Worker
}

func NewWorkerManager(workers []Worker) *WorkerManager {
	return &WorkerManager{
		Workers: workers,
	}
}

func StartWorkers(lc fx.Lifecycle, wm *WorkerManager) *cron.Cron {
	cron := cron.New()
	for _, worker := range wm.Workers {
		cron.AddFunc(worker.Periodicity(), func() {
			worker.Start(context.Background())
		})
	}
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			fmt.Println("Starting workers")
			go cron.Start()
			return nil
		},
		OnStop: func(context.Context) error {
			fmt.Println("Ending workers")
			cron.Stop()
			return nil
		},
	})
	return cron
}

type CommonWorker struct {
	name  string
	redis *redisdb.RedisDB
}

func (w *CommonWorker) Lock() error {
	if err := w.redis.RDB.Set(context.Background(), fmt.Sprintf("%s:lock", w.name), true, 0).Err(); err != nil {
		return err
	}
	return nil
}

func (w *CommonWorker) Unlock() error {
	if err := w.redis.RDB.Del(context.Background(), fmt.Sprintf("%s:lock", w.name)).Err(); err != nil {
		return err
	}
	return nil
}

func (w *CommonWorker) IsLocked() bool {
	if err := w.redis.RDB.Get(context.Background(), fmt.Sprintf("%s:lock", w.name)).Err(); err == nil {
		return true
	}
	return false
}
