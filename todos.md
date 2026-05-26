# Checklist for IIICPC Project

* [X] Initialized monorepo structure with Taskfile
* [X] Shared developer tooling configured

  * [ ] `.editorconfig`
  * [X] root `.gitignore`
  * [ ] pre-commit hooks
  * [ ] CI workflows
  * [ ] environment variable conventions
  * [ ] docker-compose for local services

---

# Go Server (`apps/api`)

* [X] Init Go module + project structure
* [X] Taskfile targets (`dev`, `build`, `test`, `lint`, `swagger`)
* [X] Config loading (`env` + `godotenv`)
* [X] Structured logging setup (`slog`)
* [X] Database migrations (`goose`) + first schema
* [X] `sqlc` setup + codegen verified
* [X] `pgx` connection pool wired up
* [X] Chi server initialized
* [] Request ID middleware
* [X] Logging middleware
* [] Recovery middleware
* [X] Centralized error handling strategy
* [X] Error response helpers
* [] Custom error types + mapping
* [ ] Validation layer (`go-playground/validator`)
* [ ] Swagger/OpenAPI setup
* [ ] Health check endpoint
* [ ] Graceful shutdown handling
* [ ] Unit test setup
* [ ] Integration test setup
* [] Linting (`golangci-lint`)
* [] GitHub Actions workflow for lint
* [ ] GitHub Actions workflow for tests
* [ ] Dockerfile
* [ ] First API endpoint
* [ ] Authentication strategy decision
* [ ] Rate limiting middleware
* [ ] CORS configuration

---

# Rust Sandbox (`apps/sandbox`)

* [ ] Initialize Cargo workspace/member
* [ ] Project structure setup
* [ ] `clippy` lint configuration
* [ ] `rustfmt` configuration
* [ ] Error handling strategy (`thiserror` / `anyhow`)
* [ ] Structured logging (`tracing`)
* [ ] Async runtime setup (`tokio`)
* [ ] Sandbox execution architecture decided
* [ ] Process isolation strategy
* [ ] Resource limits (CPU, memory, timeout)
* [ ] Secure code execution flow
* [ ] Temporary workspace management
* [ ] Language execution adapters

  * [ ] C++
  * [ ] Python
  * [ ] Java
* [ ] Test harness setup
* [ ] Unit tests
* [ ] Integration tests
* [ ] Docker sandbox environment
* [ ] GitHub Actions workflow
* [ ] Benchmarking setup
* [ ] API contract between Go ↔ Rust finalized

---

# Frontend (`apps/web`)

* [ ] Initialize frontend app
* [ ] Shared UI structure
* [ ] Routing setup
* [ ] API client setup
* [ ] Environment configuration
* [ ] State management decision
* [ ] Authentication flow
* [ ] Code editor integration
* [ ] Submission dashboard
* [ ] Contest pages
* [ ] Real-time updates strategy
* [ ] Error boundary setup
* [ ] Toast/notification system
* [ ] Loading/skeleton states
* [ ] ESLint + Prettier
* [ ] GitHub Actions workflow
* [ ] Responsive layout
* [ ] Dark mode
* [ ] Dockerfile

---

# Infra / DevOps

* [ ] Docker Compose setup
* [ ] PostgreSQL container
* [ ] Redis container
* [ ] Local development environment
* [ ] CI/CD pipeline
* [ ] Secrets management strategy
* [ ] Deployment strategy
* [ ] Monitoring/logging stack
* [ ] Backup strategy
* [ ] Production environment configs

---

# Judge System / Contest Engine

* [ ] Problem schema finalized
* [ ] Test case storage strategy
* [ ] Submission queue architecture
* [ ] Judge result pipeline
* [ ] Compilation pipeline
* [ ] Execution pipeline
* [ ] Verdict mapping
* [ ] Rejudge support
* [ ] Contest timer handling
* [ ] Leaderboard logic
* [ ] Anti-cheat considerations
* [ ] Scalability/load testing

---

# Documentation

* [ ] README
* [ ] Local setup guide
* [ ] Architecture diagram
* [ ] API documentation
* [ ] Contribution guidelines
* [ ] Coding standards
* [ ] Deployment documentation
