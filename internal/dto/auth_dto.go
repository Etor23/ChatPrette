package dto

type RegisterRequest struct {
	Username  string `json:"username"   binding:"required,min=3,max=32"`
	AvatarURL string `json:"avatar_url"`
}

type LoginResponse struct {
	User  UserResponse `json:"user"`
	IsNew bool         `json:"is_new"`
}

// reutiliza el UserResponse que ya se creó en user_dto.go
