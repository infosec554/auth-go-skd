# Auth SDK with Social Login & Integrations

Welcome to the comprehensive Auth SDK built in Golang. This SDK provides a robust authentication system including classic email/password login and 12 planned social/enterprise integrations.

## ✅ Core Features (Completed)

These core functions are fully implemented, tested, and ready to use via the Demo UI.

1.  **Register** - Sign up with Email/Password.
2.  **Login** - Authenticate and receive JWT Access/Refresh Tokens.
3.  **RefreshToken** - Securely rotate tokens.
4.  **Logout** - Invalidate sessions.
5.  **GetProfile** - Retrieve authenticated user details.
6.  **UpdateProfile** - Update user information (Name).
7.  **ChangePassword** - securely change user password.
8.  **DeleteAccount** - Remove user account and associated data.
9.  **SocialLogin (Google)** - OAuth 2.0 flow for Google.

---

## 🔐 12 Integrations Roadmap

Below is the status of the 12 planned integrations. We use a **Provider Pattern** to easily extend support for new providers.

| # | Integration | Status | Type | Description |
| :--- | :--- | :--- | :--- | :--- |
| 1️⃣ | **Google OAuth** | ✅ **DONE** | OAuth 2.0 | Most popular login method. Fully functional. |
| 2️⃣ | **GitHub OAuth** | ⏳ *Pending* | OAuth 2.0 | Essential for developer-focused tools. |
| 3️⃣ | **GitLab OAuth** | ⏳ *Pending* | OAuth 2.0 | For Enterprise / DevOps environments. |
| 4️⃣ | **LinkedIn OAuth** | ⏳ *Pending* | OAuth 2.0 | B2B & HR platforms. |
| 5️⃣ | **Facebook OAuth** | ⏳ *Pending* | OAuth 2.0 | General social media users. |
| 6️⃣ | **Twitter (X) OAuth** | ⏳ *Pending* | OAuth 2.0 | Media & Community products. |
| 7️⃣ | **Microsoft (Azure AD)**| ⏳ *Pending* | OAuth 2.0 | Corporate / Office 365 SSO. |
| 8️⃣ | **Apple Sign In** | ⏳ *Pending* | OIDC | Mandatory for iOS Apps (Privacy-first). |
| 9️⃣ | **Telegram Login** | ⏳ *Pending* | Widget | Passwordless login via Telegram Messenger. |
| 🔟 | **Twilio SMS OTP** | ⏳ *Pending* | OTP | Login via Phone Number (Passwordless). |
| 1️⃣1️⃣| **Email + Password** | ✅ **DONE** | Classic | Standard fallback login method. |
| 1️⃣2️⃣| **Email Magic Link** | ⏳ *Pending* | Passwordless | Secure link sent to email for one-click login. |

---

## 🚀 Getting Started

### 1. Prerequisites
*   Go 1.22+
*   Docker & Docker Compose
*   Make (Optional, for easy commands)

### 2. Setup Environment
Create a `.env` file in the root directory:

```bash
# App
APP_NAME=auth-service
APP_VERSION=1.0.0

# Server
HTTP_PORT=8080

# Database
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=password
POSTGRES_DB=auth_db
POSTGRES_SSL_MODE=disable

# Redis
REDIS_ADDR=localhost:6379

# Google OAuth Credentials
GOOGLE_CLIENT_ID=your_google_client_id
GOOGLE_CLIENT_SECRET=your_google_client_secret
GOOGLE_REDIRECT_URL=http://localhost:8080/api/auth/google/callback
```

### 3. Run the Project
We have a Makefile to simplify everything.

**Step 1: Start Infrastructure (Postgres & Redis)**
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

### 4. Test with Demo UI
Open your browser and navigate to:
👉 **http://localhost:8080**

You will see a Demo Panel where you can test all implemented features including Google Login.

---

## 📂 Project Structure

```
/auth-go-skd
├── cmd
│   └── main.go           # Entry point
├── config
│   └── config.go         # Configuration loader
├── internal
│   ├── domain            # Core Data Models (User, Identity, Session)
│   ├── service           # Business Logic (Auth, GoogleLogin, etc.)
│   ├── storage           # Database Layer (Postgres)
│   ├── http              # API Handlers & Routes
│   └── providers         # Social Providers (Google, GitHub...)
├── migrations            # SQL Migration files
├── public                # Static files for Demo UI
├── Makefile              # Command shortcuts
└── docker-compose.yml    # Database infrastructure
```
# auth-go-skd
