# Go Server

## This service contains following things

- go-chi/v5 Router
- structured environment variable loading with validation
- go-slog logger: with different config in the dev and prod
- sqlc(query + migration) + pgxpool for the database
- swagger with swag and http-swag and with comment annotation
- structured api error handling with error middleware
- structured api response
- module driven approach for the api's
- go-playground/validator for validator the requests and the single validator middleware to validate

```bash
.
├── Taskfile.yml
├── api
│   ├── dto. # shared dto
│   ├── middleware
│   │   ├── error-middleware
│   │   │   └── error.go
│   │   ├── request-logger
│   │   │   └── request-logger.go
│   │   └── validator
│   │       └── validator.go
│   └── modules                  
│       ├── auth
│       │   ├── auth.dto.go
│       │   ├── auth.handler.go
│       │   ├── auth.repository.go
│       │   └── auth.routes.go
│       └── route.go
├── cmd
│   └── server
│       └── main.go # main entry point
├── docs
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── go.mod
├── go.sum
├── internal
│   ├── app
│       └── app.go 
├── logs
│   └── app.log
├── pkg
│   ├── apierr
│   │   └── errors.go
│   ├── apiresponse
│   │   └── api-response.go
│   ├── config
│   │   └── config.go
│   ├── database
│   │   └── postgres.go
│   ├── logger
│   │   └── logger.go
│   ├── repository
│   │   └── gen-queries
│   │       ├── db.go
│   │       ├── models.go
│   │       └── query.sql.go
│   ├── shared
│   │   └── shared-main.go
│   └── validator
├── readme.md
├── sql
│   ├── migrations
│   │   └── 20260426170934_init_schema.sql
│   └── queries
│       └── users
│           └── query.sql
└── sqlc.yaml
```