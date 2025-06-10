# Question-Answer Game Platform

**Question-Answer** is a web-based multiplayer quiz platform built with **Golang**. It allows users to create accounts, choose a topic, and compete live by answering questions against other online players. This project features a RESTful API backend, real-time matchmaking capabilities, a robust user authentication system, and can be run as a monolith or as separate microservices.

---

## 🌟 Core Features

* **User Management**: Secure user registration and JWT-based authentication (access and refresh tokens).
* **Role-Based Access Control (RBAC)**: Differentiates between regular users and administrators with specific permissions for actions like listing users.
* **Quiz Categories**: Users can select from various quiz categories (e.g., "football", "history", "art") to be matched with others.
* **Real-time Matchmaking**: Users can join a waiting list for a specific category. A scheduler periodically runs to match online players from the waiting list.
* **User Presence**: The system tracks user online status using a gRPC-based presence service backed by Redis. Presence is updated on key actions like logging in, viewing a profile, or joining a waiting list.
* **Backoffice Operations**: Endpoints for administrative tasks like listing all users, accessible only to users with admin privileges.
* **Database Migrations**: Managed using `sql-migrate`. Migrations can be run automatically when the application starts.
* **Configuration Driven**: Application behavior is controlled via a `config.yaml` file and can be overridden by environment variables.
* **Dockerized Environment**: Comes with a `docker-compose.yml` for easy setup of MySQL and Redis services.
* **Graceful Shutdown**: The application handles interrupt signals for a clean shutdown of all its components, including the HTTP server, metrics server, and scheduler.
* **Health Check**: An endpoint to verify the status of the server and its dependencies (Database, Redis).
* **Metrics**: Exposes Prometheus metrics for monitoring application health and performance.
* **Microservice Architecture**: The project is structured to be run as a single monolithic service or deployed as individual microservices for the HTTP server, presence server, and scheduler.

---

## 🛠️ Getting Started

### Prerequisites

* [Go (1.21+ recommended)](https://go.dev/dl/)
* [Docker & Docker Compose](https://www.docker.com/get-started)
* [make (optional, for using Makefile commands)](https://www.gnu.org/software/make/)
* [sql-migrate](https://github.com/rubenv/sql-migrate) (if you wish to run migrations manually)
* [protoc](https://grpc.io/docs/protoc-installation/)

### Installation & Setup

1.  **Clone the repository:**
    ```bash
    git clone [YOUR_REPOSITORY_URL] Question-Game
    cd Question-Game
    ```

2.  **Configuration:**
    * The main configuration file is `config.yaml` located in the project root.
    * It includes settings for the HTTP server, gRPC server, database, Redis, JWT authentication, and more.
    * You can override configurations using environment variables. For example, `database_cfg.database_host` in YAML becomes `DATABASE_CFG_DATABASE_HOST` as an environment variable.
    * The default Docker setup exposes MySQL on port `3308` and Redis on `6380` on the host. Ensure these match your `config.yaml`.

3.  **Start External Services (MySQL & Redis):**
    Use the provided Docker Compose configuration to start the necessary services:
    ```bash
    docker-compose up -d
    ```
    This will start:
    * A MySQL container named `gameapp_db` on host port `3308`.
    * A Redis container named `gameapp_redis` on host port `6380`.

4.  **Database Migrations:**
    The application can automatically run migrations on startup.
    To apply migrations (create tables and seed data):
    ```bash
    go run main.go -migrate-command=up
    ```
    Other migration commands:
    * `down`: Rollback the last set of migrations.
    * `status`: Show the status of migrations.
    * `skip` (default): Skip the migration step.

5.  **Install Go Dependencies:**
    ```bash
    go mod tidy
    ```

6.  **Build Protocol Buffers:**
    If you modify any `.proto` files in the `contract/protobuf` directory, you'll need to regenerate the Go code.
    ```bash
    # For presence service
    protoc --proto_path=contract/protobuf/presence --go_out=contract/goprotobuf/presence --go_opt=paths=source_relative --go-grpc_out=contract/goprotobuf/presence --go-grpc_opt=paths=source_relative ./contract/protobuf/presence/presence.proto

    # For matching service
    protoc --proto_path=contract/protobuf/matching --go_out=contract/goprotobuf/matching --go_opt=paths=source_relative ./contract/protobuf/matching/matching.proto

    # For notification service
    protoc --proto_path=contract/protobuf/notification --go_out=contract/goprotobuf/notification --go_opt=paths=source_relative ./contract/protobuf/notification/notification.proto
    ```

### Running the Application

You can run the application as a monolith or as individual microservices.

* **Run as a Monolith (includes HTTP server and scheduler):**
    ```bash
    # This will also run migrations
    go run main.go -migrate-command=up
    ```
  To run without auto-migrating after the first time:
    ```bash
    go run main.go
    ```
  To run Presence GRPC server:
    ```bash
      go run ./cmd/presenceserver/main.go
    ```

* **Run as Individual Microservices:**
    * **HTTP Server:**
        ```bash
        go run cmd/httpserver/main.go -migrate-command=up
        ```
    * **Presence gRPC Server:**
        ```bash
        go run cmd/presenceserver/main.go
        ```
    * **Scheduler:**
        ```bash
        go run cmd/scheduler/main.go
        ```

The HTTP server will start on `127.0.0.1:8080` and the metrics server on `127.0.0.1:2112` by default.

---

## 🏗️ Project Structure

The project follows a modular structure to separate concerns and support microservices deployment.

* **`/` (root):** Contains the main entry point for the monolithic application (`main.go`), configuration files (`config.yaml`, `docker-compose.yml`), and this `README.md`.
* **`adapter/`:** Wrappers for external clients, such as Redis, gRPC clients, and message publishers.
* **`cmd/`:** Entry points for different binaries, allowing the application to be run as separate microservices.
* **`config/`:** Configuration loading and service setup logic.
* **`contract/`:** Protocol definitions, including protobuf files and generated Go packages for gRPC.
* **`delivery/`:** Presentation layer, handling communication with the outside world.
    * `httpserver/`: Echo web server setup, route definitions, middleware, and request handlers.
    * `grpcserver/`: gRPC server implementation for services like presence.
    * `metricsserver/`: Prometheus metrics server.
* **`entity/`:** Core domain models of the application (e.g., `User`, `Game`, `Permission`).
* **`pkg/`:** Shared utility packages for common functionalities like error handling, hashing, and normalization.
* **`repository/`:** Data access layer responsible for database and Redis interactions.
    * `mysql/`: MySQL database logic, including migrations.
    * `redis/`: Redis repository for caching, presence, and matchmaking waiting lists.
* **`scheduler/`:** Background job processing for tasks like matchmaking.
* **`service/`:** Business logic layer where the core application logic resides.
* **`validator/`:** Input validation logic for API requests.

---

## 🔌 API Endpoints

* `GET /health-check`: Checks the health of the application and its dependencies (MySQL, Redis).
* **User Management (`/users/`)**
    * `POST /register`: Register a new user.
    * `POST /login`: Log in an existing user and receive JWT tokens. Updates user presence.
    * `GET /profile`: Get the current user's profile. Requires authentication and updates user presence.
* **Backoffice (`/backoffice/users/`)**
    * `GET /`: List all users. Requires admin authentication and the `user-list` permission.
* **Matchmaking (`/matching-player/`)**
    * `POST /add-to-waiting-list`: Add the authenticated user to a matchmaking waiting list for a specific category. Updates user presence.

---

## ⚙️ Configuration Details

The application's behavior is managed by the `config.yaml` file.

| Section | Key | Description |
| :--- | :--- | :--- |
| **`grpc_server_cfg`** | `host`, `port`, `network` | Configuration for the gRPC server (used by the presence service). |
| **`httpserver_cfg`** | `host`, `port` | Configuration for the main HTTP server. |
| **`database_cfg`** | Various | MySQL connection details, including user, password, host, and connection pool settings. |
| **`redis_cfg`** | Various | Redis connection details. |
| **`auth_cfg`** | `sign_key`, expirations | JWT token settings, including the secret key and token expiration times. |
| **`app_cfg`** | `gracefully_shutdown_timeout` | Timeout for a graceful shutdown. |
| **`matching_cfg`** | timeouts, duration | Settings for the matchmaking service, like how long a user should wait before being considered offline. |
| **`presence_cfg`** | `expiration_time`, `prefix` | Configuration for the user presence service, such as Redis key prefix and expiration time. |
| **`scheduler_cfg`** | `crontab` | Crontab string for scheduling the matchmaking job. |
| **`logger_cfg`**| Various | Configuration for the logger, including file name, size, and rotation policy. |
| **`metrics_cfg`**| `host`, `port` | The address for the Prometheus metrics server to listen on. |