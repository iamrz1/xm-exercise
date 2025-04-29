# Company Management Microservice

A production-ready(ish) RESTful microservice for managing company information.

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
- Unit tests
- Integration tests
- Lint (golangci-lint)
- `.env` as the config file

## Requirements

- Docker and Docker Compose
- make
- Go 1.23+ (for local development)
- golangci-lint (for local development)

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
   ```
   make build
   ```

4. The API will be available at `http://localhost:8080`
5. Swagger documentation is available at `http://localhost:8080/swagger/index.html`
6. Exit the application:
   ```
   make down
   ```
7. Cleanup leftover resources:
   ```
   make cleanup
   ```

### Local Development

1. Install dependencies:
   ```
   make dep
   ```

2. Create a `.env` file with the following content to configure the app:
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

## Linting
Use `golangci-lint run` to check any linter or formatter related issue.
`golangci-lint` is also baked into the `Dockerfile` for seamless integration
in docker based implementation.

## Testing

### Run unit tests with:
   ```
   make test
   ```

### Running integration test:

#### Docker Compose
1. Run the app and integration test using the `docker-compose-e2e.yml` file.
It will ensure all the dependencies and then run integration tests by itself .
   ```
   make integration-test
   ```

2. If any error is encountered error will show up in console. If there;s no error,
`E2E Test Sequence Completed SUCCESSFULLY!` will show up.
3. Exit the application:
   ```
   make down
   ```
4. Cleanup leftover resources:
   ```
   make cleanup
   ```

#### Locally

1. Ensure all the dependencies are in place (database, kafka)
2. Set `APP_ENV` to `int`
3. Run the application:
   ```
   make run
   ```
4. If any error is encountered error will show up in console. If there;s no error,
   `E2E Test Sequence Completed SUCCESSFULLY!` will show up.
5. Press `crtl+c` to exit