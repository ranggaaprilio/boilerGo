# BoilerGo

A Go boilerplate project using Echo framework, GORM with MySQL, and Viper for configuration management. Provides a clean architecture pattern with separation of concerns between handlers, services, and repositories.

## Features

- **Clean Architecture**: Follows a modular design pattern
- **Echo Framework**: Fast and minimalist web framework
- **GORM**: Powerful ORM for MySQL
- **Viper**: Configuration management
- **Request Validation**: Built-in request validation
- **Middleware Support**: Health checks and more
- **Standardized Response Format**: Consistent API responses
- **API Documentation**: Swagger/OpenAPI documentation with interactive UI

## API Documentation

The API is documented using Swagger/OpenAPI. Once the application is running, you can access the interactive API documentation at:

```
http://localhost:8080/swagger/index.html
```

To regenerate the API documentation after making changes to your handlers:

```bash
# Using the provided script
./scripts/generate-swagger.sh

# Or using make
make swagger
```

## Prerequisites

- Go 1.24 or higher
- MySQL database
- Docker and Docker Compose (for containerized setup)

## Installation

### Clone the repository

```bash
git clone https://github.com/ranggaaprilio/boilerGo.git
cd boilerGo
```

### Set up configuration

Copy the example configuration file and modify it according to your environment:

```bash
cp config.yml.example config.yml
```

Edit the `config.yml` file to configure your database connection and other settings.

### Install dependencies

```bash
go mod download
```

Or use:

```bash
go mod tidy
```

This will download all required dependencies listed in the go.mod file.

### Build the application

```bash
go build -o boilerGo
```

### Run the application

```bash
./boilerGo
```

## Running with Docker

This project includes Docker and Docker Compose configuration for easy setup. This creates a containerized environment with both the API and MySQL database.

### Start the containers

```bash
# Build and start all services
make docker-build
make docker-up

# Or directly with docker-compose
docker-compose up -d
```

### Check logs

```bash
make docker-logs
```

### Stop the containers

```bash
make docker-down
```

The API will be available at `http://localhost:8080`, and Swagger documentation at `http://localhost:8080/swagger/index.html`.

MySQL will be accessible on port 3306 with the following credentials:

- Username: rangga
- Password: verysecret
- Database: boilergo

## Project Structure

```
├── app/                # Application core
│   └── v1/             # API version 1
│       ├── handler/    # HTTP handlers
│       └── modules/    # Business modules
├── config/             # Configuration management
├── docs/               # API documentation
│   ├── swagger.json    # Generated Swagger JSON definition
│   ├── swagger.yaml    # Generated Swagger YAML definition
│   └── swagger_guide.md # Swagger usage guide
├── exception/          # Error handling
├── helper/             # Helper functions
├── middleware/         # HTTP middleware
├── public/             # Public assets
├── routes/             # Route definitions
├── scripts/            # Utility scripts
├── bootstrap.go        # Application bootstrap
├── main.go             # Entry point
└── config.yml.example  # Example configuration
```

## API Documentation

### User Endpoints

#### Register a new user

```
POST /api/v1/users
```

Request body:

```json
{
  "name": "John Doe"
}
```

Response:

```json
{
  "code": 200,
  "message": "Success save data",
  "data": {
    "ID": 1,
    "CreatedAt": "2025-06-15T10:00:00Z",
    "UpdatedAt": "2025-06-15T10:00:00Z",
    "DeletedAt": null,
    "Name": "John Doe"
  }
}
```

Error responses:

- Bad Request (400): When form binding fails or validation fails
- Internal Server Error (500): When saving data fails

## Architecture

The project follows a clean architecture pattern with the following components:

### Handler Layer

Handles HTTP requests and responses, input validation, and calls the appropriate service methods.

### Service Layer

Contains business logic and coordinates calls between different repositories.

### Repository Layer

Handles data access operations and interactions with the database.

### Entity

Defines the data structures and models used in the application.

## Documentation

Detailed documentation is available in the `docs` directory:

- [User API Documentation](docs/user_api.md): Detailed information about the User API endpoints
- [Architecture Documentation](docs/architecture.md): Overview of the application architecture and design patterns

### API Documentation with Swagger

This project includes Swagger/OpenAPI documentation. To access the Swagger UI:

1. Start the application:

   ```bash
   ./boilerGo
   ```

2. Open the Swagger UI in your browser:

   ```
   http://localhost:8080/swagger/index.html
   ```

3. If you need to regenerate the Swagger documentation after making changes:
   ```bash
   ./scripts/generate-swagger.sh
   ```
   or you can use
   ```bash
   make swagger
   ```

## Troubleshooting

### Database Connection Issues in Docker

If you encounter a database connection error when running in Docker with a message like:

```
Failed to initialize database, got error dial tcp :0: connect: connection refused
```

This might be due to one of these issues:

1. **MySQL container isn't ready yet**: Even with the healthcheck in place, sometimes the MySQL container might not be fully ready. Try restarting just the API container:

   ```bash
   docker-compose restart api
   ```

2. **Configuration issue**: Make sure your `config.docker.yml` file has these correct settings:

   ```yaml
   database:
     dbusername: "rangga"
     dbpassword: "verysecret"
     dbhost: "mysql"
     dbport: "3306"
     dbname: "boilergo"
   ```

3. **Network issue**: Ensure both services are on the same Docker network. You can check this with:
   ```bash
   docker network inspect boilergo-network
   ```

The application has been configured to read both from the configuration file and from environment variables, so either approach should work.

## License

See the [LICENSE](LICENSE) file for details.
