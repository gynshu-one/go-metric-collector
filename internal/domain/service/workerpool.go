package service

import (
	"github.com/rs/zerolog/log"
	"sync"
)

// WorkerPool is a pool of workers that can be used to execute tasks concurrently.
// Modified by me version of Practicum example
// The pool is created with a fixed number of workers, and a queue of tasks.
// When a task is pushed to the pool, it is added to the queue.
// A worker picks up the task from the queue and executes it.
// When the task is done, the worker picks up the next task from the queue or waits.
// Usage:
//
//		wrk := workerpool.New()
//		wrk.Push(&workerpool.Task{
//			ID: "someID",
//	     Task: func()  {
//						log.Info().Msgf("start processing %s", f)
//						time.Sleep(1 * time.Second)
//						log.Info().Msgf("end processing %s", f)
//					},
type WorkerPool interface {
	Push(t *Task)
	Stop()
}
type Task struct {
	ID   string
	Task func()
}

type queue struct {
	ch       chan *Task
	wg       *sync.WaitGroup
	wrkCount int
}

func newQueue() *queue {
	return &queue{
		ch: make(chan *Task, 1),
		wg: &sync.WaitGroup{},
	}
}

func (q *queue) Push(t *Task) {
	defer func() {
		if r := recover(); r != nil {
			log.Error().Msgf("workerpool is stopped")
		}
	}()
	q.ch <- t
	q.wg.Add(1)
}

type worker struct {
	id    int
	queue *queue
}

func newWorker(id int, queue *queue) *worker {
	w := worker{
		id:    id,
		queue: queue,
	}
	return &w
}

func (w *worker) loop() {
	for {
		t, ok := <-w.queue.ch
		if !ok {
			log.Debug().Msgf("worker %d is stopping", w.id)
			w.queue.wg.Done()
			return
		}
		log.Debug().Msgf("worker %d got Task %s", w.id, t.ID)
		t.Task()
		log.Debug().Msgf("worker %d finished Task %s", w.id, t.ID)
		w.queue.wg.Done()
	}
}

// NewWorkerPool creates a new workerpool instance with the number of workers equal to the number of CPUs minus 2.
// You can push tasks to the workerpool it will execute them in parallel.
// You can stop the workerpool by calling Stop() method.
func NewWorkerPool(howMany int) *queue {
	q := newQueue()
	workers := make([]*worker, 0, howMany)
	q.wrkCount = howMany
	for i := 0; i < howMany; i++ {
		workers = append(workers, newWorker(i, q))
	}
	wg := &sync.WaitGroup{}
	wg.Add(len(workers))
	for _, w := range workers {
		go w.loop()
		wg.Done()
	}
	wg.Wait()
	return q
}

// Stop stops the workerpool Softly. It will dismiss all new tasks and return error.
func (q *queue) Stop() {
	// confirm all workers are done
	q.wg.Add(q.wrkCount)
	close(q.ch)
	q.wg.Wait()
}
