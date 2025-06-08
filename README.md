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

The project includes a [Make](./Makefile) with common development commands:


## Running the App with Docker Compose

To build and run the application and all dependencies (Postgres, Kafka, etc.) using Docker Compose:

```sh
make docker-up-build
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

## Unit Tests
  ```sh
  make test
  ```
## System Tests

System tests are located in `system_test.go` and are designed to verify the end-to-end functionality of the application
  ```sh
  make test-system
  ```

---

For more details on scripts and utilities, see the [scripts/README.md](scripts/README.md) file. 