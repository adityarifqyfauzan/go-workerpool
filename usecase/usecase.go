package usecase

import (
	"fmt"
	"playground/processor"
)

type Message struct {
	Content string
	Count   int
}

var (
	processors []*processor.Processor[Message]
)

func Initialize(singleWorker, multiWorker *processor.Processor[Message]) {
	processors = []*processor.Processor[Message]{singleWorker, multiWorker}
}

func Enqueue(processorID int) (int, Message, error) {
	if processorID < 0 || processorID >= len(processors) {
		return 0, Message{}, fmt.Errorf("invalid processor ID: %d", processorID)
	}

	p := processors[processorID]
	msg := Message{
		Content: fmt.Sprintf("Hello World! %d", p.GetQueueSize()),
		Count:   p.GetQueueSize(),
	}
	p.Enqueue(msg)
	return p.GetQueueSize(), msg, nil
}

func GetQueueSize(processorID int) (int, error) {
	if processorID < 0 || processorID >= len(processors) {
		return 0, fmt.Errorf("invalid processor ID: %d", processorID)
	}

	return processors[processorID].GetQueueSize(), nil
}
