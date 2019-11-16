# Screenshot as a Service

This is an API service which provides web screenshots of any website.

## Usage

### Running the service

To setup up the whole stack using `docker-compose` just issue the following command in the root of the project.

```shell
docker-compose up
```

The server will start listening on port 80 of the host machine.

### Taking screenshots

To create screenshots of the url you just need to issue a POST request the `api/request` endpoint with your urls. Just like this:

```shell
curl -X POST -d 'urls=https://google.com%3Bhttps://github.com%3Bhttps://stackoverflow.com' http://localhost/api/new
```

The response would be something like this:

```json
{
  "results": [
    "http://localhost/api/result/aHR0cHM6Ly9nb29nbGUuY29t",
    "http://localhost/api/result/aHR0cHM6Ly9naXRodWIuY29t",
    "http://localhost/api/result/aHR0cHM6Ly9zdGFja292ZXJmbG93LmNvbQ=="
  ]
}
```

Each url is for retreving the status of each screenshot. The API server may return one of the following HTTP status codes, depending on the status of the request.

* `200 OK` - The request is successfully finished and the resulting image url is returned
* `204 No Content` - The request has not yet been processed
* `404 Not Found` - The requested url is not found in our database
* `500 Internal Server Error` - An unexpected error occurred, you may try again
* `503 Service Unavailable` - The request has failed permanently for some reason, you may check the logs
* `422 Unprocessable Entity` - The supplied request is invalid

## Tests

There are two type of tests in this project.

### Unit Tests

To run unit tests, no setup is required and you just need the helper make command:

```shell
make test
```

### Integration Tests

To run integration tests, there needs to be a running instance of RabbitMQ server and a MongoDB server. To run them without hassle just issue the helper command:

```shell
make test-integration
```

## Architecture

This project is made of 5 main components and since none of them has any dependency to each other and they solely depend on interfaces, they can be swapped out with different implementations easily without affecting other components.

A brief explanation of each component is provided below:

### API

This component is the outermost layer which is responsible for receiving screenshot requests from user and then displaying the response. It receives requests and after validation passes it to the next layer.

### Dispatcher

Dispatcher acts as kind of coordinator. It receives requests from upper layer (api) and dispatches jobs to the job queue and then receives the results back from it. It also stores requests using the Datastore interface.

### Storage

This package provides persistence capabilities for the service. Dispatcher uses this to store requests and their statuses.  
I've used MongoDB here to avoid the hassle of SQL queries. But it is easy enough to replace it with any other database implementations.

### JobQ

JobQ does what it's name implies. It queues screenshot jobs and processes them sequentially to avoid overloading the system. It uses the `webCapture` interface to proccess jobs.

### Screenshot

Screenshot is where everything happens. This is the innermost layer which is unaware of any other layers.  I
I have used `PhantomJS` as an implementation of the `Capture` interface, but it can also use any other tool such as `Chromium Headless` or even a third-party provider.

## Scaling

Since this service uses a message queue, it is inherently scalable and since it's already containerized, it can be deployed to the cloud without much work. To scale it out, just increase the number of pods (containers) and you're good to go.

## Minio

At the last minute I thought it'd be a good idea to have a cloud native solution for file storage as well (the files were previously stored on a mounted volume) so I added Minio as the file storage service. Minio is compatible with Amazon's S3 protocol so it's a cinch to replace it with AWS Cloud storage.

Note: Having a filestore is not mandatory for the system to work and it works perfectly without it. That's why I did not bring it up in Architecture section above.

## Multi-Mode

There is a config with which you can specify how many workers should be created on each instance. By setting `WORKERS_PER_INSTANCE` env variable to a number between 1 and 10 you can utilize more system resources therefor increase the performance.
