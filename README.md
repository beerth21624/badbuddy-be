# BadBuddy Backend (public version)

A Go-based backend service for BadBuddy - a platform for managing sports venues, court bookings, and player sessions.

## Technology Stack

- Go 1.22.5
- PostgreSQL 13
- Docker & Docker Compose
- Fiber web framework
- SQLx for database operations
- WebSocket support for real-time chat
- JWT authentication

## Key Features
- User management and authentication
- Venue and court management
- Session booking and management
- Real-time chat functionality
- Player reviews and ratings
- RESTful API endpoints

## Project Structure

```
.
├── cmd/
│   └── api/           # Application entry point
├── config/            # Configuration management
├── internal/
│   ├── delivery/      # HTTP handlers and DTOs
│   ├── domain/        # Domain models
│   ├── infrastructure/# Database and server setup
│   ├── repositories/  # Data access layer
│   └── usecase/       # Business logic
└── db/
    └── migrations/    # Database migrations
```

## Getting Started

### Prerequisites

- Go 1.22 or higher
- PostgreSQL 13
- Docker (optional)

### Local Development

1. Clone the repository:
```bash
git clone https://github.com/SmoothBrain-Tech-X/badbuddy-be.git
```

2. Install dependencies:
```bash
go mod download
```

3. Set up environment variables in .env:
```env
# Database configuration
DB_HOST=         # Hostname or IP address of the database server (e.g., localhost)
DB_PORT=         # Port number of the database (e.g., 5432 for PostgreSQL)
DB_USER=         # Username for the database connection
DB_PASSWORD=     # Password for the database connection
DB_NAME=         # Name of the database to connect to
DB_SSLMODE=      # SSL mode for the database connection (e.g., 'require', 'disable', 'verify-full')

# JWT configuration
JWT_SECRET=      # Secret key for signing JWT tokens
JWT_EXPIRATION=  # Expiration time for JWT tokens

# Server configuration
PORT=            # Port number for running the application (e.g., 3000)
```

4. Run the application:
```bash
go run cmd/api/main.go
```

### Docker Deployment

Build and run using Docker Compose:

```bash
docker-compose up --build
```

## API Documentation

The API provides the following main endpoints:

- `/api/users` - User management
- `/api/venues` - Venue management
- `/api/bookings` - Booking operations
- `/api/sessions` - Session menagement
- `/api/chats` - Chat functionality
- `/ws/:chat_id` - WebSocket endpoint for real-time chat

## Testing and Development

For local development with hot reload:

```bash
air
```

Run tests:

```bash
go test ./...
```

## License

This project is licensed under the MIT License.


# badbuddy-be
