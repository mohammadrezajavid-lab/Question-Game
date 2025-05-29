# Question-Answer Game Platform

**Question-Answer** is a web-based multiplayer quiz platform built with **Golang**. It allows users to create accounts, choose a topic, and compete live by answering questions against other online players. This project features a RESTful API backend, real-time matchmaking capabilities, and a robust user authentication system.

---

## 🌟 Core Features

* **User Management**: Secure user registration and JWT-based authentication (access and refresh tokens).
* **Role-Based Access Control (RBAC)**: Differentiates between regular users and administrators with specific permissions.
* **Quiz Categories**: Users can select quiz categories (e.g., "football").
* **Real-time Matchmaking**: Users can join a waiting list for a specific category to be matched with other players. (Core matching logic in scheduler is WIP).
* **Backoffice Operations**: Endpoints for administrative tasks like listing users (requires admin privileges).
* **Database Migrations**: Managed using `sql-migrate`.
* **Configuration Driven**: Application behavior is controlled via a `config.yaml` file and environment variables.
* **Dockerized Environment**: Comes with a `docker-compose.yml` for easy setup of MySQL and Redis services.
* **Graceful Shutdown**: The application handles interrupt signals for a clean shutdown process.
* **Health Check**: Endpoint to verify the status of the server and its dependencies (DB, Redis).
* **Scheduler**: For background job processing, such as the matchmaking logic, using `gocron`.
* **User Presence**: Capability to update and check user's online presence using Redis.

---

## 🛠️ Getting Started

### Prerequisites

* [Go (1.21+ recommended)](https://go.dev/dl/)
* [Docker & Docker Compose](https://www.docker.com/get-started)
* [make (optional, for using Makefile commands)](https://www.gnu.org/software/make/)
* [sql-migrate](https://github.com/rubenv/sql-migrate) (if you wish to run migrations manually outside the application)

### Installation & Setup

1.  **Clone the repository:**
    ```bash
    git clone [YOUR_REPOSITORY_URL] Question-Game
    cd Question-Game
    ```

2.  **Configuration:**
    * The main configuration file is `config.yaml` located in the project root.
    * It includes settings for the HTTP server, database, Redis, JWT authentication, application timeouts, and various service-specific configurations like matching and presence.
    * You can override default configurations using environment variables (e.g., `DATABASE_CFG_DATABASE_HOST=mydbhost`). Variables are prefixed based on their structure in YAML (e.g., `database_cfg.database_host` becomes `DATABASE_CFG_DATABASE_HOST`).
    * Ensure the database and Redis connection details in `config.yaml` match your setup (especially if not using the provided Docker Compose for these services). The default Docker setup exposes MySQL on port `3308` and Redis on `6380` on the host.

3.  **Start External Services (MySQL & Redis):**
    The easiest way to get MySQL and Redis running is using the provided Docker Compose configuration:
    ```bash
    docker-compose up -d
    ```
    This will start:
    * A MySQL container named `gameapp_db` accessible on host port `3308`.
    * A Redis container named `gameapp_redis` accessible on host port `6380`.

4.  **Database Migrations:**
    The application uses `sql-migrate` for database migrations. Migrations are applied automatically when the application starts with the `up` command, or you can manage them manually.
    To apply migrations (create tables and seed data):
    ```bash
    go run main.go -migrate-command=up
    ```
    Or if you have built the binary:
    ```bash
    ./question-game -migrate-command=up
    ```
    Other migration commands:
    * `down`: Rollback the last set of migrations.
    * `status`: Show the status of migrations.
    * `skip` (default): Skip migration step.

5.  **Install Go Dependencies:**
    ```bash
    go mod tidy
    ```

6.  **Build the application:**
    ```bash
    go build -o question-game main.go
    ```

### Running the Application

* **Run the compiled binary:**
    ```bash
    ./question-game
    ```
  You can also specify host and port:
    ```bash
    ./question-game -host=0.0.0.0 -port=8080
    ```

* **Run directly with `go run` (includes migrations):**
    ```bash
    go run main.go -migrate-command=up
    ```
  To run without auto-migrating after the first time:
    ```bash
    go run main.go
    ```

The HTTP server will start (default: `127.0.0.1:8080`) and the matchmaking scheduler will also begin its work.

### Show Help
    ```bash
    ./question-game -help
    ```
    This will show available command-line flags.

---

## 🏗️ Project Structure

The project is organized into the following main directories:

* **`/` (root):** Contains `main.go` (primary application entry point), `config.yaml` (main configuration), `docker-compose.yml`, and this `README.md`.
* **`adapter/`:** Wrappers for external clients (e.g., `/redis` for Redis client).
* **`cmd/`:** Entry points for different binaries.
    * `httpserver/`: For running only the HTTP server component.
    * `scheduler/`: For running the scheduler as a standalone process.
* **`config/`:** Configuration loading (`/httpservercfg`) and service setup (`/setupservices`).
* **`delivery/`:** Presentation layer.
    * `httpserver/`: Echo web server setup, route definitions, middleware, and request/response handlers (`/userhandler`, `/backofficeuserhandler`, `/matchinghandler`).
* **`entity/`:** Core domain models (e.g., `User`, `Game`, `Permission`, `Category`).
* **`pkg/`:** Shared utility packages (e.g., `richerror` for error handling, `hash` for passwords, `normalize` for phone numbers).
* **`repository/`:** Data access layer.
    * `mysql/`: MySQL database interactions, including data access objects (`/usermysql`, `/accesscontrolmysql`) and migrations (`/migrations`).
    * `redis/`: Redis interactions (`/redismatching` for waiting lists, `/redispresence` for user presence).
    * `migrator/`: Database migration utility.
* **`scheduler/`:** Background job processing (e.g., matchmaking logic).
* **`service/`:** Business logic layer (e.g., `userservice`, `authenticationservice`, `matchingservice`).
* **`validator/`:** Input validation logic (`/uservalidator`, `/matchingvalidator`).

---

## 🔌 API Endpoints (Overview)

The following are some of the key API endpoints. (Note: This is based on handler route definitions).

* `GET /health-check`: Checks application health.
* **User Management (`/users/`)**
    * `POST /register`: Register a new user.
    * `POST /login`: Log in an existing user. Updates user presence.
    * `GET /profile`: Get the current user's profile (requires authentication). Updates user presence.
* **Backoffice (`/backoffice/users/`)**
    * `GET /`: List all users (requires admin authentication and `user-list` permission).
* **Matchmaking (`/matching-player/`)**
    * `POST /add-to-waiting-list`: Add the authenticated user to a matchmaking waiting list for a category. Updates user presence.

[//]: # (---)

[//]: # ()
[//]: # (## 🧪 Testing)

[//]: # ()
[//]: # (To run the unit tests included in the project &#40;and any new ones you add&#41;:)

[//]: # ()
[//]: # (```bash)

[//]: # (go test ./... -v)