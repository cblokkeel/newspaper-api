package workers

import (
	"context"
	"fmt"

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
	fmt.Println("ahh okay")
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
