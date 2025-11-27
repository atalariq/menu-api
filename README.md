# Menu Catalog API

> Live: [App](https://atalariq-menu-api.fly.dev/) | [Documentation](https://atalariq-menu-api.fly.dev/docs/index.html)

This project is a backend RESTful service built as a submission for the GDGoC Hacker Study Case.

It implements the Standard Go Project Layout and Clean Architecture (Controller-Service-Repository pattern) to ensure scalability and maintainability. The system is backed by PostgreSQL and integrates Google Gemini 2.0 Flash to provide intelligent features like auto-descriptions and context-aware menu recommendations.

## Features

### Core Functionality

- Menu Management: Full CRUD operations for menu items.
- Advanced Search & Filter: Filter by category, price range, and calories. Includes full-text search capability.
- Aggregation: Group menu items by category (supporting both count summaries and detailed lists).
- Clean Architecture: Separation of concerns between HTTP handlers, business logic, and database access.

### AI Integration (Google Gemini)

- Auto-Description: Automatically generates marketing-style descriptions for new items based on their ingredients if left empty during creation.
- Smart Recommendations: A recommendation engine that accepts natural language queries (e.g., "I need something to wake me up") and maps them to specific menu items in the database.

### Tooling

- Documentation: Auto-generated OpenAPI (Swagger) documentation.
- Task Runner: Simple workflow management using Just.

## Tech Stack

- **Language**: [Go 1.23+](https://go.dev)
- **Framework**: [Gin Web Framework](https://gin-gonic.com/)
- **Database**: [PostgreSQL](https://www.postgresql.org/)
- **ORM**: [GORM](https://gorm.io/)
- **AI SDK**: [Google Generative AI SDK](https://github.com/google/generative-ai-go)
- **Documentation**: [Swaggo](https://github.com/swaggo/swag)
- **Deployment**: [Fly.io](https://fly.io/)

## Project Structure

This project follows the standard Go layout:

```text
menu-api/
├── cmd/server/       # Application entry point
├── internal/
│   ├── controller/   # HTTP Handlers (Input parsing & validation)
│   ├── service/      # Business Logic (AI integration & core logic)
│   ├── repository/   # Database Access Layer (GORM implementation)
│   └── model/        # Domain entities & DTOs
├── docs/             # Swagger generated documentation
└── test/             # Unit tests with Mocking
```

## Getting Started

### Prerequisites

- Go 1.23 or higher
- PostgreSQL
- `just` (optional)

### Installation

1. Clone the repository:

   ```bash
   git clone [https://github.com/atalariq/menu-api.git](https://github.com/atalariq/menu-api.git)
   cd menu-api
   go mod download
   ```

2. Configuration:
   Set the required environment variables:

   ```bash
   # Copy .env.example
   cp .env.example .env

   # Edit it with your preferred text editor, mine is nvim
   nvim .env
   ```

   Or export environment variables from shell directly:

   ```bash
   # Required for AI features
   export GEMINI_API_KEY="your_google_api_key"

   # Required for PostgreSQL
   export DATABASE_URL="host=localhost user=postgres password=pass dbname=menu_api port=5432 sslmode=disable"
   ```

3. Run the application:

   Using standard Go command:

   ```bash
   # Generate Swagger docs (optional)
   swag init -g cmd/server/main.go --output docs

   # Run server
   go run cmd/server/main.go
   ```

   Or using `just`:

   ```bash
   just docs
   just run

   # Or
   just lazy
   ```

   The server will start at <http://localhost:8080>.

## API Documentation

The API comes with an interactive Swagger UI. Once the application is running, you can access it at:

<http://localhost:8080/docs/index.html>

## Testing

Unit tests focus on the Service layer to ensure business logic correctness, using mocks for Database and AI dependencies.

```bash
go test ./test/... -v

# or
just test
```
