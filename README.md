# Menu Catalog API

> This project was built as a submission for the **GDGoC Hacker Study Case**.

Menu API is a backend RESTful service designed to manage restaurant menu catalogs. This project implements the **Clean Architecture** standard (Controller-Service-Repository) and integrates **Google Gemini AI** for intelligent features.

## Features

* **Menu CRUD:** Create, Read, Update, and Delete menu data with validation.
* **Advanced Querying:** Filtering, Searching, Sorting, and Pagination.
* **Grouping & Aggregation:** Group menus by category (Count & List modes).
* **AI-Powered:**
  * **Auto Description:** Generate appetizing menu descriptions using Gemini AI.
  * **Menu Recommendation:** Contextual menu recommendations based on user mood/requests.
* **API Documentation:** Automatic documentation using Swagger UI.

## Tech Stack

* Language: [Go](https://go.dev)
* Framework: [Gin Web Framework](https://gin-gonic.com/)
* Database: [SQLite](https://sqlite.org/) (via [GORM](https://gorm.io/))
* AI SDK: Google Generative AI SDK ([`genai`](https://github.com/google/generative-ai-go))
  * Model: Gemini 2.0 Flash
* Documentation: [Swaggo](https://github.com/swaggo/swag) (Swagger v2)
* Tooling: [Just](https://github.com/casey/just) (Task Runner)

## Project Structure

```text
.
├── cmd/server/       # Entry point aplikasi
├── internal/
│   ├── controller/   # HTTP Handlers
│   ├── service/      # Business Logic & AI Integration
│   ├── repository/   # Database Access Layer
│   └── model/        # Database Structs
├── docs/             # Swagger generated files
└── test/             # Unit Tests dengan Mocking
```

## Local Development

1. Clone & Install

```bash
git clone [https://github.com/atalariq/menu-api.git](https://github.com/atalariq/menu-api.git)
cd menu-api
go mod download
```

2. Setup Environment Variable

Create a .env file (or export via terminal) for the Gemini API Key:

```bash
export GEMINI_API_KEY="isi_api_key_gemini_anda_disini"
```

3. Run the Application

Using `just`:

```bash
just run
```

Or using `go run` manually:

```bash
go run cmd/server/main.go
```

Access the API at:

\> <http://localhost:8080>

### API Documentation

Once the server is running, open your browser and access Swagger UI
\> <http://localhost:8080/swagger/index.html>
