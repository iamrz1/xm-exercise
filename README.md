# Company Management Microservice

A production-ready RESTful microservice for managing company information.

## Features

- CRUD operations for companies
- JWT authentication
- Event production using Kafka
- Containerized with Docker and docker-compose
- PostgreSQL database integration
- Input validation
- Graceful shutdown
- Clean architecture design

## Requirements

- Docker and Docker Compose
- Go 1.21+ (for local development)

## Running the Application

### Using Docker Compose

1. Clone the repository:
   ```
   git clone https://github.com/iamrz1/xm-exercise.git
   cd xm-exercise
   ```

2. Start the application:
   ```
   docker-compose up -d
   ```

3. The API will be available at `http://localhost:8080`

### Local Development

1. Install dependencies:
   ```
   go mod download
   ```

2. Create a `.env` file with the following content:
   ```
   PORT=8080
   DATABASE_URL=postgres://postgres:postgres@localhost:5432/companydb?sslmode=disable
   JWT_SECRET=your-super-secret-key-change-in-production
   JWT_EXPIRATION_HOURS=24
   KAFKA_BROKERS=localhost:9092
   ```

3. Run the application:
   ```
   go run cmd/server/main.go
   ```

## API Endpoints

### Authentication

- **POST /api/v1/auth/register** - Register a new user
- **POST /api/v1/auth/login** - Login and get JWT token

### Companies

All company endpoints require JWT authentication.

- **POST /api/v1/companies** - Create a new company
- **GET /api/v1/companies/{id}** - Get company by ID
- **PATCH /api/v1/companies/{id}** - Update company
- **DELETE /api/v1/companies/{id}** - Delete company

## Testing

Run the tests with: 