# Question-Answer Game Platform

**Question-Answer** is a web-based multiplayer quiz platform built with **Golang**. It allows users to create accounts,
choose a topic, and compete live by answering questions against other online players. This project features a RESTful
API backend, real-time matchmaking capabilities, a robust user authentication system, and can be run as a monolith or as
separate microservices.

## 🛠️ Getting Started

### Prerequisites

* [Go (1.21+ recommended)](https://go.dev/dl/)
* [Docker & Docker Compose](https://www.docker.com/get-started)
* [make (optional, for using Makefile commands)](https://www.gnu.org/software/make/)
* [sql-migrate](https://github.com/rubenv/sql-migrate) (if you wish to run migrations manually)
* [protoc](https://grpc.io/docs/protoc-installation/)

### Installation & Setup

1. **Clone the repository:**
   ```bash
   git clone https://github.com/mohammadrezajavid-lab/Question-Game.git Question-Game
   cd Question-Game
   ```

2. **Configuration:**
    * The main configuration file is `config.yaml` located in the project root.
    * It includes settings for the HTTP server, gRPC server, database, Redis, JWT authentication, and more.
    * You can override configurations using environment variables. For example, `database_cfg.database_host` in YAML
      becomes `DATABASE_CFG_DATABASE_HOST` as an environment variable.
    * The default Docker setup exposes MySQL on port `3308` and Redis on `6380` on the host. Ensure these match your
      `config.yaml`.

3. **Start External Services (MySQL & Redis):**
   Use the provided Docker Compose configuration to start the necessary services:
   ```bash
   docker-compose up -d
   ```
   This will start:
    * A MySQL container named `gameapp_db` on host port `3308`.
    * A Redis container named `gameapp_redis` on host port `6380`.

4. **Install Go Dependencies:**
   ```bash
   go mod tidy
   ```

5. **Database Migrations:**
   The application can automatically run migrations on startup.
   To apply migrations (create tables and seed data):
   ```bash
   go run main.go -migrate-command=up
   ```
   Other migration commands:
    * `down`: Rollback the last set of migrations.
    * `status`: Show the status of migrations.
    * `skip` (default): Skip the migration step.

6. **Build Protocol Buffers:**
   If you modify any `.proto` files in the `contract/protobuf` directory, you'll need to regenerate the Go code.
   ```bash
   # For presence service
   protoc --proto_path=contract/protobuf/presence --go_out=contract/goprotobuf/presence --go_opt=paths=source_relative --go-grpc_out=contract/goprotobuf/presence --go-grpc_opt=paths=source_relative ./contract/protobuf/presence/presence.proto

   # For matching service
   protoc --proto_path=contract/protobuf/matching --go_out=contract/goprotobuf/matching --go_opt=paths=source_relative ./contract/protobuf/matching/matching.proto

   # For notification service
   protoc --proto_path=contract/protobuf/notification --go_out=contract/goprotobuf/notification --go_opt=paths=source_relative ./contract/protobuf/notification/notification.proto
  
   # For game service   
   protoc --proto_path=contract/protobuf/game --go_out=contract/goprotobuf/game --go_opt=paths=source_relative ./contract/protobuf/game/created_game.proto
   ```

### Testing the Application

* **Run Unit Tests**
    ```bash
    go test -v ./...
    ```

### Running the Application

You can run the application as a monolith or as individual microservices.

* **Run as a Monolith (includes HTTP server, scheduler, gameService, metricServer, profilingServer):**
    ```bash
    # This will also run migrations
    go run main.go -migrate-command=up
    ```
    * **After migrating the project tables, run or import the following SQL script to insert questions into the question table.**
      ```
      file: ./repository/mysql/migrations/seed_question.sql
      ```
    * **To run without auto-migrating after the first time:**
      ```bash
      go run main.go
      ```
    * **To run Presence GRPC server:**
      ```bash
        go run ./cmd/presenceserver/main.go
      ```
    * **WebSocket GateWay:**
      ```bash
      go run cmd/websocketserver/main.go
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
    * **Game Service:**
        ```bash
        go run cmd/gameservice/main.go
        ```
    * **Quiz Service:**
        ```bash
        go run cmd/gameservice/main.go
        ```
    * **WebSocket GateWay:**
        ```bash
        go run cmd/websocketserver/main.go
        ```
    * **Scheduler:**
        ```bash
        go run cmd/scheduler/main.go
        ```

* **Test WebSocket Gateway:**
    * **Send heartbeat request for Upsert Presence User**
        ```bash
        websocat --header="Authorization: Bearer <jwt-token> --header="Origin: http://127.0.0.1:3000" ws://<websocat-host>:8090/ws
        {"event":"heartbeat "}
        ```

The HTTP server will start on `127.0.0.1:8080` and the metrics server on `127.0.0.1:2112` by default.

## 🔌 API Endpoints

* `GET /health-check`: Checks the health of the application and its dependencies (MySQL, Redis).
* **User Management (`/users/`)**
    * `POST /register`: Register a new user.
    * `POST /login`: Log in an existing user and receive JWT tokens. Updates user presence.
    * `GET /profile`: Get the current user's profile. Requires authentication and updates user presence.
* **Auth Management (`/auth/`)**
    * `POST /refresh`: Send Refresh Token and Response JWT tokens.
* **Backoffice (`/backoffice/users/`)**
    * `GET /`: List all users. Requires admin authentication and the `user-list` permission.
* **Matchmaking (`/matching-player/`)**
    * `POST /add-to-waiting-list`: Add the authenticated user to a matchmaking waiting list for a specific category.
      Updates user presence.
