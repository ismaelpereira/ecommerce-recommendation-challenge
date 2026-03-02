# 🛒 Ecommerce Events Analytics Service

A Go-based backend service that ingests ecommerce events and provides analytical queries for top products per store within a configurable time window.

Built with:

- Go
- Gin
- Google Cloud Bigtable (event storage)
- Google BigQuery (analytics)
- Clean layered architecture (Handler → Service → Repository)

# 📂 Project Structure

```
.
├── bigquery
│ └── client.go
├── bigtable
│ └── client.go
├── config
│ └── config.go
├── DESIGN.md
├── go.mod
├── go.sum
├── handlers
│ └── handlers.go
├── main.go
├── README.md
├── repository
│ ├── bigquery.go
│ └── bigtable.go
├── service
│ └── service.go
└── types
  └── types.go

8 directories, 13 files
```

# 🚀 Api Endpoints

## Create Event

**POST /events**

### Request Body

```json
{
   "user_id": "user_123",
   "product_id": "product_789",
   "store_id": "store_123",
   "event_type": "view"
}
```

#### Supported Event types

- view
- add_to_cart
- purchase

### Behavior

- Validate the input (all fields required, event_type must be one of: view, add_to_cart, purchase)
- Write to Bigtable
- Write to BigQuery
- Return 201 Created on success with the generated event ID
- Return appropriate error responses for invalid input or failures

## Get Top Products Per Store

**GET /analytics/top-products?store_id=<store_id>&hours=<number_of_hours>**

### Behavior

Query BigQuery for the top 10 most-viewed products for the given store within the last N hours

## Get Events From User

**GET /events/user/:user_id?limit=<max_of_events>**

### Behavior

- Read from Bigtable (not BigQuery) to fetch recent events for a user
- Return the most recent events according to the limit query, ordered newest first
- This endpoint must be fast (<100ms target for typical loads)

## Health Check

**GET /health**

### Behavior

- Returns 200 OK if the service can reach both Bigtable and BigQuery
- Returns 503 Service Unavailable with details if either is unreachable

# 🧪 Running Locally (Bigtable Emulator)

1. Start Big Table Emulator

```bash
gcloud beta emulators bigtable start
```

2. Export Environment variable

```bash
export BIGTABLE_EMULATOR_HOST=localhost:8086
```

3. Create Table

```bash
cbt -instance=local-instance createtable events
cbt -instance=local-instance createfamily events events
```

4. Run the service

```bash
go run main.go
```

# 🤖 Running Tests

Run on your terminal

```bash
go test -v ./...
```

Running Integration tests

```bash
go test ./... -v -tags=integration
```


# 🔨 Run Using Makefile

Building the service

```bash
make build
```

Running the service

```bash
make run
```

Testing the service

```bash
make test
```

Running test with coverage

```bash
make cover
```

Running Integration test

```bash
make test-integration
```