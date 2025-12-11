> This is a personal learning project. You’re welcome to use it as a reference or as a starting point for your own ideas. However, I strongly recommend using more mature and stable solutions available in the market rather than relying on this project for production needs.

# 🛡️ Aegis

**A authentication and authorization service with JWT token management, OAuth 2.0 token introspection, and token revocation.**

Aegis is designed as a centralized authentication provider that issues and manages JWT tokens for client applications. It provides secure user management, role-based access control (RBAC), and a modern web interface for administration.

## ✨ Features

### 🔐 Authentication & Authorization
- **User Registration & Login** - Secure user registration with password hashing (HMAC-SHA256 with salt/pepper)
- **JWT Token Management** - Access and refresh tokens with embedded roles and permissions
- **Token Validation** - Server-side token validation endpoint for client applications
- **OAuth 2.0 Token Introspection** - RFC 7662 compliant introspection endpoint
- **Token Revocation** - Blacklist-based token revocation for logout and security incidents
- **Role-Based Access Control** - Flexible roles and permissions system

### 👥 User Management
- Complete CRUD operations for users
- Assign/remove roles and permissions
- Password change functionality
- User listing with full details

### 🎭 Roles & Permissions
- Independent role and permission entities
- Assign multiple roles to users
- Assign multiple permissions to users
- Flexible authorization model

### 🖥️ Web Interface
- Modern dark-themed admin interface
- User management dashboard
- Role and permission management
- Real-time token validation testing

### 🏗️ Infrastructure
- Docker containerization with NGINX reverse proxy
- SQLite database with automatic migrations
- Health check endpoint
- Database seeding script for testing
- Supervisor process management

## 🏗️ Project Structure

```
aegis/
├── aegis-server/          # Backend Go API
│   ├── api/               # API endpoints
│   ├── domain/            # Domain models
│   ├── database/          # Database layer
│   ├── util/              # Utilities (hash, jwt)
│   └── main.go
├── aegis-ui/              # Frontend web interface
│   ├── index.html
│   ├── styles.css
│   └── app.js
├── config/                # Configuration files
│   ├── aegis.env          # Environment variables
│   ├── nginx.conf         # NGINX reverse proxy config
│   ├── supervisord.conf   # Supervisor configuration
│   └── seed-data.sh       # Database seeding script
├── docker-compose.yml     # Docker compose setup
└── Dockerfile             # Container image definition
```

## 🚀 Quick Start

### 🐳 Using Docker Compose (Recommended)

The easiest way to run Aegis is using Docker Compose, which runs everything in a single container with NGINX serving the UI and reverse proxying API requests to the Go backend.

```bash
# Clone and navigate to the project
cd aegis

# Start the application
sudo docker compose up -d

# View logs
sudo docker compose logs -f aegis

# Stop the application
sudo docker compose down
```

**Access the application:**
- 🌐 **UI:** http://localhost:3100
- 🔌 **API:** http://localhost:3100/api/aegis

The single container uses **supervisor** to manage both the Go backend (port 8080 internally) and NGINX (port 3100 exposed).

### ⚠️ Important: Database Persistence

**Database Location:** The database is stored at `/app/data/aegis.db` inside the container and persists via Docker volume `aegis-data`.

> **⚠️ Windows Compatibility Warning:**  
> The database path is hardcoded to `/app/data/aegis.db` (Unix-style path) which **will not work on native Windows** without modifications. If you need to run this on Windows:
> - Use **Docker Desktop** with WSL2 backend (recommended)
> - Or modify `aegis-server/database/database.go` to use Windows-compatible paths
> - Set the `AEGIS_DB_PATH` environment variable to a Windows path (e.g., `C:\aegis\data\aegis.db`)

This ensures data persists across container restarts on Linux/macOS and Docker Desktop with WSL2.

### 🌱 Seeding Test Data

After starting the container, populate the database with sample data:

```bash
./config/seed-data.sh
```

This creates:
- **6 Permissions:** read:users, write:users, delete:users, read:reports, write:reports, manage:system
- **4 Roles:** admin, manager, viewer, analyst
- **5 Users:** alice@aegis.com, bob@aegis.com, carol@aegis.com, david@aegis.com, eve@aegis.com

All test users have the password: `Password123!`

### 🏃 Running Locally (Development)

For local development without Docker:

**1. Start the Backend:**
```bash
cd aegis-server

# Build and run
go build .
./aegis
```

The backend will start on `http://localhost:8080` by default.

**2. Serve the Frontend:**

Open `aegis-ui/app.js` and modify the `API_BASE` constant to point to your local backend:

```javascript
// Change this line at the top of app.js:
const API_BASE = 'http://localhost:8080'; // For local development
```

Then serve the UI with any static file server:

```bash
cd aegis-ui

# Using Python
python -m http.server 3000

# Using Node.js http-server
npx http-server -p 3000

# Using PHP
php -S localhost:3000
```

Access the UI at `http://localhost:3000`

**3. Seed Test Data (Local):**

```bash
# Set API_BASE for local backend
API_BASE=http://localhost:8080/api ./config/seed-data.sh
```

**Note:** Remember to revert the `API_BASE` in `app.js` back to `/api` before building for Docker.

## 🌍 Environment Variables

Configure these in `.env` file or pass directly to Docker:

- `AEGIS_SERVER_PORT` - Server port (default: `8080`)
- `AEGIS_JWT_SECRET` - JWT signing secret (generates random if not set)
- `AEGIS_JWT_EXP_TIME` - JWT token expiration in minutes (default: `1440` = 24 hours)
- `AEGIS_HASH_KEY` - HMAC key for password hashing
- `AEGIS_DB_PATH` - Database file path (default: `/app/data/aegis.db`)

## 📡 API Endpoints

### 🔐 Authentication & Token Management
- `POST /aegis/api/auth/validate` - Validate JWT token and retrieve user claims
- `POST /aegis/api/auth/introspect` - OAuth 2.0 token introspection (RFC 7662)
- `POST /aegis/api/auth/revoke` - Revoke a JWT token before expiration

### 👤 User Management
- `POST /aegis/aegis/users/register` - Register a new user
- `POST /aegis/aegis/users/login` - User login (returns JWT tokens)
- `POST /aegis/aegis/users/refresh` - Refresh access token
- `PUT /aegis/aegis/users/:id/password` - Change user password
- `GET /aegis/aegis/users` - List all users
- `GET /aegis/aegis/users/:id` - Get user by ID
- `PUT /aegis/aegis/users/:id` - Update user
- `DELETE /aegis/aegis/users/:id` - Delete user

### 🎭 Roles

- `POST /aegis/aegis/roles` - Create a new role
- `GET /aegis/aegis/roles` - List all roles
- `GET /aegis/aegis/roles/:id` - Get role by ID
- `PUT /aegis/aegis/roles/:id` - Update role
- `DELETE /aegis/aegis/roles/:id` - Delete role
- `POST /aegis/aegis/users/:userId/roles/:roleId` - Assign role to user
- `DELETE /aegis/aegis/users/:userId/roles/:roleId` - Remove role from user

### 🔑 Permission Management

- `POST /aegis/aegis/permissions` - Create a new permission
- `GET /aegis/aegis/permissions` - List all permissions
- `GET /aegis/aegis/permissions/:id` - Get permission by ID
- `PUT /aegis/aegis/permissions/:id` - Update permission
- `DELETE /aegis/aegis/permissions/:id` - Delete permission
- `POST /aegis/aegis/users/:userId/permissions/:permissionId` - Assign permission to user
- `DELETE /aegis/aegis/users/:userId/permissions/:permissionId` - Remove permission from user

### 🏥 System

- `GET /aegis/health` - Service health check

## 📖 API Examples

### Token Validation

The token validation endpoint allows client applications to validate JWT tokens server-side without sharing the JWT secret.

**Validate a valid token:**

```bash
# First, login to get a token
TOKEN=$(curl -s -X POST http://localhost:3100/api/aegis/users/login \
  -H "Content-Type: application/json" \
  -d '{"subject":"admin@example.com","password":"Password123!"}' \
  | jq -r '.access_token')

# Validate the token
curl -X POST http://localhost:3100/api/aegis/auth/validate \
  -H "Content-Type: application/json" \
  -d "{\"token\":\"$TOKEN\"}"
```

Response (valid token):
```json
{
  "valid": true,
  "user": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "subject": "admin@example.com",
    "roles": ["admin"],
    "permissions": ["read:users", "write:users", "manage:system"]
  },
  "expires_at": "2025-11-29T10:00:00Z"
}
```

**Validate an invalid token:**

```bash
curl -X POST http://localhost/api/aegis/auth/validate \
  -H "Content-Type: application/json" \
  -d '{"token":"invalid.token.here"}'
```

Response (invalid token):
```json
{
  "valid": false,
  "error": "malformed token"
}
```

**Integration example:**

```bash
# Complete flow: login → validate → use authenticated endpoint
# Step 1: Login
TOKEN=$(curl -s -X POST http://localhost/api/aegis/users/login \
  -H "Content-Type: application/json" \
  -d '{"subject":"admin@example.com","password":"Password123!"}' \
  | jq -r '.access_token')

# Step 2: Validate token
VALID=$(curl -s -X POST http://localhost/api/aegis/auth/validate \
  -H "Content-Type: application/json" \
  -d "{\"token\":\"$TOKEN\"}" \
  | jq -r '.valid')

# Step 3: Check validation result
if [ "$VALID" = "true" ]; then
  echo "Token is valid, proceeding..."
  # Make authenticated API calls
  curl -X GET http://localhost/api/aegis/users \
    -H "Authorization: Bearer $TOKEN"
else
  echo "Token is invalid, please login again"
fi
```

### Token Revocation

The token revocation endpoint allows invalidating JWT tokens before their natural expiration. This is critical for logout functionality and security incidents.

**Revoke a token:**

```bash
# Login to get a token
TOKEN=$(curl -s -X POST http://localhost/api/aegis/users/login \
  -H "Content-Type: application/json" \
  -d '{"subject":"admin@example.com","password":"Password123!"}' \
  | jq -r '.access_token')

# Revoke the token
curl -X POST http://localhost/api/aegis/auth/revoke \
  -H "Content-Type: application/json" \
  -d "{\"token\":\"$TOKEN\"}"
```

Response:
```json
{
  "success": true,
  "message": "Token revoked successfully"
}
```

**Verify token is revoked:**

```bash
# Try to validate the revoked token
curl -X POST http://localhost/api/aegis/auth/validate \
  -H "Content-Type: application/json" \
  -d "{\"token\":\"$TOKEN\"}"
```

Response (revoked token):
```json
{
  "valid": false,
  "error": "token revoked"
}
```

**Complete logout flow:**

```bash
# Step 1: Login
TOKEN=$(curl -s -X POST http://localhost/api/aegis/users/login \
  -H "Content-Type: application/json" \
  -d '{"subject":"user@example.com","password":"Password123!"}' \
  | jq -r '.access_token')

# Step 2: Use the token for API calls
curl -X GET http://localhost/api/aegis/users \
  -H "Authorization: Bearer $TOKEN"

# Step 3: Logout by revoking the token
curl -X POST http://localhost/api/aegis/auth/revoke \
  -H "Content-Type: application/json" \
  -d "{\"token\":\"$TOKEN\"}"

# Step 4: Token is now invalid (verification)
curl -X POST http://localhost/api/aegis/auth/validate \
  -H "Content-Type: application/json" \
  -d "{\"token\":\"$TOKEN\"}"
# Returns: {"valid": false, "error": "token revoked"}
```

**How it works:**

1. **JTI Claim**: Each JWT token includes a unique JTI (JWT ID) claim using UUID
2. **Blacklist Storage**: Revoked token JTIs are stored in an in-memory blacklist (thread-safe)
3. **Automatic Validation**: `/api/auth/validate` and `/api/auth/introspect` automatically check the blacklist
4. **Automatic Cleanup**: A background job runs hourly to remove expired blacklist entries
5. **Expiration**: Blacklist entries are removed after the token's natural expiration time

**Security considerations:**

- Revoked tokens remain blacklisted until their original expiration time
- The blacklist is checked on every validation/introspection request
- In-memory implementation suitable for single-instance deployments
- For production with multiple instances, consider Redis-based blacklist (future enhancement)

### User Registration and Login

**Register a new user:**

```bash
curl -X POST http://localhost/api/aegis/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "subject": "newuser@example.com",
    "password": "SecurePass123!",
    "roles": ["viewer"],
    "permissions": ["read:reports"]
  }'
```

**Login:**

```bash
curl -X POST http://localhost/api/aegis/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "subject": "admin@example.com",
    "password": "Password123!"
  }'
```

Response:
```json
{
  "user": {
    "id": "uuid",
    "subject": "admin@example.com",
    "roles": [{"name": "admin"}],
    "permissions": [{"name": "manage:system"}]
  },
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_at": "2025-11-29T10:00:00Z",
  "refresh_expires_at": "2025-11-29T10:01:00Z"
}
```

**Refresh token:**

```bash
curl -X POST http://localhost/api/aegis/users/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }'
```

### Role Management

**Create a role:**

```bash
curl -X POST http://localhost/api/aegis/roles \
  -H "Content-Type: application/json" \
  -d '{
    "name": "moderator",
    "description": "Can moderate content"
  }'
```

**List all roles:**

```bash
curl -X GET http://localhost/api/aegis/roles
```

**Assign role to user:**

```bash
curl -X POST http://localhost/api/aegis/users/{user-id}/roles \
  -H "Content-Type: application/json" \
  -d '{
    "role": "moderator"
  }'
```

### Response Codes

| Status Code | Meaning | When |
|-------------|---------|------|
| 200 OK | Success | Request completed successfully |
| 201 Created | Created | Resource created successfully |
| 400 Bad Request | Invalid request | Missing required fields or validation error |
| 401 Unauthorized | Unauthorized | Invalid credentials or missing authentication |
| 404 Not Found | Not found | Resource doesn't exist |
| 409 Conflict | Conflict | Resource already exists |
| 500 Internal Server Error | Server error | Unexpected server-side error |

**Note**: The `/api/auth/validate` endpoint always returns 200 OK, using the `valid` field to indicate token validity.

### Token Introspection (RFC 7662)

The token introspection endpoint provides detailed information about JWT tokens following the OAuth 2.0 Token Introspection standard (RFC 7662).

**Introspect an active token:**

```bash
# First, login to get a token
TOKEN=$(curl -s -X POST http://localhost/api/aegis/users/login \
  -H "Content-Type: application/json" \
  -d '{"subject":"admin@example.com","password":"Password123!"}' \
  | jq -r '.access_token')

# Introspect the token
curl -X POST http://localhost/api/aegis/auth/introspect \
  -H "Content-Type: application/json" \
  -d "{\"token\":\"$TOKEN\"}"
```

Response (active token):
```json
{
  "active": true,
  "scope": "role:admin read:users write:users manage:system",
  "client_id": "aegis-default-client",
  "username": "admin@example.com",
  "token_type": "Bearer",
  "exp": 1732723200,
  "iat": 1732636800,
  "sub": "123e4567-e89b-12d3-a456-426614174000",
  "iss": "aegis",
  "roles": ["admin"],
  "permissions": ["read:users", "write:users", "manage:system"]
}
```

**Introspect an inactive token:**

```bash
curl -X POST http://localhost/api/aegis/auth/introspect \
  -H "Content-Type: application/json" \
  -d '{"token":"expired.or.invalid.token"}'
```

Response (inactive token):
```json
{
  "active": false
}
```

**With token type hint:**

```bash
# Specify token type for potential optimization
curl -X POST http://localhost/api/aegis/auth/introspect \
  -H "Content-Type: application/json" \
  -d "{\"token\":\"$TOKEN\",\"token_type_hint\":\"access_token\"}"
```

**OAuth 2.0 RFC 7662 Compliance:**

The introspection endpoint follows the OAuth 2.0 Token Introspection specification:
- Returns `active: false` for any invalid, expired, or malformed token
- Provides detailed token metadata for active tokens
- Supports optional `token_type_hint` parameter
- Uses standard RFC 7662 claim names (sub, iat, exp, iss, scope)
- Builds scope string from roles and permissions in format: `role:admin permission:read`

**Integration example:**

```bash
# API gateway using introspection to validate requests
TOKEN=$(curl -s -X POST http://localhost/api/aegis/users/login \
  -H "Content-Type: application/json" \
  -d '{"subject":"admin@example.com","password":"Password123!"}' \
  | jq -r '.access_token')

# Introspect token to get authorization info
INTROSPECT=$(curl -s -X POST http://localhost/api/aegis/auth/introspect \
  -H "Content-Type: application/json" \
  -d "{\"token\":\"$TOKEN\"}")

ACTIVE=$(echo $INTROSPECT | jq -r '.active')
SCOPE=$(echo $INTROSPECT | jq -r '.scope')

if [ "$ACTIVE" = "true" ]; then
  echo "Token is active"
  echo "Scope: $SCOPE"
  # Make authorized API calls based on scope
else
  echo "Token is inactive"
fi
```

## 🔧 Development & Deployment

### Running Tests

```bash
cd aegis-server

# Run all tests
go test ./...

# Run with coverage
go test ./... -cover

# Run specific package tests
go test ./api/auth -v
go test ./domain/token -v
```

**Test Coverage:**
- `api/auth`: 92.8% (33 tests - validation, introspection, revocation)
- `domain/token`: 100% (10 tests - blacklist operations)
- `domain/user`: 24.0% (17 tests)
- `domain/role`: 6.2% (12 tests)
- `domain/permission`: 6.2% (12 tests)
- `util/jwt`: 90.4% (19 tests)
- `util/hash`: 92.0% (19 tests)

The project includes comprehensive test coverage with over 120 tests covering:
- Authentication endpoints (validate, introspect, revoke)
- Token blacklist operations (concurrent access, cleanup)
- User management (CRUD, roles, permissions)
- JWT token generation and validation
- Password hashing and verification

Test database (`aegis-test.db`) is automatically created and cleaned up during test runs.

### Code Architecture

**Backend Structure:**
```
aegis-server/
├── api/              # HTTP handlers (Gin)
│   ├── auth/         # Token validation, introspection, revocation
│   ├── user/         # User management endpoints
│   ├── role/         # Role management endpoints
│   └── permission/   # Permission management endpoints
├── domain/           # Business logic and domain models
│   ├── token/        # Token blacklist system
│   ├── user/         # User entity and service
│   ├── role/         # Role entity and service
│   └── permission/   # Permission entity and service
├── database/         # Database initialization and migrations
└── util/             # Shared utilities
    ├── jwt/          # JWT token generation and validation
    └── hash/         # Password hashing (HMAC-SHA256)
```

**Design Patterns:**
- **Layered Architecture**: Clear separation (API → Domain → Database)
- **Repository Pattern**: Database access abstraction
- **Service Layer**: Business logic in domain packages
- **Interface-based Design**: Blacklist storage abstraction (in-memory/Redis-ready)

**Docker Architecture:**
- **Single Container**: Supervisor manages Go backend (port 8080) + NGINX (port 80)
- **Multi-stage Build**: Builder stage (Go compilation) + Runtime stage (Alpine + NGINX)
- **Database**: SQLite at `/app/data/aegis.db` (persisted via Docker volume)

### Building for Production

**Backend (Local):**
```bash
cd aegis-server
go build -o aegis .
```

**Docker (Recommended):**
```bash
# Build and start
sudo docker compose build
sudo docker compose up -d

# Seed test data
./config/seed-data.sh

# View logs
sudo docker compose logs -f

# Stop
sudo docker compose down
```

## 🛠️ Technology Stack

### Backend
- **Language:** Go 1.25+
- **Web Framework:** Gin
- **Database:** SQLite with foreign key constraints
- **Authentication:** JWT (HMAC-SHA256), Token blacklist with cleanup
- **Password Hashing:** HMAC-SHA256 with salt and pepper
- **Concurrency:** Thread-safe token blacklist (sync.RWMutex)

### Frontend
- **HTML5** with semantic markup
- **CSS3** with custom dark theme
- **Vanilla JavaScript** (no frameworks)
- **Responsive Design** (mobile-friendly)

### Infrastructure
- **Containerization:** Docker with multi-stage builds
- **Web Server:** NGINX (reverse proxy + static files)
- **Process Management:** Supervisor
- **Database Persistence:** Docker volumes

### Standards Compliance
- **RFC 7662**: OAuth 2.0 Token Introspection
- **JWT**: JSON Web Tokens for stateless authentication
- **RESTful API**: Standard HTTP methods and status codes

## 🔒 Security Features

- ✅ **Password Hashing**: HMAC-SHA256 with salt and pepper
- ✅ **JWT Tokens**: Signed tokens with expiration
- ✅ **Token Revocation**: Blacklist-based with JTI claims
- ✅ **Automatic Cleanup**: Hourly removal of expired blacklist entries
- ✅ **Thread-Safe Operations**: Concurrent access protection
- ✅ **Input Validation**: Request validation at API layer
- ✅ **HTTPS Ready**: NGINX reverse proxy configuration

## 🤷 About Aegis

Aegis is a production-ready authentication and authorization service built with Go. It provides:

- **Centralized Authentication**: Single source of truth for user authentication
- **JWT Token Management**: Issue, validate, introspect, and revoke tokens
- **OAuth 2.0 Compatible**: RFC 7662 token introspection endpoint
- **Flexible RBAC**: Role-based access control with permissions
- **Admin Interface**: Modern web UI for user and permission management

Perfect for microservices architectures, API gateways, or any system requiring centralized authentication.

**Client Integration:**

Client applications can integrate with Aegis by:
1. Directing users to Aegis for login
2. Receiving JWT tokens after successful authentication
3. Validating tokens using `/api/auth/validate` or `/api/auth/introspect`
4. Implementing their own authorization middleware based on token claims
5. Revoking tokens on logout using `/api/auth/revoke`

