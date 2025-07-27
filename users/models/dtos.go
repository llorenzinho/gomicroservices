package models

type CreateUserRequest struct {
	Username        string `json:"username" binding:"required,min=3,max=32"`
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=6"`
	ConfirmPassword string `json:"confirmPassword" binding:"required,eqfield=Password"`
}

type UpdateUserRequest struct {
	Username string `json:"username,omitempty" binding:"min=3,max=32"`
	Email    string `json:"email,omitempty" binding:"email"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password" binding:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh" binding:"required"`
}
