# Go Backend Template

This is a Go backend template for building a robust and scalable RESTful API. It comes with a set of pre-configured tools and a well-defined project structure to get you started quickly.

## Features

- **OpenAPI Specification:** Define your API using the OpenAPI 3.0 standard.
- **Code Generation:** Automatically generate server and model code from your OpenAPI specification.
- **Database Migrations:** Manage your database schema with Goose migrations.
- **Type-Safe Database Queries:** Use Jet to generate type-safe SQL query builders.
- **Configuration Management:** Easily manage your application's configuration for different environments.
- **Logging:** Structured logging with `slog`.
- **Routing:** High-performance request routing with `chi`.
- **Validation:** Request data validation with `gookit/validate`.

## Packages Used

### Dependencies

- `github.com/go-chi/chi/v5`: A lightweight, idiomatic and composable router for building Go HTTP services.
- `github.com/go-errors/errors`: Provides a simple way to create and wrap errors in Go.
- `github.com/go-jet/jet/v2`: A type-safe SQL builder for Go.
- `github.com/google/uuid`: A package for working with UUIDs.
- `github.com/gookit/validate`: A powerful and flexible data validation library for Go.
- `github.com/lib/pq`: A pure Go Postgres driver for the `database/sql` package.
- `github.com/oapi-codegen/runtime`: A collection of runtime utilities for `oapi-codegen`.
- `github.com/patrickmn/go-cache`: An in-memory key:value store/cache for Go.

### Dev Dependencies

- `github.com/go-jet/jet/v2/cmd/jet`: The command-line tool for Jet.
- `github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen`: The command-line tool for `oapi-codegen`.
- `github.com/pressly/goose/v3/cmd/goose`: The command-line tool for Goose.
- `github.com/segmentio/golines`: A Go code formatter that aligns fields in structs and other multi-line constructs.
- `golang.org/x/tools/cmd/goimports`: A tool that automatically updates your Go import lines.

## Folder Structure

```
.
├── cmd
│   ├── console
│   │   └── main.go
│   └── webserver
│       └── main.go
├── config
│   ├── app.go
│   └── provider
│       ├── cache_provider.go
│       ├── database_provider.go
│       ├── env_provider.go
│       ├── logger_provider.go
│       └── validation_provider.go
├── db
│   └── migrations
├── generated
│   ├── db
│   └── oapi
├── internal
│   ├── app
│   │   ├── app_service
│   │   ├── domain_service
│   │   ├── mutations
│   │   └── repository
│   ├── console
│   └── webserver
│       ├── handler
│       ├── middleware
│       └── webserver.go
├── scripts
└── vendor
```

- `cmd`: Entry points for your applications (e.g., webserver, console).
- `config`: Application configuration and providers for various services.
- `db/migrations`: Database migration files.
- `generated`: Generated code from `oapi-codegen` and `jet`.
- `internal`: Internal application logic, including handlers, middleware, and services.
  - `app/app_service`: This is called by the handler and is where mutations and repository functions are used. It can be used to combine and orchestrate different repositories, for example, using the category repository and product repository to get products or categories.
  - `app/domain_service`: This is for data manipulation that is not dependent on the database. It is most likely for mapping to match expected results, manipulating models, or any custom stuff that needs to be done.
  - `app/mutations`: This is for CRUD operations that use `pg` for SQL.
  - `app/repository`: This is also for CRUD operations that use `pg` for SQL.
- `scripts`: Utility scripts.
- `vendor`: Go package dependencies.

## Setup

1.  **Install Go:** Make sure you have Go 1.21 or higher installed.
2.  **Install Dependencies:** Run `go mod tidy` to download the required dependencies.
3.  **Create Environment Files:** Copy `.env.local.example` to `.env.local` and `.env.test.example` to `.env.test` and update the values as needed.
4.  **Run Migrations:** Run `./commands.sh migration:up` to run the database migrations.

## Important Commands

| Command | Description |
| --- | --- |
| `webserver` | Starts the web server. |
| `console {command}` | Runs a console command. |
| `lint` | Lints the codebase. |
| `lint:fix` | Lints the codebase and applies automated fixes. |
| `format` | Formats the codebase. |
| `test` | Runs the test suite. |
| `openapi:codegen` | Generates server and model code from `openapi.yaml`. |
| `migration:codegen {name}` | Creates a new empty SQL migration file. |
| `migration:up` | Runs all pending database migrations and regenerates the type-safe database models. |
| `migration:down` | Rolls back the last database migration and regenerates the models. |
| `migration:reset` | Rolls back all database migrations. |
| `migration:status` | Shows the status of all migrations. |

## Creating a New Endpoint

1.  **Define the Endpoint in `openapi.yaml`:** Add the new endpoint to the `paths` section of `openapi.yaml`.

    ```yaml
    paths:
      /my-new-endpoint:
        get:
          summary: "My new endpoint"
          responses:
            "200":
              description: "A successful response"
              content:
                application/json:
                  schema:
                    type: object
                    properties:
                      message:
                        type: string
    ```

2.  **Generate the Code:** Run `./commands.sh openapi:codegen` to generate the server and model code.

3.  **Implement the Handler:** Create a new handler function in the `internal/webserver/handler` package.

    ```go
    package handler

    import (
    	"net/http"

    	"github.com/labstack/echo/v4"
    )

    func (h *Handler) GetMyNewEndpoint(ctx echo.Context) error {
    	return ctx.JSON(http.StatusOK, map[string]string{"message": "Hello from my new endpoint!"})
    }
    ```

4.  **Register the Route:** The route is automatically registered by `oapi-codegen`.

5.  **Start the Server:** Run `./commands.sh webserver` to start the server and test your new endpoint.
