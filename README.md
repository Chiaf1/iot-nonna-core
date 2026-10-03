# Service **`iot-nonna-core`**

> Part of the **iot-nonna** project. For the whole system and the Docker Compose deployment, see [iot-nonna-containers](https://github.com/Chiaf1/iot-nonna-containers).

## Purpose

**`iot-nonna-core`** is the central domain service of the "nonna" IoT system.  
It is responsible for:

- **managing the IoT business model** (rooms, devices, sensors, types)
- **exposing the collected data** in a clear, consumable form
- **hiding the complexity of the database** from the rest of the system
- **acting as the single source of truth** for configuration and data reads

It is not a data ingestion service and not an active control service. It is **the information and configuration core** that all the other services rely on.

---

## Responsibilities

### ✅ What it **MUST** do

1. **Provisioning and configuration**
   - Creating and managing:
     - rooms
     - device_type
     - sensor_type
     - devices
     - device ↔ sensor associations
   - Applying:
     - validations
     - logical constraints
     - consistency rules (e.g. a device must have a valid type)
2. **Simplified data access**
   - Expose the contents of the database in a form that is:
     - readable
     - structured
     - stable over time
   - Remove the need for pgAdmin or manual queries
3. **Reading historical data**
   - Provide endpoints designed for:
     - dashboards
     - charts
     - time-based analysis
   - Wrap:
     - time filters (`from`, `to`, `limit`)
   - Return results ready for the frontend (not "raw SQL")
4. **Domain stability**
   - Be the **stable contract** between the DB and its consumers
   - Allow the DB schema to evolve without breaking the frontend

---

## Technology choices

- **Language**: Go
- **API style**: REST
- **HTTP router**: `chi`
- **Database**: PostgreSQL
- **DB access**: native driver (`pgx`)
- **DB role separation**:
  - `iot-nonna-ingest`: reads metadata + writes readings
  - `iot-nonna-core`: manages business tables + reads

### Why Go + chi

- Go gives:
  - predictable performance
  - structural simplicity
  - excellent support for domain backend services
- `chi` is:
  - minimal
  - built on `net/http`
  - free of magic and artificial patterns
  - good for clean, maintainable architectures over time

---

## Logical structure of the service

The service is **layered**, not built from "fat handlers".

Key concepts (independent of the code):

1. **HTTP layer**
   - request parsing
   - status codes
   - JSON serialization

2. **Domain / Service layer**
   - business rules
   - logical validations
   - orchestration of operations

3. **Persistence layer**
   - SQL queries
   - DB access
   - no domain logic

Each layer has **clear, non-overlapping** responsibilities.

---

## Quick start

The service needs a reachable PostgreSQL database.

```bash
cp config.example.yaml config.yaml
# edit DbURL in config.yaml with the parameters of your database
RUN_MIGRATIONS=true RUN_SEEDING=true go run ./cmd/iot-nonna-core
```

- The configuration file is `./config.yaml`; a different path can be set with the `CONFIG_PATH` variable.
- `RUN_MIGRATIONS=true` applies the migrations at startup; `RUN_SEEDING=true` runs the seeding. Without these variables neither is run.
- The API listens on port `3030`.

---

## Migrations and seeding

### Migrations (required)

Migrations are used to:

- version the database schema
- make the evolution of the DB reproducible and controlled

Every structural change (tables, columns, indexes) is:

- described in a numbered SQL file
- applied in order

This removes:

- manual changes
- differences between environments
- "I don't know how this field came to be"

Migrations are **part of the service**, not a separate service.

---

### Seeding (initial data)

Seeding is used to:

- populate the DB with initial data
- make startup and testing immediate
- avoid manual configuration

Examples:

- known sensor types
- known device types
- test configuration

Seeding can be:

- run manually
- or run at startup in non-production environments

## Planned

- Aggregation and downsampling of readings: not implemented yet. The `/readings` endpoints return the raw readings, filtered by time range and maximum number of rows.

---

## API documentation

Interactive documentation is available through Swagger UI at the endpoint:

```
GET /swagger/index.html
```

It lets you explore all the endpoints, see the request/response models
and make calls directly from the browser.

To regenerate the documentation after changing the handlers:

```bash
swag init -g cmd/iot-nonna-core/main.go
```
