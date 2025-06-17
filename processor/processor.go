package processor

import (
	"context"
	"log"
	"sync"
	"time"

	"playground/queue"
)

type Processor[T any] struct {
	stopChan    chan struct{}
	wg          sync.WaitGroup
	id          int
	workers     int
	queue       *queue.Queue[T]
	mu          sync.RWMutex
	processFunc func(T) error
}

func NewProcessor[T any](id int, workers int, processFunc func(T) error) *Processor[T] {
	return &Processor[T]{
		stopChan:    make(chan struct{}),
		id:          id,
		workers:     workers,
		queue:       queue.NewQueue[T](),
		processFunc: processFunc,
	}
}

// Enqueue adds a message to the processor's queue
func (p *Processor[T]) Enqueue(msg T) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.queue.Enqueue(msg)
}

// GetQueueSize returns the current size of the processor's queue
func (p *Processor[T]) GetQueueSize() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.queue.Size()
}

func (p *Processor[T]) Start(ctx context.Context) {
	// Start multiple workers
	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		workerID := i + 1
		go p.worker(ctx, workerID)
	}
}

func (p *Processor[T]) worker(ctx context.Context, workerID int) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic in Processor %d, Worker %d: %v", p.id, workerID, r)
		}
		p.wg.Done()
	}()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-p.stopChan:
			return
		case <-ticker.C:
			// Try to get a message from the queue
			var message T
			var ok bool

			func() {
				p.mu.Lock()
				defer p.mu.Unlock()
				message, ok = p.queue.Dequeue()
			}()

			if ok {
				// Process the message
				log.Printf("Processor %d, Worker %d processing message", p.id, workerID)
				// Use the provided processing function
				if err := p.processFunc(message); err != nil {
					log.Printf("Error processing message in Processor %d, Worker %d: %v", p.id, workerID, err)
				}
				log.Printf("Processor %d, Worker %d completed processing message", p.id, workerID)
			}
		}
	}
}

func (p *Processor[T]) Stop() {
	close(p.stopChan)
	p.wg.Wait()
}
