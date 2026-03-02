package dto

type LoginRequest struct {
	IDToken string `json:"id_token" binding:"required"`
}

type AuthResponse struct {
	User  UserResponse `json:"user"`
	IsNew bool         `json:"is_new"` // true si se registró ahora
}
