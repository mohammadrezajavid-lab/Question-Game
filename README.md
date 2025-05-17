# Question-Answer

**Question-Answer** is a web-based multiplayer quiz platform built with **Golang** and a RESTful API backend. It allows users to create accounts, choose a topic, and compete live by answering questions against other online players.

This project focuses on providing a smooth and engaging experience for real-time quiz competitions, supported by a scalable and secure backend.

---

## Core Features (Planned)

- User registration and authentication
- Topic-based quiz selection
- Real-time multiplayer question battles
- Scoring and ranking system
- RESTful API built with Go

---

## Getting Started

### Prerequisites

- [Go 1.21+](https://go.dev/dl/)
- [Docker & Docker Compose](https://www.docker.com/)

### Installation

```bash
git clone https://github.com/mohammadrezajavid-lab/Question-Game.git
cd Question-Game

# create mysql container
docker-compose up -d

# build project
go mod tidy
go build -o question-game main.go

# show help project
./question-game -help
```
