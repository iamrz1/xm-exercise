# Company Management Microservice

A production-ready RESTful microservice for managing company information.

## Features

- CRUD operations for companies
- JWT authentication
- Event production using Kafka
- Containerized with Docker and docker-compose
- SQL Database integration
- Input validation
- Graceful shutdown
- Clean architecture design
- Structured logging with Zap
- API Documentation with Swagger

## Requirements

- Docker and Docker Compose
- Go 1.21+ (for local development)
- make

## Running the Application

### Using Docker Compose

1. Clone the repository:
   ```
   git clone https://github.com/iamrz1/xm-exercise.git
   cd xm-exercise
   ```

2. Start the application:
   ```
   make up
   ```

3. Once an image fo this app is built, a new image is not built from code again.
Run the following to build a new image if any changes is made to the source code:
4. Run the application:
   ```
   make build
   ```

5. The API will be available at `http://localhost:8080`
6. Swagger documentation is available at `http://localhost:8080/swagger/index.html`
7. Exit the application:
   ```
   make down
   ```

### Local Development

1. Install dependencies:
   ```
   make dep
   ```

2. Create a `.env` file with the following content:
   ```
   PORT=8080
   APP_ENV=dev
   LOG_LEVEL=debug
   DATABASE_DIALECT=postgres
   DATABASE_URL=postgres://postgres:postgres@localhost:5432/companydb?sslmode=disable
   JWT_SECRET=your-super-secret-key-change-in-production
   JWT_EXPIRATION_HOURS=24
   KAFKA_BROKERS=localhost:9092
   ```

3. Generate Swagger documentation:
   ```
   make doc
   ```

4. Ensure dependencies:Ensure 
   Make sure that the database and kafka configuration provided via env or config file are valid and present

5. Run the application:
   ```
   make run
   ```

## API Documentation

The API is documented using Swagger. When the application is running, you can access 
the Swagger UI at the `{host}/swagger/index.html` endpoint.


## API Endpoints Overview

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