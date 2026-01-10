# Driver & Driven Architecture in Go

This document explains the **Driver / Driven** concept in Go backend architecture, based on **Clean Architecture** and **Hexagonal (Ports & Adapters)** principles.

---

## 1. Overview

In Go systems, **Driver** and **Driven** describe **how data and control flow move through the application**.

> **Driver → Core (Business Logic) → Driven**

* **Drivers** initiate actions into the system
* **Driven** components are dependencies the system calls out to
* **Core logic must not depend on either**

This separation keeps business logic **testable, maintainable, and framework-agnostic**.

---

## 2. Definitions

### Driver (Inbound)

**Drivers push requests into the system.**

Typical drivers:

* HTTP handlers
* gRPC handlers
* CLI commands
* Cron jobs
* Message queue consumers (Kafka, SQS)

📌 Drivers **depend on the core**, never the other way around.

---

### Driven (Outbound)

**Driven components are external dependencies used by the core.**

Typical driven components:

* Databases (Postgres, MongoDB, TigerBeetle)
* Cache (Redis)
* External APIs
* Message brokers
* File systems

📌 The core depends only on **interfaces**, not concrete implementations.

---

## 3. Dependency Rule (Most Important)

```text
Driver  ───▶  Core  ───▶  Driven
```

Allowed:

* Driver → Core
* Core → Interface (Port)
* Driven → Interface (Implementation)

❌ Not allowed:

* Core → HTTP
* Core → SQL
* Core → Redis
* Core → gRPC

If your business logic imports `net/http`, `database/sql`, or vendor SDKs, the boundary is broken.

---

## 4. Driver Side (Inbound Example)

### HTTP Handler as Driver

```go
type UserHandler struct {
    svc UserService
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    err := h.svc.CreateUser(ctx, CreateUserInput{
        Email: r.FormValue("email"),
    })
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
}
```

**Why this is a Driver**

* Receives external input
* Translates protocol (HTTP → domain)
* Calls core use case

---

## 5. Core (Business Logic)

```go
type UserService interface {
    CreateUser(ctx context.Context, input CreateUserInput) error
}
```

```go
type userService struct {
    repo UserRepository
}

func (s *userService) CreateUser(ctx context.Context, input CreateUserInput) error {
    user := User{Email: input.Email}
    return s.repo.Save(ctx, user)
}
```

**Core rules**

* No HTTP / gRPC / SQL imports
* Pure Go logic
* Depends only on interfaces

---

## 6. Driven Side (Outbound Example)

### Repository Interface (Port)

```go
type UserRepository interface {
    Save(ctx context.Context, user User) error
}
```

### Concrete Implementation (Adapter)

```go
type PostgresUserRepository struct {
    db *sql.DB
}

func (r *PostgresUserRepository) Save(ctx context.Context, user User) error {
    _, err := r.db.ExecContext(
        ctx,
        "INSERT INTO users(email) VALUES ($1)",
        user.Email,
    )
    return err
}
```

**Why this is Driven**

* Called by the core
* Talks to external systems
* Implements a core-defined interface

---

## 7. Ports & Adapters Mapping

| Concept     | Role           | Go Representation             |
| ----------- | -------------- | ----------------------------- |
| Driver Port | Inbound        | Use case interface            |
| Driven Port | Outbound       | Repository / client interface |
| Adapter     | Implementation | HTTP handler, DB repo         |

---

## 8. Project Structure Example

### Clean Architecture Style

```text
/internal
  /domain
  /usecase              ← CORE
  /ports
    inbound.go           ← Driver ports
    outbound.go          ← Driven ports
  /adapters
    /http                ← Driver adapters
    /postgres            ← Driven adapters
```

### Simpler Go Layout

```text
/internal
  /app
    /service              ← Core
    /repository           ← Interfaces
  /transport
    /http                 ← Drivers
  /infra
    /postgres             ← Driven
```

---

## 9. Anti-patterns to Avoid

❌ Business logic importing infrastructure:

```go
import "database/sql"
```

❌ HTTP logic inside use cases
❌ Repositories returning HTTP errors
❌ Passing framework types (`*http.Request`) into core

---

## 10. Benefits

* ✅ Easy unit testing
* ✅ Infrastructure can change without rewriting logic
* ✅ Clear ownership of responsibilities
* ✅ Better long-term maintainability

---

## 11. Real-world Mapping (Financial Systems)

| Component           | Role   |
| ------------------- | ------ |
| gRPC / REST API     | Driver |
| Balance calculation | Core   |
| TigerBeetle client  | Driven |
| Kafka producer      | Driven |

---

## 12. Rule of Thumb

> **Drivers depend on Core**
> **Core depends on interfaces**
> **Driven implements interfaces**

If this rule is followed, the architecture is correct.

---
