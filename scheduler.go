package stix

import (
	"context"
	"log"
	"time"

	"github.com/su5kk/stix/pkg/container"
)

type Task struct {
	ID          string
	Name        string
	Description string
	Timeout     time.Duration
	Execute     func() error
}

type Scheduler struct {
	tasks     *container.SyncMap[string, *Task]
	addpipe   chan *Task
	retrypipe chan string
}

func NewScheduler(addpipe chan *Task, retrypipe chan string) *Scheduler {
	return &Scheduler{
		addpipe:   addpipe,
		retrypipe: retrypipe,
		tasks:     container.NewSyncMap[string, *Task](),
	}
}

// Start starts the Scheduler and returns Shutdown func
func (s *Scheduler) Start(ctx context.Context) func() {
	return s.runScheduler(ctx)
}

func (s *Scheduler) runScheduler(ctx context.Context) func() {
	ctx, cancel := context.WithCancel(ctx)
	done := make(chan struct{})

	go func() {
		defer close(done)
		s.iterLoop(ctx)
	}()

	return func() {
		cancel()
		<-done
	}
}

func (s *Scheduler) iterLoop(ctx context.Context) {
	iter := func(ctx context.Context, run func(task *Task)) bool {
		select {
		case <-ctx.Done():
			// TODO: add graceful shutdown
			log.Println("shutting down scheduler")
			return false
		case t := <-s.addpipe:
			log.Println("received new task through pipe", t.ID)
			s.tasks.Store(t.ID, t)
			run(t)
		case tID := <-s.retrypipe:
			log.Println("retrying task with id", tID)
			if t, ok := s.tasks.Load(tID); ok {
				run(t)
			} else {
				log.Println("task with id not found", tID)
			}
		}
		return true
	}

	for iter(ctx, func(task *Task) {
		err := task.Execute()
		if err != nil {
			log.Println("task execution failed")
			return
		}
		s.tasks.Delete(task.ID)
	}) {
	}
}
