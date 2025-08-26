# Snippetbox

A snippet-sharing web application built in **Go**, following the *Let's Go* book by Alex Edwards.

This project was my introduction to building a full-featured web app in Go. It covers everything from routing and templates to user authentication and middleware. While simple in scope, it gave me a solid foundation in **idiomatic Go web development** and helped me learn the fundamentals needed for production-ready applications.

---

## 🚀 Features

- Create and view text snippets
- User signup and login with secure session management
- Flash messages for user feedback
- Middleware for authentication and logging
- Server-side HTML rendering with `html/template`
- PostgreSQL database integration
- Environment-based configuration

---

## 🧠 What I Learned

Through this project, I gained hands-on experience with:

- **Project Organization**: Structuring a Go web app for clarity and maintainability.
- **HTTP in Go**: Using the standard `net/http` package and custom handlers.
- **Templates**: Building layouts, partials, and safe HTML rendering.
- **Forms**: Parsing and validating user input.
- **Sessions & Authentication**: Implementing secure login, logout, and session handling.
- **Middleware**: Writing reusable middleware for security, error handling, and logging.
- **PostgreSQL**: Writing SQL queries directly and connecting via `database/sql`.
- **Error Handling**: Defining and returning meaningful errors at each layer.
- **Testing Basics**: Adding tests for handlers and helper functions.

---

## 🛠️ Tech Stack

- **Language**: Go (net/http, html/template, database/sql)
- **Database**: PostgreSQL
- **Frontend**: Server-side rendered templates
- **Sessions**: Secure cookie-based sessions
- **Tooling**: `go test`, `go vet`, and `golangci-lint`

---

## ⚡ Getting Started

### Prerequisites
- Go 1.22+
- PostgreSQL

### Setup

```bash
# Clone the repo
git clone https://github.com/yourusername/snippetbox.git
cd snippetbox

# Set up environment variables (DB connection, secret keys, etc.)
export SNIPPETBOX_DB_DSN="postgres://user:pass@localhost/snippetbox"
export SNIPPETBOX_SECRET_KEY="your-secret-key"

# Run migrations (schema.sql provided)
psql -d snippetbox -f ./migrations/schema.sql

# Run the server
go run ./cmd/web
