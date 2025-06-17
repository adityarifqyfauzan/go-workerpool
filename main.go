package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"playground/processor"
	"playground/usecase"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create two different processor configurations
	processors := make([]*processor.Processor[usecase.Message], 2)

	// processor 1: single worker
	processors[0] = processor.NewProcessor(1, 1, func(msg usecase.Message) error {
		// Simulate processing time
		time.Sleep(5 * time.Second)
		return nil
	})
	processors[0].Start(ctx)

	// processor 2: multiple workers (3 workers)
	processors[1] = processor.NewProcessor(2, 3, func(msg usecase.Message) error {
		// Simulate processing time
		time.Sleep(5 * time.Second)
		return nil
	})
	processors[1].Start(ctx)

	usecase.Initialize(processors[0], processors[1])

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		cancel()

		for _, p := range processors {
			p.Stop()
		}
		os.Exit(0)
	}()

	// hello world endpoint
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello World!"})
	})

	r.GET("/enqueue/:processor", func(c *gin.Context) {
		processorID, err := strconv.Atoi(c.Param("processor"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid processor ID"})
			return
		}

		total, msg, err := usecase.Enqueue(processorID - 1)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":    msg.Content,
			"count":      msg.Count,
			"queue_size": total,
		})
	})

	r.GET("/queue_size/:processor", func(c *gin.Context) {
		processorID, err := strconv.Atoi(c.Param("processor"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid processor ID"})
			return
		}

		size, err := usecase.GetQueueSize(processorID - 1)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"queue_size": size,
		})
	})

	r.Run(":8080")
}
