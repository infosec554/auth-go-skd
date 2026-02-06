# Auth Go SDK

![Go Version](https://img.shields.io/badge/go-1.22+-00ADD8?style=flat&logo=go)
![Docker](https://img.shields.io/badge/docker-available-2496ED?style=flat&logo=docker)
![License](https://img.shields.io/badge/license-MIT-green)

A production-ready, standalone **Authentication SDK** built with Golang. It provides a secure, flexible, and scalable RESTful API logic for handling user authentication, session management, and social logins.

Designed as a **Library** (`go-pkgz/auth` style) to be easily integrated into any Go application (Chi, Gin, Stdlib).

---

## 🌟 Features

- **Framework Agnostic**: Works with `net/http` handlers and standard middleware.
- **Provider Pattern**: Plug-and-play support for Google, GitHub, etc.
- **Token Management**: Built-in JWT generation, validation, and Cookie management.
- **Context Helper**: Easily retrieve user info in your handlers with `token.GetUserInfo(r)`.
- **Avatar Storage**: Pluggable storage for user avatars (LocalFS, AWS S3, etc.).

### 🌐 Supported Integrations (Roadmap)

Built with a **Pluggable Provider Pattern**, allowing for easy addition of new providers.

| # | Integration | Status | Type | Description |
| :--- | :--- | :--- | :--- | :--- |
| 1️⃣ | **Google OAuth** | ✅ **DONE** | OAuth 2.0 | Most popular login method. Fully implemented. |
| 2️⃣ | **GitHub OAuth** | ⏳ *Pending* | OAuth 2.0 | Essential for developer-focused tools. |
| 3️⃣ | **GitLab OAuth** | ⏳ *Pending* | OAuth 2.0 | For Enterprise / DevOps environments. |
| 4️⃣ | **LinkedIn OAuth** | ⏳ *Pending* | OAuth 2.0 | B2B & HR platforms. |
| 5️⃣ | **Facebook OAuth** | ⏳ *Pending* | OAuth 2.0 | General social media users. |
| 6️⃣ | **Twitter (X) OAuth** | ⏳ *Pending* | OAuth 2.0 | Media & Community products. |
| 7️⃣ | **Microsoft (Azure AD)**| ⏳ *Pending* | OAuth 2.0 | Corporate / Office 365 SSO. |
| 8️⃣ | **Apple Sign In** | ⏳ *Pending* | OIDC | Mandatory for iOS Apps (Privacy-first). |
| 9️⃣ | **Telegram Login** | ⏳ *Pending* | Widget | Passwordless login via Telegram Messenger. |
| 1️⃣0️⃣| **Twilio SMS OTP** | ⏳ *Pending* | OTP | Login via Phone Number (Passwordless). |
| 1️⃣1️⃣| **Email + Password** | ✅ **DONE** | Classic | Standard fallback login method. |
| 1️⃣2️⃣| **Email Magic Link** | ⏳ *Pending* | Passwordless | Secure link sent to email for one-click login. |

---

## 🚀 Installation

```bash
go get github.com/infosec554/auth-go-skd
```

---

## 💻 Usage Example

Here is how to use `auth-go-skd` in your main application:

```go
package main

import (
    "net/http"
    "time"
    
    "github.com/go-chi/chi/v5"
    "auth-go-skd/auth"
    "auth-go-skd/avatar"
    "auth-go-skd/token"
    "auth-go-skd/provider/google"
)

func main() {
    // 1. Configure the Service
    opts := auth.Opts{
        SecretReader: func(id string) (string, error) {
            return "super-secret-key", nil
        },
        TokenDuration:  time.Minute * 15,
        CookieDuration: time.Hour * 24,
        Issuer:         "my-app",
        URL:            "http://localhost:8080",
        AvatarStore:    avatar.NewLocalFS("/tmp/avatars"),
    }

    // 2. Initialize Service & Providers
    service := auth.NewService(opts)
    service.AddCustomProvider(google.New(config.Google{...}))

    // 3. Mount Handlers
    r := chi.NewRouter()
    authHandler, avatarHandler := service.Handlers()
    r.Mount("/auth", authHandler)
    r.Mount("/avatar", avatarHandler)

    // 4. Protect Routes
    m := service.Middleware()
    r.Group(func(r chi.Router) {
        r.Use(m.Auth)
        r.Get("/private", func(w http.ResponseWriter, r *http.Request) {
            user := token.MustGetUserInfo(r)
            w.Write([]byte("Hello " + user.Name))
        })
    })

    http.ListenAndServe(":8080", r)
}
```

---

## 📂 Project Structure

```
auth-go-skd/
├── auth/                  # Core Authentication Logic (Service, Handlers, Middleware)
├── provider/              # OAuth Provider Interfaces & Implementations
│   ├── google/            # Google Provider
│   └── ...                # Other providers (Github, Facebook, etc.)
├── token/                 # JWT Token Management & Context Helpers
├── avatar/                # User Avatar Storage Layer
├── store/                 # Storage Repositories (Postgres, Redis Interfaces)
├── data/                  # Core Data Models (User, Session, Identity)
├── config/                # Configuration Loader
└── cmd/                   # Example Application entry point
```

## 🤝 Contributing

Contributions are welcome! Please fork the repository and submit a Pull Request.

## 📄 License

This project is licensed under the MIT License.
