package dto

type CreateUserRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Username  string `json:"username" binding:"required"`
	AvatarURL string `json:"avatar_url,omitempty"`
}

type UserResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatar_url,omitempty"`
	Birthdate string `json:"birthdate,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
}

// UpdateProfileRequest para cambiar username o birthdate
type UpdateProfileRequest struct {
Username  string `json:"username,omitempty"`  // opcional
Birthdate string `json:"birthdate,omitempty"` // opcional
AvatarURL string `json:"avatar_url,omitempty"` // opcional
}
