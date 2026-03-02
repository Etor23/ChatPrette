package dto

type CreateUserRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Username  string `json:"username" binding:"required"`
	AvatarURL string `json:"avatar_url,omitempty"`
	BirthDate string `json:"birth_date" binding:"required"`
}

type UpdateUserRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Username  string `json:"username" binding:"required"`
	AvatarURL string `json:"avatar_url,omitempty"`
	BirthDate string `json:"birth_date" binding:"required"`
}

type UserResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatar_url,omitempty"`
	BirthDate string `json:"birth_date"`
}
