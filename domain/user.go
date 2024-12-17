package domain

// User represents a user in db.
//
// TODO: use another struct for user creation, without ID.
type User struct {
	ID       string `json:"id"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password"  validate:"required,min=8"`
}
