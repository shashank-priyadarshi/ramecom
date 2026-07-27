# E-Commerce Microservices in Go

A production-inspired E-Commerce backend built using **Golang**, **Microservices**, **gRPC**, **Kafka**, **Docker**, **Prometheus**, and **Grafana**.

This project demonstrates how modern backend systems are designed using service-oriented architecture, asynchronous communication, centralized API routing, and observability.

---

## Architecture

```
                    +------------------+
                    |   API Gateway    |
                    |      :8080       |
                    +--------+---------+
                             |
      -------------------------------------------------
      |                    |                         |
      |                    |                         |
+-------------+     +--------------+        +----------------+
| Auth Service|     | Order Service| -----> | Payment Service|
|   REST API  |     | REST + gRPC  |        | gRPC Server    |
+-------------+     +--------------+        +----------------+
                                                  |
                                                  |
                                              Kafka Event
                                                  |
                                                  ▼
                                         +----------------------+
                                         | Notification Service |
                                         | Kafka Consumer       |
                                         +----------------------+

                     +----------------------+
                     |     Prometheus       |
                     +----------+-----------+
                                |
                                ▼
                        +---------------+
                        |    Grafana    |
                        +---------------+
```

Apps are a **single binary** selected at runtime with `-app` (`gateway`, `auth`, `user`, `order`, `payment`, `notification`).

---

# Tech Stack

| Area          | Choice                                                    |
| ------------- | --------------------------------------------------------- |
| Language      | Go 1.26                                                   |
| HTTP          | Echo (REST)                                               |
| RPC           | gRPC                                                      |
| Messaging     | Apache Kafka (segmentio/kafka-go)                         |
| Database      | MariaDB (auth / user)                                     |
| Config        | `build/.env/.<app>.env` + process env                     |
| Observability | Prometheus, Grafana                                       |
| Infra         | Docker, Docker Compose                                    |
| Dev tools     | Make, [Air](https://github.com/air-verse/air) live reload |

---

# Microservices

| Service              | Flag (`-app`)  | Port  | Description                                   |
| -------------------- | -------------- | ----- | --------------------------------------------- |
| API Gateway          | `gateway`      | 8080  | Single entry point / reverse proxy            |
| Auth Service         | `auth`         | 8081  | JWT authentication APIs, MariaDB              |
| User Service         | `user`         | 8082  | User CRUD APIs, MariaDB                       |
| Order Service        | `order`        | 8083  | Creates orders; calls Payment via gRPC        |
| Payment Service      | `payment`      | 50051 | gRPC server; publishes Kafka payment events   |
| Notification Service | `notification` | —     | Kafka consumer; simulates email notifications |

Shared code lives under `libs/` (config, db, schema, proto, monitoring) and `internal/` (handlers, services, servers).

---

# Communication

## REST

- Client → API Gateway
- API Gateway → Auth / User / Order

## gRPC

```
Order Service
        |
        ▼
Payment Service
```

## Kafka

```
Payment Service
        |
        ▼
Kafka (payment-events)
        |
        ▼
Notification Service
```

---

# Folder Structure

```
ecommerce-microservices/
├── cmd/                      # main: go run ./cmd -app <name>
├── internal/
│   ├── app.go                # app wiring by -app flag
│   ├── handler/
│   ├── middleware/
│   ├── model/
│   ├── repo/
│   ├── server/
│   └── service/
├── libs/
│   ├── config/               # env loading (LoadForApp)
│   ├── db/
│   ├── monitoring/
│   ├── proto/
│   └── schema/
├── build/
│   ├── .air.toml             # Air live-reload config
│   ├── .env/
│   │   ├── .gateway.env
│   │   ├── .auth.env
│   │   ├── .user.env
│   │   ├── .order.env
│   │   ├── .payment.env
│   │   └── .notification.env
│   └── monitoring/           # Prometheus + Grafana configs
├── Dockerfile                # multi-stage single binary (ecom)
├── docker-compose.yml        # infra + all six apps
├── Makefile
├── go.mod
└── README.md
```

---

# Running the Project

## Clone

```bash
git clone https://github.com/rajabhishekmaurya/ecommerce-microservices.git
cd ecommerce-microservices
```

## Quick reference (Make)

```bash
make help                 # all targets
make test                 # go test ./...
make build                # bin/ecom
make infra                # zookeeper, kafka, mariadb, prometheus, grafana
make run APP=gateway      # host process (uses build/.env/.<app>.env)
make watch APP=auth       # live reload via Air
make up                   # full stack: docker compose up --build -d
make down
make run-docker APP=order # one app container
make logs SERVICE=gateway
make tools                # install air
```

`APP` must be one of: `gateway`, `auth`, `user`, `order`, `payment`, `notification`.

---

## Full stack with Docker Compose

One command starts infrastructure **and** all six apps:

```bash
make up
# or: docker compose up --build -d
```

| Service      | Port(s) | Notes                        |
| ------------ | ------- | ---------------------------- |
| gateway      | 8080    | HTTP API entrypoint          |
| auth         | 8081    | HTTP, MariaDB                |
| user         | 8082    | HTTP, MariaDB                |
| order        | 8083    | HTTP → payment gRPC          |
| payment      | 50051   | gRPC, Kafka producer         |
| notification | (none)  | Kafka consumer only          |
| kafka        | 9092    | Host clients use `localhost` |
| mariadb      | 3306    |                              |
| prometheus   | 9090    |                              |
| grafana      | 3000    |                              |

Apps load config from `build/.env/.<app>.env` via Compose `env_file`. Compose also sets Docker DNS overrides (`DB_HOST=mariadb`, `KAFKA_BROKER=kafka:29092`, service URLs, etc.). Host tools can still reach Kafka at `localhost:9092` (dual listeners).

Stop:

```bash
make down
# or: docker compose down
```

---

## Local host development

Start infrastructure only:

```bash
make infra
# or: docker compose up -d zookeeper kafka mariadb prometheus grafana
```

Run apps on the host (localhost-oriented env files under `build/.env/`):

```bash
make run APP=gateway
make run APP=auth
# or: go run ./cmd -app gateway
```

Live reload (Air config: `build/.air.toml`):

```bash
make tools                 # once
make watch APP=gateway
```

Build a host binary:

```bash
make build
./bin/ecom -app gateway
```

---

# Config

- **Host:** `config.LoadForApp` reads `build/.env/.<app>.env` (e.g. `build/.env/.auth.env`).
- **Docker:** Compose `env_file` injects the same files; `environment:` overrides hostnames for the compose network.
- Process env always wins over file values.

---

# Example Flow

1. Client sends request to API Gateway.
2. API Gateway forwards the request to Order Service.
3. Order Service creates an order.
4. Order Service invokes Payment Service via gRPC.
5. Payment Service processes payment.
6. Payment Service publishes a Kafka event.
7. Notification Service consumes the event.
8. Notification Service simulates sending an email notification.

---

# Observability

| UI         | URL                   |
| ---------- | --------------------- |
| Prometheus | http://localhost:9090 |
| Grafana    | http://localhost:3000 |

Scrape targets (in-compose): `gateway:8080`, `auth:8081`, `user:8082`, `order:8083`. Metrics middleware/routes may still be incomplete on some HTTP services.

---

# Features Implemented

- Single-binary multi-app monorepo (`-app` flag)
- API Gateway (reverse proxy)
- JWT Authentication
- User CRUD APIs + MariaDB
- Order Service
- Payment Service (gRPC)
- Kafka producer / consumer
- Docker Compose full stack
- Makefile (test, build, run, docker, air)
- Air live reload
- Prometheus + Grafana config layout

---

# Planned Features

- Product Service
- Inventory Service
- Redis Caching
- Distributed Tracing
- Request Correlation IDs
- Business Metrics
- Broader unit / integration tests
- GitHub Actions CI
- Kubernetes Deployment
- Wire `/metrics` on all HTTP services

---

# Learning Objectives

This project demonstrates:

- Microservice architecture in a Go monorepo
- REST API development
- gRPC communication
- Event-driven architecture with Kafka
- API Gateway pattern
- Docker & Docker Compose
- Config via env files + Compose overrides
- Prometheus / Grafana observability layout
- Productive local DX (Make + Air)

---

# Author

**Raj Abhishek Maurya**

Backend Engineer | Golang | Microservices | Distributed Systems | Kafka
