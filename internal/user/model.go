package user

// CreateUserRequest represents the request to create a user
type CreateUserRequest struct {
	Name     string `json:"name" validate:"required,min=1,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password,omitempty" validate:"omitempty,min=8,max=72"`
}

// UserResponse represents a user in responses
type UserResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
