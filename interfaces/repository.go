package interfaces

import (
	"context"
	"github.com/mauleyzaola/validation/domain"
)

//go:generate moq -out ../mocks/repository.go -pkg mocks . Repository
type Repository interface {
	// CreateUser inserts one user in db and returns its generated uuid v4, and true if operation was successful.
	CreateUser(ctx context.Context, user domain.User) (string, bool, error)
}
