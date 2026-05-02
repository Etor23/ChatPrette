package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	ContextUserID   = "user_id"
	ContextEmail    = "email"
	ContextUsername = "username"
)

// Middleware verifica el JWT propio del backend en cada request
func Middleware(manager *JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {

		// 1. Obtener header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Falta el header Authorization",
			})
			return
		}

		// 2. Extraer token: "Bearer eyJhb..."
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Formato inválido. Usa: Bearer <token>",
			})
			return
		}

		// 3. Verificar JWT
		claims, err := manager.VerifyToken(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Token inválido o expirado",
			})
			return
		}

		// 4. Guardar en contexto para que los handlers lo usen
		c.Set(ContextUserID, claims.UID)
		c.Set(ContextEmail, claims.Email)
		c.Set(ContextUsername, claims.Username)

		c.Next()
	}
}

// Helpers para extraer datos del contexto en los handlers
func GetUserID(c *gin.Context) string {
	uid, _ := c.Get(ContextUserID)
	if uid == nil {
		return ""
	}
	return uid.(string)
}

func GetEmail(c *gin.Context) string {
	email, _ := c.Get(ContextEmail)
	if email == nil {
		return ""
	}
	return email.(string)
}

func GetUsername(c *gin.Context) string {
	username, _ := c.Get(ContextUsername)
	if username == nil {
		return ""
	}
	return username.(string)
}

// Backward compatible alias
func GetFirebaseUID(c *gin.Context) string {
	return GetUserID(c)
}
