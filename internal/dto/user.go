package dto

// CreateUserRequest represents the request body for creating a new user
type CreateUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

// CreateUserResponse represents the
type CreateUserResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}
