# Validation Example

This project demonstrates a custom validation middleware implementation in Go using reflection and struct tags, without external dependencies.

## Features

- Custom struct validation using tags
- HTTP middleware for request validation
- Support for multiple validation rules:
  - `required`: Field cannot be empty
  - `email`: Must match email format
  - `min`: Minimum length requirement
- Clean separation of concerns
- Fully tested with unit and integration tests

## Usage

### Unit Test

```bash
make test
```

### Execution

```bash
go run .
```

The server will start on port 8080.

### API Endpoints

#### Create User
```bash
POST /users
Content-Type: application/json

{
    "Email": "test@example.com",
    "Password": "password123"
}
```

### Validation Rules

Add validation rules to your structs using tags:

```go
type User struct {
    Email    string `validate:"required,email"`
    Password string `validate:"required,min=8"`
}
```

### Testing

Run all tests:
```bash
go test ./...
```

## Examples

### Valid Request
```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{
    "Email": "test@example.com",
    "Password": "password123"
  }'
```

### Invalid Request (Validation Error)
```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{
    "Email": "invalid-email",
    "Password": "short"
  }'
```

