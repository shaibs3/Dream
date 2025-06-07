# Go Hello World Application

A simple HTTP server that responds with "Hello, World!" on the root endpoint.

## Prerequisites

- Go 1.21 or later
- Make (for using Makefile commands)

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