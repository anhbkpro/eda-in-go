# Event-Driven Architecture in Go - Documentation

This documentation provides comprehensive guides and references for understanding and working with this event-driven e-commerce system built in Go.

## 📚 Documentation Overview

| Document | Description |
|----------|-------------|
| **[Application Messaging (AM) Module](application-messaging.md)** | Core messaging abstraction layer for event-driven communication |
| **[Driver & Driven Architecture](architecture/driver-driven.md)** | Clean Architecture principles and dependency flow |
| **[Product Price Flow](architecture/product-price-flow.md)** | Business process documentation for pricing workflows |
| **[Basket Store Event Handlers](basket-store-event-handlers.md)** | Event handling patterns for basket and store services |

## 🏗️ Architecture Documentation

### C4 Model Diagrams
Located in `docs/architecture/c4/` - PlantUML diagrams showing system architecture at different levels:

- **System Context**: External actors and system boundaries
- **Container Diagrams**: Microservices and data flows
- **Component Diagrams**: Internal structure of key services

### How to View C4 Diagrams
1. Install VS Code PlantUML extension
2. Install Graphviz (`brew install graphviz` on macOS)
3. Open `.puml` files and preview

## 🛠️ Key Components

### Core Modules
- **`internal/am`**: Application Messaging abstraction layer
- **`internal/ddd`**: Domain-Driven Design primitives
- **`internal/es`**: Event Sourcing infrastructure
- **`internal/jetstream`**: NATS JetStream message transport
- **`internal/registry`**: Event serialization/deserialization

### Services
- **`baskets/`**: Shopping basket management
- **`customers/`**: Customer data and profiles
- **`depot/`**: Inventory and warehouse management
- **`ordering/`**: Order processing and fulfillment
- **`payments/`**: Payment processing
- **`search/`**: Product search and catalog
- **`stores/`**: Store and product management
- **`notifications/`**: Event-driven notifications

## 🚀 Getting Started

1. **Architecture Overview**: Start with [Driver & Driven Architecture](architecture/driver-driven.md)
2. **Messaging**: Learn about [Application Messaging](application-messaging.md)
3. **System Design**: Review C4 diagrams in `architecture/c4/`
4. **Business Processes**: Check specific flows like [Product Price Flow](architecture/product-price-flow.md)

## 📋 Development Guidelines

- Follow Clean Architecture principles (Driver → Core → Driven)
- Use the AM module for all inter-service communication
- Implement domain events for business logic changes
- Use event sourcing for critical business entities
- Test business logic independently of infrastructure

## 🔧 Infrastructure

- **Message Broker**: NATS JetStream
- **Database**: PostgreSQL
- **Protocol**: gRPC for service communication
- **Event Format**: Protocol Buffers
- **Deployment**: Docker Compose for local development

---

For code examples and API documentation, see the individual service modules and their protobuf definitions.
