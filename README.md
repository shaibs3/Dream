# Dream Project

Dream is a Go microservice for ingesting, parsing, and storing process data from different operating systems (Linux, Windows, macOS) for research analytics. It supports analytics queries by user, faculty, OS, and time, and is built with Kafka, PostgreSQL, and Docker Compose.


## High level architecture

## Technology Choices

### Why Go?
Go was chosen for its simplicity, performance, and strong support for concurrency. Its lightweight goroutines make it ideal for building microservices that handle high-throughput data processing efficiently.

### Why Kafka?
Kafka was selected as the messaging system due to its ability to handle large volumes of data with high reliability and scalability. It provides robust features for distributed systems, making it suitable for real-time data ingestion and processing.

### Why PostgreSQL?
PostgreSQL was chosen as the database for its advanced features, reliability, and strong support for complex queries. It is well-suited for handling structured data and supports analytics queries efficiently.

### Why Docker Compose Locally?
Docker Compose was chosen for local development to simplify the setup and management of dependencies like Kafka and PostgreSQL. It allows developers to run the entire stack with minimal configuration.

### Kubernetes in Production
While Docker Compose is used locally, Kubernetes (K8s) is planned for production deployment to ensure scalability, fault tolerance, and efficient resource management in a distributed environment.

> **Note:** Due to the scope and requirements of this project, cloud storage was not used. Instead, the database schema uses a string field to hold the full command output for each process record. This simplifies local development and testing, but may be revisited for scalability in a production environment.

![Image description](img/Arch.png)


## Prerequisites

- [Go 1.21 or later](https://go.dev/doc/install)
- [Make](./Makefile) (for using Makefile commands)
- [Docker](https://www.docker.com/get-started) (required for running the app and dependencies locally)
- [Docker Compose](https://docs.docker.com/compose/) (for orchestrating multi-container setup)

## Running the Application

### Using Make

The project includes a Makefile with common development commands:

```bash
# Build the application
make build

# Run the application
make run

# Run tests
make test

# Clean build files
make clean

# Install dependencies
make deps

# Run linter
make lint

# Show all available commands
make help
```

### Manual Run

Alternatively, you can run the application directly:

1. Clone this repository
2. Run the application:
   ```bash
   go run main.go
   ```
3. Open your browser and visit `http://localhost:8080`

## Features

- Simple HTTP server
- Root endpoint that returns "Hello, World!"
- Proper error handling and logging
- Makefile for common development tasks

## Running the App with Docker Compose

To build and run the application and all dependencies (Postgres, Kafka, etc.) using Docker Compose:

```sh
docker-compose build
docker-compose up -d
```

- The app will be available at http://localhost:8080
- Logs can be viewed with:
  ```sh
  docker-compose logs -f app
  ```
- To stop all services:
  ```sh
  docker-compose down
  ```

## System Tests

System tests are located in `system_test.go` and are designed to verify the end-to-end functionality of the application.

### How to Add a System Test
1. Open or create `system_test.go` in the project root.
2. Write your test using Go's `testing` package and `require`/`assert` from `testify` for assertions.
3. Use the provided helpers and data in the `testdata/` directory as needed.

### How to Run System Tests

1. Ensure the app and dependencies are running (see above).
2. In a new terminal, run:
   ```sh
   go test -v system_test.go
   ```
   or to run all tests:
   ```sh
   go test -v ./...
   ```

- The test will automatically set up the required environment variables for the test database.
- Make sure your test database is accessible at `localhost:5432` (or adjust the connection string in the test).

---

For more details on scripts and utilities, see the `scripts/README.md` file. 