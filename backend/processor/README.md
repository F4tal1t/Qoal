# Qoal File Processor Microservice

A Go-based microservice for processing file conversion jobs.

## Features
- File upload/download from S3
- Redis queue for job processing
- REST API for job management

## Setup
1. Install dependencies:
```
go mod download
```

2. Configure environment variables:
```
cp .env.example .env
```

3. Run locally:
```
make run
```

## Docker
Build and run with Redis:
```
docker-compose up --build
```

## API Endpoints
- `POST /process`: Submit a new file processing job
- `GET /status/:id`: Check job status

## Development
- Run tests: `make test`
- Build binary: `make build`