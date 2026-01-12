# freerider-rest-api

A REST API that acts as a wrapper around the Freerider public API. It provides enhanced filtering capabilities and a "watch" feature to monitor for new rides matching specific criteria.

## Features

- **Advanced Filtering**: Search for rides using multiple origins or destinations.
- **Trip Watching**: Real-time notifications via Server-Sent Events (SSE) when new rides matching your criteria are found.
- **Location Discovery**: Easily retrieve all available pickup and drop-off locations.

## API Endpoints

### 1. Get Trips
Returns a list of available trips based on the provided filters.

- **URL**: `/trips`
- **Method**: `GET`
- **Query Parameters**:
    - `origin` (optional, multiple): The starting point(s) of the ride.
    - `destination` (optional, multiple): The destination(s) of the ride.
    - `startDate` (optional): The earliest date and time (Format: `2006-01-02T15:04:05` or `2006-01-02`).
    - `endDate` (optional): The latest date and time (Format: `2006-01-02T15:04:05` or `2006-01-02`).

### 2. Watch Trips
Monitor for new trips. This endpoint uses Server-Sent Events (SSE) to stream data back to the client.

- **URL**: `/watch`
- **Method**: `POST`
- **Body** (JSON):
  ```json
  {
    "origin": "Stockholm",
    "destination": "Oslo",
    "minDate": "2026-01-12T10:00:00",
    "maxDate": "2026-01-15T10:00:00"
  }
  ```
- **Response**: A stream of `ride-found` events.

### 3. Get Locations
Returns all valid locations known to the Freerider system.

- **URL**: `/locations`
- **Method**: `GET`

## Run

### Prerequisites
- Go 1.x

### Running the server
```bash
go run cmd/api-server/main.go
```