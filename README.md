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