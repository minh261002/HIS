<div align="center">

[![ServBay](public/logo.png)](https://www.servbay.com)

# Hospital Information System (HIS)

**Modern Hospital Management Backend Built with Go**

[![Go Version](https://img.shields.io/badge/Go-1.24.4-00ADD8?style=flat&logo=go)](https://go.dev/)
[![MySQL](https://img.shields.io/badge/MySQL-8.0-4479A1?style=flat&logo=mysql&logoColor=white)](https://www.mysql.com/)
[![GORM](https://img.shields.io/badge/GORM-ORM-00ADD8?style=flat)](https://gorm.io/)
[![Gin](https://img.shields.io/badge/Gin-Web_Framework-00ADD8?style=flat)](https://gin-gonic.com/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

**Developed with [ServBay](https://www.servbay.com) - The Ultimate Local Development Environment**

[Features](#-features) • [Architecture](#-architecture) • [Getting Started](#-getting-started) • [API Documentation](#-api-documentation) • [Modules](#-modules)

</div>

---

## 📋 Overview

A comprehensive Hospital Information System (HIS) backend built with **Go**, featuring clean architecture, JWT authentication, MySQL database with GORM ORM, and RESTful APIs. Designed for modern healthcare facilities with modular, scalable architecture that can be deployed in phases.

### Why This Project?

- **Production-Ready**: Built with enterprise-grade patterns and best practices
- **Modular Design**: Easy to extend and customize for specific hospital needs
- **Clean Architecture**: Clear separation of concerns across layers
- **Comprehensive**: Covers patient management, clinical workflows, billing, and more
- **Developer-Friendly**: Hot reload, Docker support, comprehensive documentation

---

## ✨ Features

### 🔐 Authentication & Security

- **JWT-based authentication** with access and refresh tokens
- **Role-Based Access Control (RBAC)** with granular permissions
- **Bcrypt password hashing** for secure credential storage
- **Configurable token expiry** via environment variables
- **CORS middleware** for cross-origin security

### 🏥 Clinical Modules

- **Patient Management**: Demographics, insurance, emergency contacts, medical history, allergies
- **Appointments & Visits**: Complete scheduling and visit lifecycle management
- **Diagnoses**: ICD-10 code integration and diagnosis tracking
- **Prescriptions**: Medication management with dispensing workflow
- **Lab Tests**: Test templates, requests, sample tracking, and results
- **Imaging**: Imaging templates, requests, scheduling, and reporting
- **Inpatient Care**: Admissions, bed management, nursing notes

### 💊 Pharmacy & Inventory

- **Medication Catalog**: Comprehensive drug database
- **Inventory Management**: Stock levels, batch tracking, expiry monitoring
- **Dispensing**: Prescription fulfillment and tracking
- **Alerts**: Low stock and expiring medication notifications

### 💰 Billing & Finance

- **Invoice Management**: Service-based billing with detailed line items
- **Payment Processing**: Multiple payment methods and tracking
- **Insurance Claims**: Internal claim management and processing
- **Financial Reports**: Revenue tracking and analytics

### 🏢 System Administration

- **Department Management**: Organizational structure and hierarchy
- **Medical Services Catalog**: Service types, pricing, and departments
- **User & Role Management**: Staff accounts and permission assignment
- **Audit Logging**: Complete activity tracking with IP and user agent

### 🛠️ Developer Experience

- **Hot Reload**: Live development with Air
- **Docker Compose**: MySQL, Redis, and Adminer pre-configured
- **Database Migrations**: Version-controlled schema with golang-migrate
- **Structured Logging**: JSON logging with Zap
- **Makefile**: Common tasks automated
- **API Documentation**: OpenAPI/Swagger specification included

---

## 🏗 Architecture

### Clean Architecture Layers

```
┌─────────────────────────────────────────────────────────┐
│                    HTTP Layer (Gin)                     │
│                  Handlers + Middleware                  │
└────────────────────┬────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────┐
│                   Service Layer                         │
│              Business Logic + Use Cases                 │
└────────────────────┬────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────┐
│                 Repository Layer                        │
│              Data Access + GORM Queries                 │
└────────────────────┬────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────┐
│                   Domain Layer                          │
│            Models + Business Entities                   │
└─────────────────────────────────────────────────────────┘
```

### Project Structure

```
his/
├── cmd/api/              # Application entry point
│   └── main.go           # Dependency injection, routing, server setup
├── internal/
│   ├── config/           # Configuration loading and database initialization
│   ├── domain/           # Domain models (31 GORM models)
│   ├── repository/       # Data access layer (GORM queries, search, stats)
│   ├── service/          # Business logic and use cases
│   ├── handler/          # HTTP handlers (Gin controllers)
│   ├── dto/              # Data Transfer Objects (requests/responses)
│   ├── middleware/       # Auth, RBAC, CORS, logging, recovery
│   └── pkg/              # Shared utilities (JWT, logger, response helpers)
├── migrations/           # Database migrations (78 migration files)
├── docker/               # Docker Compose configuration
├── docs/                 # API documentation (OpenAPI spec)
├── public/               # Static assets (logos, images)
├── .env.example          # Environment variable template
├── Makefile              # Development and build automation
└── README.md
```

### Request Flow

```
1. Client Request
   ↓
2. Middleware Chain (CORS → Logger → Recovery → Auth → RBAC)
   ↓
3. Handler (DTO binding, validation, context extraction)
   ↓
4. Service (Business logic execution)
   ↓
5. Repository (Database operations via GORM)
   ↓
6. Domain Models (GORM entities)
   ↓
7. Response (JSON via helper functions)
```

---

## 🚀 Getting Started

### Prerequisites

- **Go**: 1.21 or higher ([Download](https://go.dev/dl/))
- **Docker & Docker Compose**: For local services ([Download](https://www.docker.com/))
- **ServBay**: Recommended for local development ([Download](https://www.servbay.com))
- **golang-migrate**: For database migrations

  ```bash
  # macOS
  brew install golang-migrate

  # Or via Go
  go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
  ```

- **Air** (optional): For hot reload
  ```bash
  go install github.com/air-verse/air@latest
  ```

### Installation

#### 1. Clone the Repository

```bash
git clone <your-repo-url> his
cd his
```

#### 2. Install Dependencies

```bash
go mod download
```

#### 3. Configure Environment

```bash
cp .env.example .env
```

Edit `.env` with your configuration:

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=13306
DB_USER=his_user
DB_PASSWORD=his_password
DB_NAME=hospital_db

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=16379
REDIS_PASSWORD=
REDIS_DB=0

# JWT Configuration
JWT_SECRET=your-super-secret-key-change-this
JWT_ACCESS_TOKEN_EXPIRY=15m
JWT_REFRESH_TOKEN_EXPIRY=168h

# Server Configuration
SERVER_PORT=8080
SERVER_MODE=debug
SERVER_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173

# Logging
LOG_LEVEL=debug
LOG_FORMAT=json
```

#### 4. Start Docker Services

```bash
make docker-up
```

This starts:

- **MySQL** on port `13306`
- **Redis** on port `16379`
- **Adminer** on port `8081` (database management UI)

#### 5. Run Database Migrations

```bash
make migrate-up
```

This creates all tables and seeds initial data (ICD-10 codes, medications, lab templates, etc.)

#### 6. Run the Application

**Development mode (with hot reload):**

```bash
make run
```

**Production build:**

```bash
make build
./bin/api
```

The API will be available at: **http://localhost:8080**

---

## 🗄️ Database Management

### Migrations

**Create a new migration:**

```bash
make migrate-create name=create_new_table
```

**Apply migrations:**

```bash
make migrate-up
```

**Rollback migrations:**

```bash
make migrate-down
```

### Adminer (Database UI)

Access the database management interface:

- **URL**: http://localhost:8081
- **System**: MySQL
- **Server**: mysql
- **Username**: his_user
- **Password**: his_password
- **Database**: hospital_db

---

## 📚 Modules

### 1. Authentication & User Management

**Domain Models**: `User`, `Role`, `Permission`

**Endpoints**:

- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/refresh` - Refresh access token
- `GET /api/v1/auth/profile` - Get current user profile
- `/api/v1/users` - User CRUD operations
- `/api/v1/roles` - Role management
- `/api/v1/permissions` - Permission management

**Features**:

- JWT token management
- Many-to-many role-permission relationships
- RBAC enforcement via middleware

---

### 2. Patient Management

**Domain Models**: `Patient`, `PatientAllergy`, `PatientMedicalHistory`

**Key Endpoints**:

- `POST /api/v1/patients` - Create patient
- `GET /api/v1/patients` - List patients (paginated)
- `GET /api/v1/patients/search` - Search patients
- `GET /api/v1/patients/stats` - Patient statistics
- `GET /api/v1/patients/code/:code` - Get by patient code
- `GET /api/v1/patients/:id` - Get patient details
- `GET /api/v1/patients/:id/allergies` - Patient allergies
- `GET /api/v1/patients/:id/medical-history` - Medical history
- `GET /api/v1/patients/:id/visits` - Patient visits
- `GET /api/v1/patients/:id/appointments` - Patient appointments

**Features**:

- Demographics and insurance information
- Emergency contact management
- Allergy tracking
- Medical history records
- Comprehensive patient search

---

### 3. Appointments & Visits

**Domain Models**: `Appointment`, `Visit`

**Appointment Lifecycle**:
`SCHEDULED → CONFIRMED → IN_PROGRESS → COMPLETED / CANCELLED / NO_SHOW`

**Visit Lifecycle**:
`ACTIVE → COMPLETED / CANCELLED`

**Key Endpoints**:

- `/api/v1/appointments` - Appointment CRUD
- `/api/v1/visits` - Visit management
- `/api/v1/doctors/:id/schedule` - Doctor's schedule
- `/api/v1/doctors/:id/visits` - Doctor's visits
- `/api/v1/doctors/:id/available-slots` - Available time slots

---

### 4. Diagnoses & ICD-10

**Domain Models**: `Diagnosis`, `ICD10Code`

**Features**:

- ICD-10 code database (seeded via migrations)
- Search ICD-10 codes by code or description
- Link diagnoses to visits and patients
- Primary and secondary diagnosis support

---

### 5. Medications & Prescriptions

**Domain Models**: `Medication`, `Prescription`, `PrescriptionItem`

**Prescription Status**: `ACTIVE → DISPENSED / COMPLETED / CANCELLED`

**Key Endpoints**:

- `/api/v1/medications` - Medication catalog
- `/api/v1/prescriptions` - Prescription management
- `/api/v1/patients/:id/prescriptions` - Patient prescriptions

---

### 6. Inventory & Dispensing

**Domain Models**: `Inventory`, `Dispensing`

**Features**:

- Stock level tracking
- Batch management with expiry dates
- Low stock alerts
- Expiring medication alerts
- Dispensing records linked to prescriptions

**Key Endpoints**:

- `/api/v1/inventory` - Inventory management
- `/api/v1/inventory/low-stock` - Low stock alerts
- `/api/v1/inventory/expiring-soon` - Expiring items
- `/api/v1/dispensing` - Dispensing records

---

### 7. Lab Tests

**Domain Models**: `LabTestTemplate`, `LabTestTemplateParameter`, `LabTestRequest`, `LabTestResult`

**Workflow**: `PENDING → SAMPLE_COLLECTED → IN_PROGRESS → COMPLETED / CANCELLED`

**Features**:

- Test templates with parameters
- Sample collection tracking
- Result entry and reporting
- Patient and doctor views

---

### 8. Imaging

**Domain Models**: `ImagingTemplate`, `ImagingRequest`, `ImagingResult`

**Workflow**: `PENDING → SCHEDULED → IN_PROGRESS → COMPLETED / CANCELLED`

**Features**:

- Imaging modality templates (X-Ray, CT, MRI, Ultrasound)
- Request scheduling
- Report generation
- DICOM integration ready

---

### 9. Inpatient Care

**Domain Models**: `Admission`, `Bed`, `BedAllocation`, `NursingNote`

**Bed Types**: `STANDARD`, `ICU`, `EMERGENCY`, `ISOLATION`, `MATERNITY`, `PEDIATRIC`

**Key Endpoints**:

- `/api/v1/beds/available` - Available beds
- `/api/v1/admissions` - Admission management
- `/api/v1/admissions/:id/discharge` - Discharge patient
- `/api/v1/admissions/:id/transfer` - Transfer bed
- `/api/v1/admissions/:id/nursing-notes` - Nursing documentation

---

### 10. Billing & Finance

**Domain Models**: `Invoice`, `InvoiceItem`, `Payment`, `InsuranceClaim`

**Features**:

- Service-based invoicing
- Multiple payment methods
- Insurance claim management
- Payment tracking and reconciliation

**Key Endpoints**:

- `/api/v1/invoices` - Invoice management
- `/api/v1/invoices/code/:code` - Get by invoice code
- `/api/v1/payments` - Payment processing
- `/api/v1/insurance-claims` - Claim management

> **Note**: Current insurance support is internal claim management only. No direct integration with national health insurance gateways yet.

---

### 11. System Administration

**Domain Models**: `Department`, `MedicalService`, `AuditLog`

**Features**:

- Department hierarchy management
- Medical service catalog with pricing
- Comprehensive audit logging

**Audit Actions**: `CREATE`, `UPDATE`, `DELETE`, `VIEW`, `LOGIN`, `LOGOUT`

**Key Endpoints**:

- `/api/v1/system/departments` - Department management
- `/api/v1/system/medical-services` - Service catalog
- `/api/v1/system/audit-logs` - Audit log queries (filterable by user, resource, date range)

---

## 📖 API Documentation

Comprehensive API documentation is available in the `docs/` directory:

- **[API.md](docs/API.md)** - Detailed endpoint documentation
- **[openapi.yaml](docs/openapi.yaml)** - OpenAPI 3.0 specification

You can import the OpenAPI spec into tools like:

- [Swagger Editor](https://editor.swagger.io/)
- [Postman](https://www.postman.com/)
- [Insomnia](https://insomnia.rest/)

---

## 🔒 Security

### Authentication

- **JWT tokens** with configurable expiry
- **Access tokens** (short-lived, default 15 minutes)
- **Refresh tokens** (long-lived, default 7 days)
- Token secrets configured via environment variables

### Password Security

- **Bcrypt hashing** with appropriate cost factor
- No plain-text password storage
- Secure password reset flows

### Authorization

- **RBAC middleware** enforces permissions on protected routes
- **Granular permissions** (e.g., `patients.view`, `appointments.create`, `invoices.update`)
- **Role-based access** with many-to-many role-permission mapping

### API Security

- **CORS middleware** with configurable allowed origins
- **Request validation** via Gin binding tags
- **Parameterized queries** via GORM (SQL injection prevention)
- **Rate limiting** ready for implementation

---

## 🧪 Testing

Run the test suite:

```bash
make test
```

Run tests with coverage:

```bash
go test -v -cover -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

---

## 🔧 Development Commands

The `Makefile` provides convenient commands for common tasks:

```bash
make help          # Show all available commands
make run           # Run with hot reload (Air)
make build         # Build production binary
make test          # Run tests
make docker-up     # Start Docker services
make docker-down   # Stop Docker services
make migrate-up    # Apply database migrations
make migrate-down  # Rollback migrations
make migrate-create name=migration_name  # Create new migration
make clean         # Remove build artifacts
make fmt           # Format code (gofmt)
make tidy          # Tidy dependencies
make lint          # Run linter (requires golangci-lint)
make deps          # Download dependencies
```

---

## 🗺️ Roadmap

### ✅ Phase 1 (Completed)

- Authentication and authorization
- User, role, and permission management
- Patient management
- Appointments and visits
- ICD-10 integration
- Medications and prescriptions

### ✅ Phase 2 (Completed)

- Lab tests and imaging
- Inventory and dispensing
- Low stock and expiry alerts

### ✅ Phase 3 (Completed)

- Inpatient admissions
- Bed management and allocation
- Nursing notes
- Billing and invoicing
- Payment processing
- Insurance claim management

### 🚧 Phase 4 (In Progress)

- Surgery module
- Operating room scheduling
- Surgical procedure tracking

### 📋 Phase 5 (Planned)

- Reporting and dashboards
- Analytics and BI integration
- Data export capabilities
- National health insurance gateway integration
- FHIR API support
- Mobile app backend APIs

---

## 🤝 Contributing

We welcome contributions! Here's how you can help:

1. **Fork** the repository
2. **Create** a feature branch (`git checkout -b feature/amazing-feature`)
3. **Commit** your changes (`git commit -m 'Add amazing feature'`)
4. **Push** to the branch (`git push origin feature/amazing-feature`)
5. **Open** a Pull Request

### Contribution Guidelines

- Write clear, descriptive commit messages
- Follow Go best practices and conventions
- Add tests for new features
- Update documentation as needed
- Ensure all tests pass before submitting PR
- Keep PRs focused on a single feature or fix

---

## 📄 License

This project is licensed under the **MIT License** - see the [LICENSE](LICENSE) file for details.

---

## 💬 Support

- **Issues**: [GitHub Issues](https://github.com/minhtran/his/issues)
- **Discussions**: [GitHub Discussions](https://github.com/minhtran/his/discussions)
- **Email**: support@example.com

---

## 🙏 Acknowledgments

- **[ServBay](https://www.servbay.com)** - The ultimate local development environment that powers this project
- **[Gin](https://gin-gonic.com/)** - High-performance HTTP web framework
- **[GORM](https://gorm.io/)** - Fantastic ORM library for Go
- **[golang-migrate](https://github.com/golang-migrate/migrate)** - Database migration tool
- **[Zap](https://github.com/uber-go/zap)** - Blazing fast, structured logging

---

<div align="center">

**Built with ❤️ using [ServBay](https://www.servbay.com)**

[![ServBay](https://img.shields.io/badge/Powered_by-ServBay-00ADD8?style=for-the-badge)](https://www.servbay.com)

</div>
