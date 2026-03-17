# Shipment Management Service
## 🚀 Quick Start

### Prerequisites

* Docker and Docker Compose
* Go (1.25+)
* Protocol Buffers Compiler (protoc) with Go plugins
* Make (optional, but recommended for automation)

### Code Generation
Before running the service or tests, you must generate the gRPC code from the Protocol Buffer definitions:

```
make proto
```

### Running the Service

The entire environment is containerized. To build and start the service (including the PostgreSQL database and automatic schema migrations):
```
make up
```

If you don't have make installed:

```
docker-compose up -d --build
```

The service will be available at localhost:50051. To view application logs:
```
make logs
```

Stopping the Service
```
make down
```

## 🧪 Testing

**Running Unit Tests**

To run the project tests with verbose output:
```
go test -v ./...
```

**Code Coverage**

To check the test coverage:
```
go test -cover ./...
```

## 🏗 Architecture Overview

The project is structured according to **Clean Architecture (Domain-Driven Design)** to ensure separation of concerns, high testability, and independence from external libraries:

* Domain (internal/domain): Contains core business entities (Shipment, ShipmentEvent) and custom errors. This layer has zero dependencies on other layers.

* UseCase (internal/usecase): Contains business logic and the shipment state machine. It defines repository interfaces to remain storage-agnostic.

* Delivery (internal/delivery/grpc): Handles gRPC transport. It maps incoming Protobuf requests to domain models and translates internal errors into gRPC Status Codes.

* Infrastructure (internal/infrastructure): Implements external interfaces.

    * Postgres: Persistent storage implementation.

    * Memory: In-memory implementation for testing.

    * Logger: Structured logging with log/slog.

## 💡 Design Decisions

* Strict State Machine: Shipment status transitions (e.g., PENDING -> PICKED_UP) are strictly validated in the UseCase layer. This prevents inconsistent data states and ensures business rule compliance.

* PostgreSQL UPSERT Strategy: The repository uses an UPSERT (ON CONFLICT) strategy for saving shipments. This allows the system to handle both creation and updates through a single method efficiently.

* Structured Logging: Using log/slog for JSON-formatted logs. This ensures the service is ready for modern log management systems like ELK or Datadog.

* Server Reflection: gRPC Reflection is enabled in main.go to facilitate easy testing with tools like Postman or Evans without requiring manual .proto file imports.

* Dependency Injection (DI): All layers are connected in main.go. This decouples the business logic from the infrastructure, allowing easy swapping of components (e.g., switching from Postgres to Memory DB).

* Graceful Shutdown: The service handles SIGINT and SIGTERM signals to safely stop the gRPC server and close database connections.

## 📋 Assumptions

* Linear Workflow: It is assumed that shipments follow a linear lifecycle: PENDING -> PICKED_UP -> IN_TRANSIT -> DELIVERED. Skipping stages or moving backwards is restricted.

* Identification: UUIDs (v4) are used for all record IDs to ensure uniqueness and prevent ID enumeration.

* Authentication: For the scope of this test task, authentication and authorization layers are omitted.

* Database Initial State: The database is initialized via a Docker volume and an init.sql script. If the database schema changes, a manual volume reset (docker-compose down -v) is required for local development.

## 🛠 Project Structure
```
.
├── api/proto/           # Protobuf definitions and generated Go code
├── cmd/          # Application entrypoint (main.go)
├── internal/
│   ├── delivery/        # gRPC Handlers & Protocol Mappers
│   ├── domain/          # Core Business Entities & Domain Errors
│   ├── infrastructure/  # DB Repositories (Postgres, Memory) & Logger
│   └── usecase/         # Business Logic & Service Interfaces
├── migrations/          # SQL initialization scripts (init.sql)
├── Dockerfile           # Multi-stage Docker build
├── docker-compose.yml   # Infrastructure orchestration
└── Makefile             # Utility commands for development
```

## 📡 API Testing (Postman)

    1. Create a new gRPC Request in Postman.

    2. Set URL to localhost:50051.

    3. Disable TLS (click the lock icon to unlock).

    4. Use Server Reflection to load methods.

Create Shipment (CreateShipment) 
```
{
  "reference_number": "REF-100",
  "origin": "Astana",
  "destination": "Almaty",
  "driver_details": "John Doe",
  "unit_details": "Truck A1",
  "shipment_amount": 1500.50,
  "driver_revenue": 1000.00
}
```

Get Shipment Details (GetShipment)

```
{
  "id": "PASTE_YOUR_UUID_HERE"
}
```

Update Status (UpdateShipmentStatus)

**Valid transitions: PENDING ➔ PICKED_UP ➔ IN_TRANSIT ➔ DELIVERED.**

```
{
  "id": "PASTE_YOUR_UUID_HERE",
  "new_status": "STATUS_PICKED_UP",
  "note": "Cargo loaded at warehouse"
}
```

Get History (GetShipmentHistory)
```
{
  "id": "PASTE_YOUR_UUID_HERE"
}
```