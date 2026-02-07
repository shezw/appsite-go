# Appsite Go Framework Design & Development Guidelines

This document serves as the primary "skill" or knowledge base for the development of the Appsite Go full-stack framework. All subsequent development must adhere to these guidelines to ensure Go best practices and support for both monolithic and microservices architectures.

## 1. Architectural Overview

The framework is designed based on a Clean Architecture approach, decoupling the core business logic from transport layers and external dependencies. This ensures that the system can be deployed as a unified monolith or split into independent microservices with minimal refactoring.

### Core Principles
- **Interface-Driven Development**: All services must be defined by interfaces. This allows implementations to be direct method calls (in monolith) or RPC/HTTP clients (in microservices).
- **Standard Go Project Layout**: We follow the community-accepted directory structure.
- **Modularity**: The system is divided into distinct functional domains (Core, Services, Utils, etc.) as depicted in the design chart.

## 2. Directory Structure

The project structure maps the design nodes to a standard Go layout:

```
appsite-go/
├── cmd/                    # Entry points
│   ├── appsite-monolith/   # Single binary containing all services
│   ├── appsite-admin/      # Admin portal specific entry
│   └── services/           # Individual service entry points (for microservice deployment)
├── internal/               # Private application and library code
│   ├── core/               # "Core" node: Fundamental infrastructure
│   │   ├── model/          # Base models and interfaces
│   │   ├── route/          # Routing infrastructure
│   │   ├── setting/        # Configuration management
│   │   ├── error/          # Standardized error handling
│   │   └── log/            # Logging wrapper
│   ├── services/           # "Services" node: Business Logic Domains
│   │   ├── access/         # Operation, Permission, Token, Verify
│   │   ├── user/           # Account, Collect, Comment, Group, Info, Pocket, Preference
│   │   ├── contents/       # Category, Article, Banner, Page, Tag, Media
│   │   ├── commerce/       # Coupon, Order, Payment, Product, Shipping, Stock
│   │   ├── finance/        # Deal, Withdraw, Point
│   │   ├── form/           # Contract, Request, Verify
│   │   ├── message/        # Chat, Announcement, Notification
│   │   ├── relation/       # Friendships/Follows
│   │   ├── shieldword/     # Content filtering
│   │   └── world/          # Area, Company, District, Subway, Industry, Saas
│   ├── apis/               # "APIS" node: Public HTTP Adapters/Controllers
│   │   ├── account/
│   │   ├── content/
│   │   └── redirect/
│   └── admin/              # "Admin" node: Admin Handling Logic
│       ├── install/
│       ├── dashboard/
│       ├── user/
│       ├── contents/
│       ├── commerce/
│       └── ...
├── pkg/                    # Public library code (can be imported by other projects)
│   ├── utils/              # "Utils" node: Common utilities
│   │   ├── i18n/
│   │   ├── file/
│   │   ├── simpleimage/
│   │   ├── timeconvert/
│   │   ├── orm/            # Database abstraction
│   │   ├── graphql/
│   │   └── redis/
│   ├── thirdparty/         # "ThirdParty" node: External API integrations
│   │   ├── aliyunoss/
│   │   ├── baiduocr/
│   │   ├── google/
│   │   ├── paypal/
│   │   ├── stripe/
│   │   └── ...
│   └── extra/              # "Extra" node: Service adapters
│       ├── cloudstorage/
│       ├── sms/
│       ├── smtp/
│       └── ...
├── api/                    # OpenAPI/Swagger specs, Protocol definitions
├── configs/                # Configuration templates
└── web/                    # "WebSite" and "Admin" frontend assets
```

## 3. Development Guidelines

### Go Best Practices
- **Error Handling**: Use explicit error checking. Wrap errors with context usage.
- **Concurrency**: Use Channels and Goroutines for async tasks (e.g., sending SMS/Emails in `extra` services). Avoid shared mutable state.
- **Context**: Pass `context.Context` as the first argument to all functions in the call chain that involve I/O or long-running processes.

### Testing
- **Mandatory Unit Tests**: Every module implementation (non-interface/abstract code) must be accompanied by a corresponding `_test.go` file.
- **Coverage**: Aim for high code coverage for core logic and utilities.
- **mocks**: Use `gomock` or similar to mock interfaces when testing services dependent on other layers.

### Database (ORM)
- The `Core -> Model` and `Utils -> ORM` suggest a centralized data handling approach.
- Models should use struct tags for validation and JSON marshaling.
- Repository pattern should be used within `internal/services` to access data, keeping the database choice decoupled.

### Microservices Readiness
- **Service Interfaces**: Define service boundaries clearly in `internal/services`.
  - Example: `type UserService interface { Register(...) }`
- **Dependency Injection**: Use a wire/DI framework or manual injection in `cmd/` to wire up services.
- **Configurability**: Ensure `Core -> Setting` can load config from efficient sources (Env vars, Consol, Etcd) for diverse environments.

## 4. Frontend Architecture

### Admin Panel (Management)
- **Role**: Backend Management System for Admin Users.
- **Technology**: Single Page Application (SPA).
- **Framework**: React.js (with Ant Design Pro or generic scaffold).
- **Interaction**: Consumes RESTful APIs from `internal/admin/*`.
- **Deployment**:
  - Source code resides in `web/admin`.
  - Build artifacts (HTML/CSS/JS) are embedded into the Go binary using `embed`.
  - Served via a static file server middleware in `cmd/appsite-monolith`.

### Web Site (Content Presentation)
- **Role**: Public facing website (CMS, Blog, Portal).
- **Requirements**: High SEO (SSR), Interactivity (SPA), Configurable rendering.
- **Framework**: Next.js (React).
- **Rendering Strategy**: Application supports hybrid rendering modes configurable per page type:
  - **ISR (Incremental Static Regeneration)**: For high-traffic public content (Articles, Lists).
  - **SSR (Server-Side Rendering)**: For dynamic pages requiring SEO (User Profiles).
  - **CSR (Client-Side Rendering)**: For user-specific interactive dashboards.
- **Architecture**:
  - **Headless CMS Pattern**: Next.js acts as the "Head", fetching data from the "Headless" Go API (`internal/apis/*`).
  - **Deployment**:
    - **Mode A (Node.js)**: Deployed as a standalone Node.js server for full SSR/ISR capabilities.
    - **Mode B (Static Export)**: Exported as static HTML/JS and embedded in Go binary (loses ISR/SSR, degrade to CSR).

## 5. Module Specifications

### Core
Crucial for system stability. `Error` and `Log` packages must be established first to be used everywhere.

### Services (Business Logic)
This is the heart of the application. Business rules reside here, NOT in the HTTP handlers (`apis/` or `admin/`).
- **User**: Central identity management.
- **Commerce**: Handling complex transaction flows.

### Utils & ThirdParty
Keep these stateless. They should be initialized with configuration and safe for concurrent use.

---
*This file is to be referenced for all architectural decisions.*
