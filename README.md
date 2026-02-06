# Auth Service

This is a standalone Authentication Service built in Golang. It provides a RESTful API for managing user authentication, sessions, and social logins.

## ✅ Core Features

1.  **Register** - `POST /api/auth/register`
2.  **Login** - `POST /api/auth/login`
3.  **RefreshToken** - `POST /api/auth/refresh`
4.  **Logout** - `POST /api/auth/logout`
5.  **GetProfile** - `GET /api/user/profile/{id}`
6.  **UpdateProfile** - `PUT /api/user/profile/{id}`
7.  **ChangePassword** - `PUT /api/user/change-password/{id}`
8.  **DeleteAccount** - `DELETE /api/user/profile/{id}`
9.  **Social Login** - Google OAuth support.

## 🚀 Getting Started

### 1. Prerequisites
*   Go 1.22+
*   Docker & Docker Compose
*   Make

### 2. Setup Environment
Ensure `.env` file exists with DB and OAuth credentials.

### 3. Run the Project

**Step 1: Start Infrastructure**
```bash
make docker-up
```

**Step 2: Run Migrations**
```bash
make migrate-up
```

**Step 3: Start the Server**
```bash
make run
```

The server will start on port `8080`. You can interact with it using `curl` or Postman.

## 📂 Project Structure

```
/auth-go-skd
├── cmd
│   └── main.go           # Entry point
├── config
│   └── config.go         # Configuration loader
├── internal
│   ├── domain            # Core Data Models
│   ├── service           # Business Logic
│   ├── storage           # Database Layer
│   ├── http              # API Handlers
│   └── providers         # Social Providers
├── migrations            # SQL Migrations
└── docker-compose.yml    # Infrastructure
```
