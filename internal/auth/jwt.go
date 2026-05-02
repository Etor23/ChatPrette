package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// JWTManager genera y valida JWTs propios del backend
type JWTManager struct {
	secretKey string
	expiresIn time.Duration
}

// CustomClaims contiene los claims personalizados de nuestro JWT
type CustomClaims struct {
	UID      string `json:"uid"`
	Email    string `json:"email"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// NewJWTManager crea un nuevo manager de JWTs
// secretKey: clave secreta para firmar tokens (usa una clave fuerte, idealmente desde variables de ambiente)
// expiresIn: duración de expiración del token
func NewJWTManager(secretKey string, expiresIn time.Duration) *JWTManager {
	return &JWTManager{
		secretKey: secretKey,
		expiresIn: expiresIn,
	}
}

// GenerateToken genera un nuevo JWT con los datos del usuario
func (m *JWTManager) GenerateToken(uid, email, username string) (string, error) {
	now := time.Now()
	expiresAt := now.Add(m.expiresIn)

	claims := &CustomClaims{
		UID:      uid,
		Email:    email,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			Issuer:    "chat-back",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(m.secretKey))
	if err != nil {
		return "", fmt.Errorf("error signing token: %w", err)
	}

	return signed, nil
}

// VerifyToken valida un JWT y retorna los claims si es válido
func (m *JWTManager) VerifyToken(tokenString string) (*CustomClaims, error) {
	claims := &CustomClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Validar que el algoritmo sea el esperado
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.secretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("error parsing token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
