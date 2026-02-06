# Auth Go SDK

![Go Version](https://img.shields.io/badge/go-1.22+-00ADD8?style=flat&logo=go)
![Docker](https://img.shields.io/badge/docker-available-2496ED?style=flat&logo=docker)
![License](https://img.shields.io/badge/license-MIT-green)

A production-ready, standalone **Authentication Service (SDK)** built with Golang. It provides a secure, flexible, and scalable RESTful API for handling user authentication, session management, and OAuth2 social logins.

Designed with **Clean Architecture** principles, making it easy to extend, test, and maintain.

---

## 🌟 Features

### 🔐 Core Authentication
- **Registration**: Secure user sign-up with bcrypt password hashing.
- **Login**: JWT-based authentication (Access & Refresh Tokens).
- **Session Management**: Secure token rotation and session tracking (Postgres/Redis backed).
- **Logout**: Safe session invalidation.

### 👤 User Management
- **Profile Management**: Retrieve and update user details.
- **Security**: Password change and account deletion functionalities.
- **Roles**: Basic role-based access control (RBAC) ready.

### 🌐 Social Integrations (OAuth2)
Built with a pluggable Provider Pattern. Currently supports:
- [x] **Google**
- [ ] GitHub (Ready to implement)
- [ ] Facebook (Ready to implement)
- [ ] ... and extensible for any OAuth2 provider.

---

## 🛠 Tech Stack

- **Language**: Go (Golang) 1.22+
- **Framework**: `chi` (Lightweight, idiomatic router)
- **Database**: PostgreSQL (with `pgx` driver)
- **Caching**: Redis (for session/rate limiting - optional)
- **Config**: `cleanenv` (YAML + Environment variables)
- **Migrations**: `golang-migrate`
- **Logging**: `uber-go/zap` (Structured logging)

---

## 🚀 Getting Started

### Prerequisites
- Go 1.22 or higher
- Docker & Docker Compose
- Make (optional, for convenience)

### 1. Installation

Clone the repository:
```bash
git clone https://github.com/infosec554/auth-go-skd.git
cd auth-go-skd
```

### 2. Configuration

Create a `.env` file in the root directory (or use `config/config.yaml`):

```bash
# App
APP_NAME=auth-service
APP_PORT=8080

# Database
POSTGRES_HOST=localhost
POSTGRES_PORT=5435
POSTGRES_USER=postgres
POSTGRES_PASSWORD=password
POSTGRES_DB=auth_db

# OAuth (Google)
GOOGLE_CLIENT_ID=your_client_id
GOOGLE_CLIENT_SECRET=your_client_secret
GOOGLE_REDIRECT_URL=http://localhost:8080/api/auth/google/callback
```

### 3. Running the Service

We provide a `Makefile` to simplify common tasks.

**Option A: Using Docker (Recommended)**
Start Postgres and Redis containers:
```bash
make docker-up
```

**Option B: Manual Run**
Ensure Postgres is running locally, then apply migrations:
```bash
make migrate-up
```
Start the server:
```bash
make run
```

The server will be available at `http://localhost:8080`.

---

## 📡 API Documentation

### Authentication

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `POST` | `/api/auth/register` | Register a new user |
| `POST` | `/api/auth/login` | Login with email/password |
| `POST` | `/api/auth/refresh` | Refresh access token |
| `POST` | `/api/auth/logout` | Logout user |

### Social Auth

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `GET` | `/api/auth/{provider}/login` | Initiate generic OAuth login (e.g., /google/login) |
| `GET` | `/api/auth/{provider}/callback` | OAuth callback handler |

### User Profile (Protected)

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `GET` | `/api/user/profile/{id}` | Get user profile |
| `PUT` | `/api/user/profile/{id}` | Update user profile |
| `PUT` | `/api/user/change-password/{id}` | Change password |
| `DELETE` | `/api/user/profile/{id}` | Delete account permanently |

---

## 📂 Project Structure

```
auth-go-skd/
├── cmd/
│   └── main.go            # Application entry point
├── config/                # Environment configuration
├── internal/
│   ├── domain/            # Entities & Business Interfaces (Pure)
│   ├── service/           # Business Logic Implementation
│   ├── storage/           # Database Repositories (Postgres/Redis)
│   ├── http/              # HTTP Handlers (REST Adapter)
│   └── providers/         # OAuth Provider Implementations (Google, etc.)
├── migrations/            # SQL Database Migrations
└── docker-compose.yml     # Container orchestration
```

## 🤝 Contributing

Contributions are welcome! Please fork the repository and submit a Pull Request.

## 📄 License

This project is licensed under the MIT License.
