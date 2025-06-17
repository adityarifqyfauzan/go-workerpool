# Hello Gin

A simple Go web application using the Gin framework.

## Prerequisites

- Go 1.21 or later

## Running the Application

1. Download dependencies:
```bash
go mod download
```

2. Run the application:
```bash
go run main.go
```

The server will start on `http://localhost:8080`. You can test it by visiting the URL in your browser or using curl:

```bash
curl http://localhost:8080
```

You should see a JSON response:
```json
{"message":"Hello World!"}
```

## How to Test

### 1. Testing Single Worker vs Multiple Workers

The application has two processors:
- Processor 1: Single worker
- Processor 2: Three workers

#### Test Single Worker (Processor 1)
```bash
# Add messages to Processor 1's queue
curl http://localhost:8080/enqueue/1
curl http://localhost:8080/enqueue/1
curl http://localhost:8080/enqueue/1

# Check Processor 1's queue size
curl http://localhost:8080/queue_size/1
```

#### Test Multiple Workers (Processor 2)
```bash
# Add messages to Processor 2's queue
curl http://localhost:8080/enqueue/2
curl http://localhost:8080/enqueue/2
curl http://localhost:8080/enqueue/2

# Check Processor 2's queue size
curl http://localhost:8080/queue_size/2
```

### 2. Performance Comparison

To compare the performance between single and multiple workers:

1. Add multiple messages simultaneously:
```bash
# Add 10 messages to Processor 1 (Single worker)
for i in {1..10}; do curl http://localhost:8080/enqueue/1 & done

# Add 10 messages to Processor 2 (Multiple workers)
for i in {1..10}; do curl http://localhost:8080/enqueue/2 & done
```

2. Watch the logs to observe:
   - Processor 1: Messages are processed one at a time
   - Processor 2: Multiple messages are processed simultaneously

### 3. Expected Results

#### Processor 1 (Single Worker)
- Messages are processed sequentially
- Each message takes 5 seconds to process
- Total time for 10 messages: ~50 seconds

#### Processor 2 (Multiple Workers)
- Messages are processed in parallel
- Each message takes 5 seconds to process
- Total time for 10 messages: ~20 seconds (roughly 3x faster)

### 4. Response Format

Each enqueue request returns:
```json
{
    "message": "Hello World! X",
    "count": X,
    "queue_size": Y
}
```

Where:
- `message`: The content of the message
- `count`: The message count
- `queue_size`: Current size of the queue

### 5. Error Handling

The application handles various error cases:
- Invalid processor ID
- Queue processing errors
- System shutdown

Example error response:
```json
{
    "error": "Invalid processor ID"
}
``` 

