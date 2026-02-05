## Hospital Information System (HIS) – Golang

Hospital Information System (HIS) backend implemented in Go, using clean architecture, JWT authentication, MySQL + GORM and RESTful APIs.  
The goal is to provide a modern, modular HIS backend that can be rolled out in phases in real hospital environments.

---

## 🎯 Key Features

- **Clean, extensible architecture**
  - Clear layers: `domain` (business models), `repository` (data access), `service` (business logic), `handler` (HTTP), `dto` (request/response), `middleware` (cross‑cutting).
  - Manual dependency wiring in `cmd/api/main.go` for full control.

- **Authentication & Authorization**
  - **JWT**: access + refresh tokens, configurable via environment variables.
  - **RBAC** based on permission codes (e.g. `patients.view`, `appointments.create`, `invoices.view`, `departments.update`, `audit.view`).
  - `AuthMiddleware` + `RBACMiddleware` on protected routes.

- **Database & migrations**
  - **MySQL** + **GORM** ORM.
  - SQL migrations in `migrations/` managed via `golang-migrate`.
  - Seed data for ICD‑10, medications, lab templates, imaging templates, beds, etc.

- **Infrastructure & Dev Experience**
  - **Docker Compose**: MySQL, Redis, Adminer.
  - **Air** for hot reload in development.
  - Structured logging with **Zap** and a small logging wrapper.

- **Business modules already implemented**
  - **User / Role / Permission** management, auth, profile.
  - **Patients**, medical history and allergies.
  - **Appointments** and **visits**.
  - **Diagnoses** and **ICD‑10** catalog.
  - **Medications, prescriptions, inventory, dispensing.**
  - **Clinical support**: lab tests and imaging.
  - **Inpatient module**: admissions, beds & bed allocation, nursing notes.
  - **Billing & accounting**: invoices, payments, insurance claims.
  - **System module**: departments, medical services, audit logging.

---

## 🏗 Architecture & Folder Structure

```text
his/
├── cmd/api/              # Application entry point
│   └── main.go           # Init config, DB, logger, DI, routes
├── internal/
│   ├── config/           # Config loading, DB initialization
│   ├── domain/           # Domain models (GORM models, enums, hooks, ...)
│   ├── repository/       # Data access (GORM queries, search, stats)
│   ├── service/          # Business logic / use cases
│   ├── handler/          # Gin HTTP handlers, DTO binding, responses
│   ├── middleware/       # Auth, RBAC, CORS, logger, recovery
│   ├── dto/              # DTOs for requests and responses
│   └── pkg/              # Shared utilities (jwt, logger, response)
├── migrations/           # Database migrations (golang-migrate)
├── docker/               # Docker Compose stack
├── .env.example          # Environment variable template
├── Makefile              # Dev/build/migration helpers
└── README.md
```

### Request Lifecycle (high‑level)

1. Client sends HTTP request to Gin router.
2. Request goes through middleware chain: **CORS → Logger → Recovery → Auth (JWT) → RBAC (if applied)**.
3. Handler binds JSON / query / path params into **DTOs**, validates them, and reads `userID` from context.
4. Handler calls corresponding **Service**.
5. Service executes business logic, talks to **Repository** using **domain models** (GORM).
6. Service maps results to **DTO responses** (where applicable) and returns to handler.
7. Handler returns JSON using helper functions like `response.Success`, `response.SuccessPaginated`, `response.BadRequest`, etc.

---

## 📦 Main Modules

### 1. Authentication & User Management

- **Auth**
  - Endpoints:
    - `POST /api/v1/auth/login`
    - `POST /api/v1/auth/refresh`
    - `GET /api/v1/auth/profile`
  - Tokens are managed via `internal/pkg/jwt`.

- **User / Role / Permission**
  - Domain: `User`, `Role`, `Permission` with many‑to‑many relationships.
  - Endpoints:
    - `/api/v1/users` – CRUD + role assignment.
  - Used by RBAC middleware to enforce permissions.

### 2. Patient & Clinical Workflow

- **Patients**
  - Manage demographics, insurance info, emergency contacts, history and allergies.
  - Endpoints (examples):
    - `POST /api/v1/patients`
    - `GET /api/v1/patients`
    - `GET /api/v1/patients/search`
    - `GET /api/v1/patients/stats`
    - `GET /api/v1/patients/code/:code`
    - `GET /api/v1/patients/:id`
    - Sub‑routes:  
      `/:id/allergies`, `/:id/medical-history`, `/:id/visits`, `/:id/appointments`,  
      `/:id/prescriptions`, `/:id/admissions`, `/:id/lab-tests`,  
      `/:id/imaging-requests`, `/:id/invoices`, `/:id/dispensing`.

- **Appointments & Visits**
  - Appointment lifecycle: schedule → confirm → start → complete → cancel → no‑show.
  - Visit lifecycle: create → update → complete → cancel.
  - Doctor‑specific views:
    - `GET /api/v1/doctors/:id/schedule`
    - `GET /api/v1/doctors/:id/visits`
    - `GET /api/v1/doctors/:id/available-slots`

- **Diagnoses & ICD‑10**
  - Search ICD‑10 codes and assign diagnoses to visits/patients.

### 3. Medications, Prescriptions, Inventory & Dispensing

- **Medications & Prescriptions**
  - Search medications, create/update prescriptions, track dispensing status (dispensed / complete / canceled).

- **Inventory & Dispensing**
  - Manage stock levels and batches.
  - Alerts for **low stock** and **expiring soon**.
  - Dispensing records linked to prescriptions and patients.

### 4. Paraclinical: Lab Tests & Imaging

- **Lab Tests**
  - Test templates, parameters, test requests, and results.
  - Workflow: create request → collect sample → start processing → complete → enter results → cancel.

- **Imaging**
  - Imaging templates, imaging requests, and reports.
  - Workflow: create request → schedule → start → complete → report / cancel.

### 5. Inpatient Module: Beds & Nursing Notes

- **Beds & Bed Allocation**
  - Domain `Bed` with enums `BedType`, `BedStatus`, `DepartmentType`.
  - Endpoints:
    - `GET /api/v1/beds/available`
    - `GET /api/v1/beds/:number`

- **Admissions & Nursing Notes**
  - Admissions are linked to patients, visits and beds.
  - Endpoints:
    - `/api/v1/admissions` – create, get by ID/code, discharge, transfer bed, list active admissions.
    - `POST /api/v1/admissions/:id/nursing-notes`
    - `GET /api/v1/admissions/:id/nursing-notes`

### 6. Billing, Accounting & Insurance

- **Invoices, Payments, Insurance Claims**
  - Create invoices from services / prescriptions / procedures.
  - Record payments.
  - Create and manage insurance claims (approve / reject) linked to invoices.
  - Example endpoints:
    - `/api/v1/invoices`, `/api/v1/invoices/code/:code`, `/api/v1/patients/:id/invoices`
    - `/api/v1/payments`, `/api/v1/invoices/:id/payments`
    - `/api/v1/insurance-claims`, `/api/v1/invoices/:id/insurance-claims`

> Note: current insurance support is **internal claim management only** – there is no direct integration yet with national health insurance gateways.

### 7. System Module (Departments, Services, Audit Logs)

- **Departments**
  - Manage hospital departments, head of department and active status.
  - API responses use `DepartmentResponse` DTO, optionally including a simplified head‑doctor user object.

- **Medical Services**
  - Catalog of medical services: type (consultation, lab test, imaging, procedure, other), base price, department.
  - Responses use `MedicalServiceResponse` DTO, embedding `DepartmentResponse`.

- **Audit Logs**
  - Track system actions: CREATE / UPDATE / DELETE / VIEW / LOGIN / LOGOUT per resource and user, including IP and user agent.
  - `GET /api/v1/system/audit-logs` supports:
    - Filtering by `user_id`, `resource`, `from_date`, `to_date`.
    - Pagination metadata.

---

## ⚙️ System Requirements

- Go **1.21+**
- Docker & Docker Compose
- Make (optional but recommended)
- `golang-migrate` CLI for DB migrations:
  - `brew install golang-migrate` (macOS), or
  - `go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest`

---

## 🚀 Getting Started

### 1. Clone & install dependencies

```bash
git clone <your-repo-url> his
cd his

go mod download
```

### 2. Configure environment

```bash
cp .env.example .env
# Edit .env to match your setup (DB, JWT_SECRET, LOG_LEVEL, etc.)
```

### 3. Start Docker services

```bash
make docker-up
# or:
cd docker && docker-compose up -d
```

### 4. Run database migrations

```bash
make migrate-up
```

### 5. Run the application

#### Development (with hot reload)

```bash
make run
# or, if you installed Air:
air
```

#### Production build

```bash
make build
./bin/api
```

By default, the API listens on: `http://localhost:8080`

---

## 🗄️ Database & Adminer

### Migrations

- Create a new migration:

```bash
make migrate-create name=create_patients_table
```

- Apply migrations:

```bash
make migrate-up
```

- Rollback:

```bash
make migrate-down
```

### Access via Adminer

- URL: `http://localhost:8081`
- System: `MySQL`
- Server: `mysql`
- Username: `his_user`
- Password: `his_password`
- Database: `hospital_db`

---

## 🔒 Security

- User passwords are hashed with **bcrypt**.
- JWT expiry and secrets are configurable via environment variables.
- CORS is enforced via middleware and configuration.
- Request validation is handled by Gin binding with `binding` tags.
- GORM and parameterized queries help prevent SQL injection.

---

## 🧪 Testing

```bash
make test
```

---

## 🔧 Useful Make Commands

```bash
make help          # List all commands
make run           # Run in dev mode with hot reload (Air)
make build         # Build binary
make test          # Run tests
make docker-up     # Start Docker services
make docker-down   # Stop Docker services
make migrate-up    # Apply DB migrations
make migrate-down  # Rollback DB migrations
make clean         # Remove build artifacts
make fmt           # gofmt
make tidy          # go mod tidy
```

---

## 🧭 Roadmap (high level)

> This is an indicative roadmap and may evolve with real deployments.

- **Phase 1**: Auth, user/role/permission, patients, appointments, visits, ICD‑10, medications, prescriptions.
- **Phase 2**: Paraclinical (lab tests, imaging), inventory & dispensing.
- **Phase 3**: Inpatient (admissions, beds, nursing notes), billing & accounting, internal insurance flows.
- **Phase 4** (future):
  - Surgery module.
  - Dedicated reporting & dashboards.
  - Data exports / aggregations for BI & analytics.

---

## 🤝 Contributing

1. Fork this repository.
2. Create a feature/bugfix branch.
3. Commit changes in small, logical steps with clear messages.
4. Push to your fork and open a Pull Request.

---

## 📄 License

This project is licensed under the **MIT License**
