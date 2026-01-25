# Hospital Information System (HIS) - Golang

A modern Hospital Information System built with Golang, featuring clean architecture, JWT authentication, and RESTful APIs.

## 🚀 Features

- **Clean Architecture**: Separation of concerns with domain, repository, service, and handler layers
- **JWT Authentication**: Secure token-based authentication with access and refresh tokens
- **Database**: MySQL with GORM ORM and migration support
- **Logging**: Structured logging with Zap
- **Middleware**: CORS, authentication, request logging, and panic recovery
- **Docker**: Containerized MySQL and Redis services
- **Hot Reload**: Development with Air for instant code reloading

## 📋 Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose
- Make (optional, for convenience commands)
- golang-migrate (for database migrations)

## 🛠️ Installation

1. **Clone the repository**

```bash
cd /Users/minhtran/Glang/his
```

2. **Install dependencies**

```bash
go mod download
```

3. **Install golang-migrate**

```bash
# macOS
brew install golang-migrate

# Or using Go
go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

4. **Install Air for hot reload (optional)**

```bash
go install github.com/air-verse/air@latest
```

5. **Setup environment**

```bash
cp .env.example .env
# Edit .env with your configuration
```

6. **Start Docker services**

```bash
make docker-up
# Or: cd docker && docker-compose up -d
```

7. **Run database migrations**

```bash
make migrate-up
```

## 🏃 Running the Application

### Development (with hot reload)

```bash
make run
# Or: air
```

### Production build

```bash
make build
./bin/api
```

The API will be available at `http://localhost:8080`

## 📚 API Endpoints

### Health Check

- `GET /health` - Server health status

### Authentication

- `POST /api/v1/auth/register` - Register new user
- `POST /api/v1/auth/login` - Login user
- `POST /api/v1/auth/refresh` - Refresh access token
- `GET /api/v1/auth/profile` - Get current user profile (protected)

### Example Requests

**Register**

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "email": "john@example.com",
    "password": "password123",
    "full_name": "John Doe",
    "phone_number": "0123456789"
  }'
```

**Login**

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "password": "password123"
  }'
```

**Get Profile**

```bash
curl -X GET http://localhost:8080/api/v1/auth/profile \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

## 🗄️ Database

### Migrations

**Create new migration**

```bash
make migrate-create name=create_patients_table
```

**Run migrations**

```bash
make migrate-up
```

**Rollback migrations**

```bash
make migrate-down
```

### Access Database

**Using Adminer** (Web UI)

- URL: http://localhost:8081
- System: MySQL
- Server: mysql
- Username: his_user
- Password: his_password
- Database: hospital_db

## 🏗️ Project Structure

```
his/
├── cmd/api/              # Application entry point
├── internal/
│   ├── config/          # Configuration management
│   ├── domain/          # Domain models
│   ├── repository/      # Data access layer
│   ├── service/         # Business logic
│   ├── handler/         # HTTP handlers
│   ├── middleware/      # HTTP middleware
│   ├── dto/             # Data transfer objects
│   └── pkg/             # Shared utilities
├── migrations/          # Database migrations
├── docker/              # Docker configuration
├── .env.example         # Environment template
├── Makefile            # Development commands
└── README.md
```

## 🧪 Testing

```bash
make test
```

## 🔧 Available Make Commands

```bash
make help          # Show all available commands
make run           # Run with hot reload
make build         # Build binary
make test          # Run tests
make docker-up     # Start Docker services
make docker-down   # Stop Docker services
make migrate-up    # Run migrations
make migrate-down  # Rollback migrations
make clean         # Clean build artifacts
make fmt           # Format code
make tidy          # Tidy dependencies
```

## 🛡️ Security

- Passwords are hashed using bcrypt
- JWT tokens with configurable expiry
- CORS protection
- Request validation
- SQL injection protection via GORM

## 📝 Environment Variables

See `.env.example` for all available configuration options.

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## 📄 License

This project is licensed under the MIT License.
