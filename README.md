# Validation Example

This project demonstrates a custom validation middleware implementation in Go using reflection and struct tags.

## Usage

### Execution

```bash
make docker-start docker-logs
```

The server will start on port 8080 (read from [.env](.env)).

To stop services, run this command:

```bash
make docker-stop
```

### Unit Test

```bash
make test
```

### API Endpoints

#### Create User
```
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


## Examples

### Valid Request
```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

### Invalid Request (Validation Error)
```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "invalid-email",
    "password": "short"
  }'
```

