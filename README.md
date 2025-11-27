# 🛡️ Aegis

Simple and secure user authentication and authorization system with a modern web interface.

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
- 🌐 **UI:** http://localhost
- 🔌 **API:** http://localhost/api

The single container uses **supervisor** to manage both the Go backend (port 8080 internally) and NGINX (port 80 exposed).

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

## ✨ Features

### Backend (aegis-server)
- 👤 User registration and login with HMAC-SHA256 password hashing
- 🔑 JWT token-based authentication with access and refresh tokens
- ⏱️ Configurable token expiration (default 24 hours)
- 🔄 Refresh token support
- 🎭 Role-based access control (RBAC)
- 🔐 Permission management
- 💾 SQLite database with foreign key constraints
- 🏥 Health check endpoint
- 📝 Comprehensive logging

### Frontend (aegis-ui)
- 🎨 Modern, responsive dark-themed interface
- 📊 Real-time service health monitoring
- 👥 User management dashboard with sortable tables
- 🎭 Role viewing and management
- 🔐 Permission viewing and management
- 📱 Mobile-friendly responsive design
- 🔄 Real-time data refresh
- 📋 Dynamic tab counters showing data counts

## 🌍 Environment Variables

Configure these in `.env` file or pass directly to Docker:

- `AEGIS_SERVER_PORT` - Server port (default: `8080`)
- `AEGIS_JWT_SECRET` - JWT signing secret (generates random if not set)
- `AEGIS_JWT_EXP_TIME` - JWT token expiration in minutes (default: `1440` = 24 hours)
- `AEGIS_HASH_KEY` - HMAC key for password hashing

## 📡 API Endpoints

### 🏥 Health Check

- `GET /health` - Health check endpoint (returns service status)

### 🎭 Roles

- `POST /roles` - Create a new role
- `GET /roles` - List all roles
- `GET /roles/:name` - Get role by name
- `PUT /roles/:name` - Update role description
- `DELETE /roles/:name` - Delete a role

### 🔐 Permissions

- `POST /permissions` - Create a new permission
- `GET /permissions` - List all permissions
- `GET /permissions/:name` - Get permission by name
- `PUT /permissions/:name` - Update permission description
- `DELETE /permissions/:name` - Delete a permission

### 👤 Users

- `POST /users/register` - Register a new user
- `POST /users/login` - Login and receive JWT tokens
- `POST /users/refresh` - Refresh access token
- `GET /users` - List all users
- `GET /users/:id` - Get user by ID
- `PUT /users/:id` - Update user details
- `DELETE /users/:id` - Delete a user
- `POST /users/:id/password` - Change user password
- `POST /users/:id/roles` - Add a role to a user
- `DELETE /users/:id/roles/:role` - Remove a role from a user
- `POST /users/:id/permissions` - Add a permission to a user
- `DELETE /users/:id/permissions/:permission` - Remove a permission from a user
## 🔧 Development

### Project Structure Details

```
aegis/
├── aegis-server/          # Backend Go API
│   ├── api/               # HTTP handlers
│   │   ├── user/          # User endpoints
│   │   ├── role/          # Role endpoints
│   │   └── permission/    # Permission endpoints
│   ├── domain/            # Business logic
│   │   ├── user/          # User domain model
│   │   ├── role/          # Role domain model
│   │   └── permission/    # Permission domain model
│   ├── database/          # Data persistence
│   ├── util/              # Utilities
│   │   ├── hash/          # Password hashing
│   │   └── jwt/           # JWT token management
│   ├── main.go            # Application entry point
│   ├── go.mod             # Go dependencies
│   └── Dockerfile         # Backend container (multi-stage build)
├── aegis-ui/              # Frontend web interface
│   ├── index.html         # Main HTML structure
│   ├── styles.css         # Dark theme styling
### Building for Production

**Backend:**
```bash
cd aegis-server
CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o aegis .
```

**Full Stack (Docker):**
```bash
# Build and start
sudo docker compose build
sudo docker compose up -d

# Seed with test data
./seed-data.sh

# View logs
sudo docker compose logs -f

# Stop
sudo docker compose down
```

### Docker Architecture

The project uses a **single container** approach with:
- **Supervisor** as process manager (PID 1)
- **Aegis Backend** running as user `aegis` on port 8080
- **NGINX** serving static UI files and reverse proxying `/api/*` to backend
- **SQLite Database** at `/app/aegis.db` (persisted via Docker volume at `/app/data`)

The multi-stage Dockerfile:
1. **Stage 1 (builder):** Compiles Go binary with CGO for SQLite
2. **Stage 2 (runtime):** Alpine with NGINX, supervisor, and compiled binary

### Testing

The project includes 83 comprehensive tests covering:
- API integration tests (12 tests)
- Domain logic tests (41 tests)
- Utility tests (38 tests - hash and JWT)

Test database (`aegis-test.db`) is automatically created and cleaned up during test runs.

### Building for Production

**Backend:**
```bash
cd aegis-server
go build -o aegis .
```

**Full Stack (Docker):**
```bash
docker-compose build
docker-compose up -d
```

## 🛠️ Technology Stack

- **Backend:** Go 1.25+, Gin Web Framework, SQLite
- **Frontend:** HTML5, CSS3, Vanilla JavaScript
- **Deployment:** Docker, NGINX (reverse proxy)
- **Authentication:** JWT with HMAC-SHA256

## 📝 License

This is a personal learning project. Feel free to use it as a reference or starting point for your own projects.

## 🤷 Why Aegis?

Aegis is a personal project created to learn Go while building something useful. If you find it helpful, that's great! For production use, consider established solutions like Auth0, Keycloak, or similar services.
