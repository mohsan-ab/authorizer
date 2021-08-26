package async

import (
	"github.com/mohsanabbas/authorizer/commons/eventsource"
)

var workerPool = make(chan chan eventsource.Command)

// Worker contains the basic info to manage commands
type Worker struct {
	WorkerPool     chan chan eventsource.Command
	JobChannel     chan eventsource.Command
	CommandHandler eventsource.CommandHandlerRegister
}

// Bus stores the command handler
type Bus struct {
	CommandHandler eventsource.CommandHandlerRegister
	maxWorkers     int
}

// Start initialize a worker ready to receive jobs
func (w *Worker) Start() {
	go func() {
		for {
			w.WorkerPool <- w.JobChannel

			job := <-w.JobChannel
			handler, err := w.CommandHandler.Get(job)
			if err != nil {
				continue
			}

			if !job.IsValid() {
				continue
			}

			if err = handler.Handle(job); err != nil {
				//TODO: log the error

			}
		}
	}()
}

// NewWorker initialize the values of worker and start it
func NewWorker(commandHandler eventsource.CommandHandlerRegister) {
	w := Worker{
		WorkerPool:     workerPool,
		CommandHandler: commandHandler,
		JobChannel:     make(chan eventsource.Command),
	}

	w.Start()
}

// HandleCommand ad a job to the queue
func (b *Bus) HandleCommand(command eventsource.Command) {
	go func(c eventsource.Command) {
		workerJobQueue := <-workerPool
		workerJobQueue <- c
	}(command)
}

// NewBus return a bus with command handler register
func NewBus(register eventsource.CommandHandlerRegister, maxWorkers int) *Bus {
	b := &Bus{
		CommandHandler: register,
		maxWorkers:     maxWorkers,
	}

	// start the bus
	b.Start()
	return b
}

// Start the bus
func (b *Bus) Start() {
	for i := 0; i < b.maxWorkers; i++ {
		NewWorker(b.CommandHandler)
	}
}
