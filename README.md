# Home API

A simple backend service that integrates with the PokéAPI, stores data in MySQL, and implements Redis caching.

## Features

- External API integration with PokéAPI
- `/sync` endpoint to fetch and store data from the external API
- `/items` endpoint to retrieve stored data with Redis caching
- Background job that refreshes data every 15 minutes
- Idempotent writes to prevent duplicates
- Graceful handling of API failures with retry and backoff
- Docker setup for MySQL and Redis

## Tech Stack

- Go Lang
- MySQL
- Redis

## How to Run the Project

### Using Makefile (Recommended)

The project includes a Makefile that simplifies both Docker and local development workflows.

#### Docker Commands

```bash
# Build Docker images
make build

# Start all services (MySQL, Redis, API)
make up

# View logs from all containers
make logs

# Stop and remove all containers
make down
```

#### Local Development Commands

```bash
# Install dependencies
make install

# Set up environment file
make env-setup

# Run the application locally
make run

# Run tests
make test
```

For a complete list of available commands:

```bash
make help
```

### Using Docker Directly

1. Clone the repository
2. Navigate to the project directory
3. Start the services using Docker Compose:

```bash
docker-compose up -d
```

The API will be available at http://localhost:8080

### Manual Setup

1. Clone the repository
2. Set up MySQL and Redis locally
3. Copy `.env.example` to `.env` and update the configuration
4. Build and run the application:

```bash
go build -o home-api ./cmd/api
./home-api
```

## API Endpoints

- `POST /sync` - Fetches data from the external API and stores it in MySQL
- `GET /items` - Returns the stored data from MySQL (with Redis caching)

## SQL Queries for Orders

The following SQL queries are implemented in the `order_repository.go` file:

### Number of orders and total amount per status in the last 30 days

```sql
SELECT 
    status, 
    COUNT(*) as count, 
    SUM(amount) as total_amount
FROM orders
WHERE created_at >= DATE_SUB(NOW(), INTERVAL 30 DAY)
GROUP BY status
ORDER BY status
```

### Top 5 customers by total spend

```sql
SELECT 
    customer_id, 
    SUM(amount) as total_spend
FROM orders
WHERE status = 'PAID'
GROUP BY customer_id
ORDER BY total_spend DESC
LIMIT 5
```

## Assumptions Made

1. The external API (PokéAPI) is publicly accessible and doesn't require authentication
2. The data structure from the external API is consistent
3. MySQL and Redis are running on default ports
4. The application will handle a moderate amount of traffic

## Trade-offs and Future Improvements

With more time, the following improvements could be made:

1. **Authentication and Authorization**: Add JWT-based authentication for API endpoints
2. **Pagination**: Implement pagination for the `/items` endpoint to handle large datasets
3. **Metrics and Monitoring**: Add Prometheus metrics and Grafana dashboards
4. **Rate Limiting**: Implement rate limiting for API endpoints
5. **Unit and Integration Tests**: Add comprehensive test coverage
6. **API Documentation**: Add Swagger documentation for the API endpoints
7. **Circuit Breaker**: Implement a circuit breaker pattern for external API calls
8. **Distributed Locking**: Use Redis for distributed locking in the background job
9. **Structured Logging**: Implement structured logging with correlation IDs
10. **Health Checks**: Add health check endpoints for the application and dependencies