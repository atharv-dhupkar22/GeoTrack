
# GeoTrack – Geofencing & Real-time Alert System

GeoTrack is a full-stack vehicle tracking and geofencing platform designed to monitor vehicle locations, detect geofence entry and exit events, and deliver real-time alerts.

The application combines a Go REST API, PostgreSQL with PostGIS, a React frontend, and WebSocket-based notifications.

## Features

### Geofence Management
- Create and manage geographic zones.
- Support delivery, restricted, toll, and customer-area categories.
- Validate polygon coordinates.
- Detect whether a vehicle is inside or outside a geofence.

### Vehicle Management
- Register vehicles and driver information.
- Retrieve registered vehicle records.
- Track vehicle locations using GPS coordinates.

### Live GPS Tracking
- Submit vehicle latitude and longitude.
- View recorded vehicle locations on an interactive map.
- Detect geofence entry and exit transitions.

### Real-time Alerts
- Configure entry, exit, or both event notifications.
- Receive WebSocket notifications.
- Store alert events for later retrieval.

### Violation History
- Retrieve recorded geofence violations.
- Filter by vehicle, geofence, and event type.
- Support pagination.

### Dashboard
- View vehicle and geofence statistics.
- Monitor recent alerts and violations.
- Access the platform's main features.

## Technology Stack

| Layer | Technology |
|---|---|
| Frontend | React, TypeScript, Vite |
| UI Styling | Tailwind CSS |
| Maps | Leaflet, React Leaflet |
| Backend | Go |
| API | REST |
| Real-time | WebSocket |
| Database | PostgreSQL |
| Spatial Database | PostGIS |
| Database Driver | pgx |
| Containerization | Docker |
| Orchestration | Docker Compose |
| Web Server | Nginx |

## Architecture

```text
                 GeoTrack Platform
                         |
                 React Frontend
                    Port 3000
                         |
                  REST API / WS
                         |
                    Go Backend
                    Port 8080
                         |
                  PostgreSQL
                   + PostGIS
                    Port 5432
```

## API Endpoints

| Method | Endpoint | Purpose |
|---|---|---|
| GET | /health | Backend health |
| POST | /geofences | Create geofence |
| GET | /geofences | List geofences |
| POST | /vehicles | Register vehicle |
| GET | /vehicles | List vehicles |
| POST | /vehicles/location | Submit GPS location |
| GET | /vehicles/location/{vehicle_id} | Get latest location |
| POST | /alerts/configure | Configure alerts |
| GET | /alerts/configurations | List alert configurations |
| GET | /alerts | Retrieve alert history |
| GET | /violations/history | Retrieve violation history |
| WS | /ws/alerts | Real-time alert connection |

API responses include a `time_ns` field.

## Project Structure

```text
GeoTrack/
│
├── backend/
│   ├── cmd/
│   ├── internal/
│   │   ├── config/
│   │   ├── database/
│   │   ├── geofence/
│   │   ├── handlers/
│   │   ├── middleware/
│   │   ├── models/
│   │   ├── repository/
│   │   ├── services/
│   │   └── websocket/
│   ├── migrations/
│   ├── tests/
│   ├── Dockerfile
│   ├── .dockerignore
│   ├── .env
│   └── go.mod
│
├── frontend/
│   ├── src/
│   │   ├── components/
│   │   ├── context/
│   │   ├── hooks/
│   │   ├── pages/
│   │   ├── services/
│   │   ├── types/
│   │   └── utils/
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

## Quick Start with Docker

### Requirements
- Docker Desktop
- Docker Compose

### 1. Clone the repository

```powershell
git clone YOUR_GITHUB_REPOSITORY_URL
cd GeoTrack
```

### 2. Configure environment variables

Copy the example environment file:

```powershell
Copy-Item .env.example .env
```

Edit `.env` and set a secure database password.

### 3. Start the application

```powershell
docker compose up --build -d
```

### 4. Access the application

- Frontend: http://localhost:3000
- Backend: http://localhost:8080
- Health: http://localhost:8080/health
- WebSocket: ws://localhost:8080/ws/alerts

### 5. Check containers

```powershell
docker compose ps
```

### 6. Stop the application

```powershell
docker compose down
```

The database volume is retained when using this command.

## Local Development

### Backend

```powershell
cd backend
go mod download
go test ./...
go run ./cmd
```

### Frontend

Open another terminal:

```powershell
cd frontend
npm install
npm run dev
```

Configure the frontend environment variables in `frontend/.env` as described in `SETUP.md`.

## Testing

Backend tests:

```powershell
cd backend
go test ./... -v
```

Static analysis:

```powershell
go vet ./...
```

Backend build:

```powershell
go build ./...
```

Frontend production build:

```powershell
cd frontend
npm run build
```

## Security Notes

- Do not commit `.env` files or production credentials.
- Use strong database passwords in deployed environments.
- Configure CORS for the actual frontend domain before deployment.
- Use HTTPS and secure WebSocket connections in production.
- Configure authentication and authorization before exposing sensitive endpoints publicly.

## Project Status

The local development environment and Docker Compose deployment have been started and verified.

Further production deployment, automated testing expansion, and submission preparation remain separate activities.

## License

This project was developed as a full-stack engineering assessment project.
