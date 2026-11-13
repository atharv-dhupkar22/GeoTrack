
# GeoTrack – Setup and Installation Guide

This document explains how to configure, install, run, test, and troubleshoot the GeoTrack application on Windows.

## 1. Prerequisites

Install the following software:

| Software | Purpose |
|---|---|
| Git | Version control |
| Go | Backend development |
| Node.js and npm | Frontend development |
| Docker Desktop | Containerization |
| Visual Studio Code | Code editor |
| PowerShell | Windows terminal |

For Docker deployment, Docker Desktop with Docker Compose is sufficient to run the application.

## 2. Clone the Repository

Open PowerShell:

```powershell
git clone YOUR_GITHUB_REPOSITORY_URL

cd GeoTrack
```

Replace `YOUR_GITHUB_REPOSITORY_URL` with the actual repository URL.

## 3. Configure Environment Variables

### Docker environment

Create a local `.env` file from the example:

```powershell
Copy-Item .env.example .env
```

Edit `.env` and configure the database credentials:

```env
POSTGRES_USER=geotrack
POSTGRES_PASSWORD=replace_with_a_secure_password
POSTGRES_DB=geotrack

APP_NAME=GeoTrack
APP_ENV=development
APP_PORT=8080
FRONTEND_URL=http://localhost:3000
```

Use the same database credentials consistently.

Never commit the actual `.env` file.

### Backend local development environment

Create `backend/.env`:

```env
APP_NAME=GeoTrack
APP_ENV=development
APP_PORT=8080
FRONTEND_URL=http://localhost:5173
DATABASE_URL=postgres://geotrack:YOUR_PASSWORD@127.0.0.1:5432/geotrack?sslmode=disable
```

Replace `YOUR_PASSWORD` with the password configured for PostgreSQL.

### Frontend local development environment

Create `frontend/.env`:

```env
VITE_API_BASE_URL=http://localhost:8080
VITE_WS_URL=ws://localhost:8080/ws/alerts
```

These frontend variables are used by Vite when building and running the frontend.

## 4. Start with Docker Compose

Docker Compose runs PostgreSQL/PostGIS, the Go backend, and the React frontend.

From the project root:

```powershell
docker compose up --build -d
```

Check the running containers:

```powershell
docker compose ps
```

View backend logs:

```powershell
docker compose logs -f backend
```

View database logs:

```powershell
docker compose logs -f postgres
```

View frontend logs:

```powershell
docker compose logs -f frontend
```

Stop the application:

```powershell
docker compose down
```

The named PostgreSQL volume is retained when using this command.

Do not use `docker compose down -v` unless you intentionally want to delete the database volume and its data.

## 5. Database Setup and Migrations

GeoTrack uses PostgreSQL with PostGIS spatial functionality.

The migration files are located in:

```text
backend/migrations/
    001_init.sql
    002_create_alerts.sql
```

The first migration creates the core application tables. The second creates the alerts table.

### Apply migrations to a fresh Docker database

Start PostgreSQL:

```powershell
docker compose up -d postgres
```

Wait until PostgreSQL is healthy:

```powershell
docker compose ps
```

Copy the migration files into the database container:

```powershell
docker compose cp backend/migrations/001_init.sql geotrack-postgres:/tmp/001_init.sql

docker compose cp backend/migrations/002_create_alerts.sql geotrack-postgres:/tmp/002_create_alerts.sql
```

Apply them in order:

```powershell
docker compose exec postgres psql -U geotrack -d geotrack -f /tmp/001_init.sql

docker compose exec postgres psql -U geotrack -d geotrack -f /tmp/002_create_alerts.sql
```

If you changed the database username or database name in `.env`, replace `geotrack` in these commands accordingly.

Verify the tables:

```powershell
docker compose exec postgres psql -U geotrack -d geotrack -c "\dt"
```

**Important:** These migration commands are for a fresh database. Do not rerun them against an already initialized database unless the migration scripts have been made safely repeatable. The development database used during project verification already has both migrations applied.

## 6. Run Backend Locally

Make sure PostgreSQL is running and `backend/.env` contains the correct database connection string.

Open a PowerShell terminal:

```powershell
cd backend

go mod download

go run ./cmd
```

The backend starts on:

```text
http://localhost:8080
```

Expected startup messages include:

```text
PostgreSQL connection established successfully
GeoTrack backend starting on port 8080
Environment: development
```

## 7. Run Frontend Locally

Open a second PowerShell terminal:

```powershell
cd frontend

npm install

npm run dev
```

The Vite development server is available at:

```text
http://localhost:5173
```

Ensure the backend is running before using the frontend.

## 8. Application URLs

### Docker deployment

| Service | URL |
|---|---|
| React frontend | http://localhost:3000 |
| Go backend | http://localhost:8080 |
| Health endpoint | http://localhost:8080/health |
| WebSocket | ws://localhost:8080/ws/alerts |
| PostgreSQL | localhost:5432 |

### Frontend development mode

| Service | URL |
|---|---|
| Vite frontend | http://localhost:5173 |
| Go backend | http://localhost:8080 |

## 9. API Reference

| Method | Endpoint | Description |
|---|---|---|
| GET | /health | Check backend health |
| POST | /geofences | Create a geofence |
| GET | /geofences | Retrieve geofences |
| POST | /vehicles | Register a vehicle |
| GET | /vehicles | Retrieve vehicles |
| POST | /vehicles/location | Submit vehicle GPS location |
| GET | /vehicles/location/{vehicle_id} | Retrieve latest vehicle location |
| POST | /alerts/configure | Configure geofence alerts |
| GET | /alerts/configurations | Retrieve alert configurations |
| GET | /alerts | Retrieve notification history |
| GET | /violations/history | Retrieve violation history |
| WS | /ws/alerts | Receive real-time alert events |

API responses include a `time_ns` field.

For request body formats and query parameters, refer to the backend handlers and models.

## 10. Testing and Validation

### Backend unit tests

```powershell
cd backend

go test ./... -v
```

### Backend static analysis

```powershell
go vet ./...
```

### Backend build

```powershell
go build ./...
```

### Frontend production build

```powershell
cd frontend

npm install

npm run build
```

### Docker validation

```powershell
cd ..

docker compose config

docker compose ps
```

### Health endpoint test

```powershell
Invoke-RestMethod http://localhost:8080/health |
    ConvertTo-Json -Depth 5
```

## 11. Troubleshooting

### PostgreSQL connection error

Check whether the database container is running:

```powershell
docker compose ps
```

Inspect the database logs:

```powershell
docker compose logs postgres
```

Confirm that the database credentials and `DATABASE_URL` match.

### Backend is not starting

Inspect backend logs:

```powershell
docker compose logs backend
```

Check that PostgreSQL is healthy and that the backend environment variables are configured correctly.

### Frontend cannot connect to the backend

Verify:

- Backend is available at port 8080.
- `VITE_API_BASE_URL` is correct.
- `VITE_WS_URL` points to the correct WebSocket endpoint.
- CORS permits the frontend origin.

If you change frontend build arguments, rebuild the frontend image:

```powershell
docker compose up --build -d frontend
```

### Port already in use

Check which process is using a port:

```powershell
netstat -ano | findstr :8080
netstat -ano | findstr :3000
netstat -ano | findstr :5432
```

Stop or reconfigure the conflicting service before restarting GeoTrack.

### React routes return 404 after refreshing

Confirm that the frontend Nginx configuration includes:

```nginx
location / {
    try_files $uri $uri/ /index.html;
}
```

## 12. Security Considerations

- Keep actual environment files out of version control.
- Use strong database passwords.
- Restrict database access in production.
- Configure production CORS origins.
- Use HTTPS and WSS for production traffic.
- Configure authentication and authorization before exposing sensitive operations publicly.
- Never publish database credentials or API secrets.

## 13. Project Structure

```text
GeoTrack/
│
├── backend/
│   ├── cmd/
│   ├── internal/
│   ├── migrations/
│   ├── tests/
│   ├── Dockerfile
│   ├── .dockerignore
│   ├── .env
│   └── go.mod
│
├── frontend/
│   ├── src/
│   ├── Dockerfile
│   ├── nginx.conf
│   └── package.json
│
├── docker-compose.yml
├── .env.example
├── .gitignore
├── README.md
└── SETUP.md
```

## 14. Stopping the Application

To stop all containers:

```powershell
docker compose down
```

To restart:

```powershell
docker compose up -d
```

To rebuild after code changes:

```powershell
docker compose up --build -d
```

The PostgreSQL named volume is retained unless explicitly removed.

---

**GeoTrack – Geofencing & Real-time Alert System**

Full-Stack Developer Assessment Project
