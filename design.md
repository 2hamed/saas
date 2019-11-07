# SaaS
Screenshot as a service

## Components

1. API
2. Screenshot
3. Storage
4. Job queue
5. UI

### API
The API component is responsible for serving requests to generate screenshots from user provided urls.

### Screenshot
This is the main service which receives a url and generates a screenshot of the page.

### Storage
This component is responsible for storing the processed requests (screenshots) inside a database.

### Job queue
This component receives user requests and queues them to be processed by the Screenshot service.

### UI
Maybe a simple webpage with a form inside to receive user requests.