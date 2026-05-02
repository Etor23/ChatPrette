package dto

// Solicitud de registro: email, password, username
type RegisterRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	Username  string `json:"username" binding:"required,min=3,max=32"`
	AvatarURL string `json:"avatar_url"`
	Birthdate string `json:"birthDate,omitempty"` // formato ISO 8601 o YYYY-MM-DD
}

// Solicitud de login: email, password
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// Respuesta de registro/login: token + user + is_new
type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
	IsNew bool         `json:"is_new"`
}

// Respuesta de logout
type LogoutResponse struct {
	Message string `json:"message"`
}

// Respuesta de refresh token
type RefreshTokenResponse struct {
	Token     string `json:"token"`
	ExpiresIn int64  `json:"expires_in"`
}
