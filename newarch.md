# New Microservice Architecture

## Services

### 1. API Gateway (Heimdall)
This service receives capture requests from endusers and queues using the Queue service

### 2. Queue Dispatch (Odin)
This services receives capture requests from the API Gateway and pushes it inside a message broker.

### 3. Queue Consumer (Odin)
This services listens on the Queue for new capture jobs and makes an internal request to the Capture service.

### 4. Capture (Thor)
This service performs the capture job and pushes the result back to the queue dispatcher.

### 5. Persistence (Loki)
This service persistes capture reuqests and their statuses to its backing database.


## Notes

### Request IDs
Requests IDs are generated using UUID.

### Service to Service communication
For internal communication between services, gRPC is used.

### Data Storage
For storage Google Spanner will be used.

### Queue (Message Broker)
Google Pub/Sub will be used as the queue.