# Qoal Backend Authentication System

## Overview
This document describes the JWT-based authentication system implemented for the Qoal file processing backend.

## Features
- JWT-based authentication with 24-hour token expiration
- User registration and login
- Protected API endpoints
- PostgreSQL database integration with GORM
- CORS support for frontend integration

## API Endpoints

### Public Endpoints

#### Register User
```
POST /api/v1/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123",
  "name": "John Doe"
}

Response:
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "name": "John Doe",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  },
  "expires_at": "2024-01-02T00:00:00Z"
}
```

#### Login User
```
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}

Response:
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "name": "John Doe",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  },
  "expires_at": "2024-01-02T00:00:00Z"
}
```

### Protected Endpoints

All protected endpoints require the Authorization header:
```
Authorization: Bearer <your-jwt-token>
```

#### Get User Profile
```
GET /api/v1/auth/profile
Authorization: Bearer <token>

Response:
{
  "user": {
    "id": 1,
    "email": "user@example.com",
    "name": "John Doe",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

#### Create Processing Job (Protected)
```
POST /api/v1/process
Authorization: Bearer <token>
Content-Type: application/json

{
  "input_file": "/path/to/file.pdf",
  "target_format": "docx"
}
```

#### Get Job Status (Protected)
```
GET /api/v1/status/{job_id}
Authorization: Bearer <token>
```

## Environment Variables

Add these to your `.env` file:

```env
# Database Configuration
DATABASE_URL=postgres://username:password@localhost:5432/qoal_db
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
REDIS_URL=localhost:6379
```

## Database Schema

The authentication system uses a `users` table with the following schema:

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    name VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## Security Features

1. **Password Hashing**: All passwords are hashed using bcrypt
2. **JWT Tokens**: Tokens expire after 24 hours
3. **CORS Protection**: Configured for frontend integration
4. **Input Validation**: Request validation using Gin bindings

## Testing

Run the authentication tests:
```bash
go test ./tests -v
```

## Error Responses

All endpoints return consistent error responses:

```json
{
  "error": "Error message description"
}
```

Common HTTP status codes:
- `200 OK`: Successful request
- `201 Created`: Resource created successfully
- `400 Bad Request`: Invalid request data
- `401 Unauthorized`: Invalid or missing authentication
- `500 Internal Server Error`: Server error

## Next Steps

The authentication system is now complete. The next steps are:

1. **File Upload System**: Implement multipart file upload endpoints
2. **Local Storage**: Replace S3 with local file storage
3. **Rate Limiting**: Add rate limiting middleware
4. **User Limits**: Implement user-specific processing limits
5. **File Processing**: Integrate actual file processing with authentication