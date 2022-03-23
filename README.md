# Go ARRP

Implementation of the [async request-response pattern](https://docs.microsoft.com/en-us/azure/architecture/patterns/async-request-reply) in Go with Gin and Cadence.

## Demo

A short video demoing the implementation of the async request-response pattern:

![Go ARRP Demo GIF](./docs/Go-ARRP.gif)

View the [higher quality video](./docs/Go-ARRP.mp4).

## Architecture

Go ARRP uses Cadence to orchestrate workflows and manages workers to register various workflows and activities.

A Gin HTTP server allows users to use a simple REST API to start jobs, send signals and query job status.

The following routes are available in the REST API:

- `POST /api/v1/jobs/hello-world`: Starts a new job.
- `POST /api/v1/jobs/hello-world/:workflowID/age`: Send the user's age to the workflow.
- `GET /api/v1/jobs/hello-world/:workflowID`: Get the job status.

## Setup

Cadence needs to be started to manage jobs and workers. Start it using Docker Compose:

```sh
$ docker-compose up --build
```

A Cadence domain needs to be created using the Cadence CLI to group workflows:

```sh
$ docker run --network=host --rm ubercadence/cli:master --do go-arrp domain register -rd 1
```

> The `-rd 1` flag is used to set the retention days for the domain. If workflow logs should be kept for longer, increase this value.

Register and start the workers in Go:

```sh
$ go run ./cmd/workers/
```

> Note: Workers should be built and started as binaries in production.

Finally, start the Gin HTTP server in a separate terminal window:

```sh
$ go run ./cmd/server/
```

> Note: Just like workers the server should be built and started as a binary in production.

## Cadence Web UI

The Cadence Web UI is included in the Docker Compose to allow users to view the progress of workflows and activities. To access the it, navigate to http://localhost:8088/domains/go-arrp.

## Tests

In order to run tests Cadence and the workers need to be running. See above on how to start them.

Use `go test` to run all the tests:

```sh
$ go test -v ./...
```

## Configuration

Go ARRP uses Viper to read configuration from the environment variables and the config file in [`configs/application.yml`](./configs/application.yml). Modify this file if you want to change the default values.
