# Dormitory Backend

Go + Gin + PostgreSQL backend for the dormitory management system.

## Run

Prepare the database first from the repository root:

```bash
PGPASSWORD=10928 psql -U turing -d student_dormitory -f sql/001_create_student_dormitory_schema.sql
```

Start the API:

```bash
cd backend
DORM_JWT_SECRET='e671eb179af89abbd4bb1e264799e47a3888bc5ac10d13add8e24d2a4cc3ed15' go run ./cmd/server
```

The service listens on `http://localhost:8080` by default. Health check:

```bash
curl http://localhost:8080/api/health
```

## Configuration

Defaults live in `config/config.yaml`. Environment variables use the `DORM_` prefix, for example:

- `DORM_SERVER_PORT=8081`
- `DORM_DATABASE_HOST=localhost`
- `DORM_DATABASE_PASSWORD=10928`
- `DORM_JWT_SECRET=...`

## API Documentation

Frontend integration guide: [`docs/API.md`](docs/API.md)

## Implemented Areas

- JWT login and refresh
- Role middleware for student, repair staff, cleaning staff, dormitory manager, and system admin
- Student profile, survey, roommate, request overview, leave, late return, room change, off-campus, repair, cleaning, and payment endpoints
- Repair and cleaning status flows
- Manager review endpoints for applications and work orders
- Allocation approval transaction with bed occupation
- Room change and off-campus approval transactions with bed updates
- Notifications, dashboard views, available beds, room balance
- Attachment upload/download with PostgreSQL `BYTEA` storage
