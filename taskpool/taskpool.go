package taskpool

import (
	"context"
	"log"
)

type task func() error

// TaskPool manages concurrent processes. No more than size tasks run at once,
// and when a task returns an error, the Pool's context is cancelled.
type TaskPool struct {
	add  chan task
	wait chan error
}

// New returns a TaskPool waiting for tasks to run.
func New(ctx context.Context, size int) (TaskPool, context.Context) {
	ctx, cancel := context.WithCancel(ctx)

	t := TaskPool{
		add:  make(chan task),
		wait: make(chan error),
	}
	go t.start(ctx, cancel, size)
	return t, ctx
}

func (tp TaskPool) start(ctx context.Context, cancel context.CancelFunc, size int) {
	result := make(chan error)
	workerC := make(chan task)

	defer cancel()
	defer close(workerC) // Make sure worker goroutines get cleaned up
	defer close(tp.add)  // Closing this is pointless unless there's a bug

	// Start worker pool
	for i := 0; i < size; i++ {
		go func() {
			for t := range workerC {
				select {
				case <-ctx.Done():
					log.Println("37")
					return
				case result <- t():
				}
			}
		}()
	}

	var (
		nextTask task
		tasks    []task
		taskC    chan task
		waitC    chan error
		running  int
	)

	for {
		select {
		case <-ctx.Done():
			tp.wait <- ctx.Err()
			log.Println("58")
			return

		case t := <-tp.add:
			tasks = append(tasks, t)

		case taskC <- nextTask:
			tasks = tasks[1:]
			running++

		case waitC <- nil:
			return

		case err := <-result:
			if err != nil {
				tp.wait <- err
				log.Println("72")
				return
			}
			running--
		}

		waitC = nil
		if len(tasks) > 0 {
			nextTask = tasks[0]
			taskC = workerC
		} else {
			if running == 0 {
				waitC = tp.wait
			}
			nextTask = nil
			taskC = nil
		}
	}
}

// Go adds a task to the pool of tasks to complete.
func (tp TaskPool) Go(t func() error) {
	tp.add <- t
}

// Wait waits for all tasks added to the pool to finish running.
func (tp TaskPool) Wait() error {
	return <-tp.wait
}
