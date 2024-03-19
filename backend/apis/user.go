package apis

import "fmt"

type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Name     string `json:"name" binding:"required"`
}

func (r *CreateUserRequest) Validate() error {
	if r.Username == "" {
		return fmt.Errorf("username is required")
	}

	if r.Email == "" {
		return fmt.Errorf("email is required")
	}

	if r.Password == "" {
		return fmt.Errorf("password is required")
	}

	if r.Name == "" {
		return fmt.Errorf("name is required")
	}

	return nil
}

type CreateUserResponse struct {
	UserID string `json:"user_id"`
}

type GetUserResponse struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Name     string `json:"name"`
}

type UpdateUserRequest struct {
	Name string `json:"name"`
}

type UpdateUserResponse struct {
	UserID string `json:"user_id"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token     string `json:"token"`
	Type      string `json:"type"`
	ExpiresIn int    `json:"expires_in"`
}
